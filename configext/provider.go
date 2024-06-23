package configext

import (
	"context"
	"fmt"
	"github.com/aesoper101/x/fileutil"
	"io/fs"
	"log/slog"
	"net/url"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/inhies/go-bytesize"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

type tuple struct {
	Key   string
	Value interface{}
}

type eventChanel chan struct{}

type Provider struct {
	l      sync.RWMutex
	logger *slog.Logger
	*koanf.Koanf
	immutables, exceptImmutables []string

	onChanges []func(event interface{}, err error)

	forcedValues []tuple
	baseValues   []tuple
	files        []string

	flags *pflag.FlagSet

	enableEnvLoading bool
	envPrefix        string

	providers     []koanf.Provider
	userProviders []koanf.Provider

	decodeHookFunc mapstructure.DecodeHookFunc
}

const (
	FlagConfig = "config"
	Delimiter  = "."
)

func RegisterConfigFlag(flags *pflag.FlagSet, fallback []string) {
	flags.StringSliceP(FlagConfig, "c", fallback, "Config files to load, overwriting in the order specified.")
}

func New(ctx context.Context, options ...ProviderOptions) (*Provider, error) {
	p := &Provider{
		logger: slog.Default(),
		decodeHookFunc: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToIPHookFunc(),
			mapstructure.StringToNetIPAddrHookFunc(),
			mapstructure.StringToNetIPAddrPortHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToSliceHookFunc(","),
			StringToMailAddressHookFunc(),
			StringToRegexpHookFunc(),
			StringToURLHookFunc(),
			JsonUnmarshalerHookFunc(),
			TextUnmarshalerHookFunc(),
		),
		Koanf: koanf.NewWithConf(koanf.Conf{Delim: Delimiter, StrictMerge: true}),
	}
	for _, option := range options {
		option(p)
	}

	providers, err := p.createProviders(ctx)
	if err != nil {
		return nil, err
	}

	p.providers = providers

	k, err := p.newKoanf()
	if err != nil {
		return nil, err
	}

	p.replaceKoanf(k)

	return p, nil
}

func (p *Provider) createProviders(ctx context.Context) (providers []koanf.Provider, err error) {
	for _, t := range p.baseValues {
		providers = append(providers, NewKoanfConfmap([]tuple{t}))
	}

	paths := p.files
	if p.flags != nil {
		p, _ := p.flags.GetStringSlice(FlagConfig)
		paths = append(paths, p...)
	}

	p.logger.DebugContext(ctx, "Adding config files.", "files", paths)

	c := make(eventChanel)
	go p.watchForFileChanges(ctx, c)

	files := getFiles(paths, isConfigFile)
	for _, path := range files {
		fp, err := NewKoanfFile(path)
		if err != nil {
			return nil, err
		}
		err = fp.Watch(func(event interface{}, err error) {
			if err != nil {
				p.logger.ErrorContext(ctx, "Failed to watch file.", "file", path)
			} else {
				c <- struct{}{}
			}
		})
		if err != nil {
			return nil, err
		}
		providers = append(providers, fp)
	}

	providers = append(providers, p.userProviders...)

	if p.flags != nil {
		pp := posflag.Provider(p.flags, Delimiter, p.Koanf)
		providers = append(providers, pp)
	}

	if p.enableEnvLoading {
		envProvider := env.Provider(p.envPrefix, Delimiter, func(s string) string {
			return strings.Replace(strings.ToLower(
				strings.ReplaceAll(strings.TrimPrefix(s, p.envPrefix), " ", "_")), "_", ".", -1)
		})

		providers = append(providers, envProvider)
	}

	for _, t := range p.forcedValues {
		providers = append(providers, NewKoanfConfmap([]tuple{t}))
	}

	return providers, nil
}

func (p *Provider) newKoanf() (_ *koanf.Koanf, err error) {
	k := koanf.New(Delimiter)

	for _, provider := range p.providers {
		if _, ok := provider.(*posflag.Posflag); ok {
			provider = posflag.Provider(p.flags, ".", k)
		}

		var opts []koanf.Option
		if _, ok := provider.(*env.Env); ok {
			opts = append(opts, koanf.WithMergeFunc(MergeAllTypes))
		}

		if err := k.Load(provider, nil, opts...); err != nil {
			return nil, err
		}
	}

	return k, nil
}

func (p *Provider) replaceKoanf(k *koanf.Koanf) {
	p.Koanf = k
	p.replaceVars(k)
}

func (p *Provider) replaceVars(k *koanf.Koanf) {
	keys := k.Keys()

	for _, key := range keys {
		value := k.String(key)
		if value == "" {
			continue
		}

		reg := regexp.MustCompile(`\${([^}]+)}`)
		matches := reg.FindAllStringSubmatch(value, -1)
		for _, match := range matches {
			if len(match) != 2 {
				continue
			}
			// 替换值中存在的 ${} 包裹的变量
			valPath := strings.TrimSpace(match[1])

			if k.Exists(valPath) {
				value = strings.Replace(value, match[0], k.String(valPath), -1)
			} else if p.enableEnvLoading {
				envKey := strings.Replace(strings.ToLower(
					strings.TrimPrefix(valPath, p.envPrefix)), "_", ".", -1)
				envVal := k.String(envKey)
				if envVal != "" {
					value = strings.Replace(value, match[0], envVal, -1)
				}
			}
		}

		_ = k.Set(key, value)
	}

}

func (p *Provider) watchForFileChanges(ctx context.Context, c eventChanel) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c:
			p.logger.DebugContext(ctx, "File changed.")
			p.reload()
		}
	}
}

func (p *Provider) reload() {
	p.l.Lock()

	var err error
	defer func() {
		p.l.Unlock()
		p.runOnChanges(nil, err)
	}()

	nk, err := p.newKoanf()
	if err != nil {
		return
	}

	oldImmutables, newImmutables := p.Koanf.Copy(), nk.Copy()
	deleteOtherKeys(oldImmutables, p.immutables)
	deleteOtherKeys(newImmutables, p.immutables)

	for _, key := range p.exceptImmutables {
		oldImmutables.Delete(key)
		newImmutables.Delete(key)
	}
	if !reflect.DeepEqual(oldImmutables.Raw(), newImmutables.Raw()) {
		for _, key := range p.immutables {
			if !reflect.DeepEqual(oldImmutables.Get(key), newImmutables.Get(key)) {
				err = NewImmutableError(key, fmt.Sprintf("%v", p.Koanf.Get(key)), fmt.Sprintf("%v", nk.Get(key)))
				return // unlocks & runs changes in defer
			}
		}
	}

	p.replaceKoanf(nk)
}

func (p *Provider) runOnChanges(e interface{}, err error) {
	for k := range p.onChanges {
		p.onChanges[k](e, err)
	}
}

func deleteOtherKeys(k *koanf.Koanf, keys []string) {
outer:
	for _, key := range k.Keys() {
		for _, ik := range keys {
			if key == ik {
				continue outer
			}
		}
		k.Delete(key)
	}
}

// DirtyPatch patches individual config keys without reloading the full config
//
// WARNING! This method is only useful to override existing keys in string or number
// format. DO NOT use this method to override arrays, maps, or other complex types.
//
// This method DOES NOT validate the config against the config JSON schema. If you
// need to validate the config, use the Set method instead.
//
// This method can not be used to remove keys from the config as that is not
// possible without reloading the full config.
func (p *Provider) DirtyPatch(key string, value any) error {
	p.l.Lock()
	defer p.l.Unlock()

	t := tuple{Key: key, Value: value}
	kc := NewKoanfConfmap([]tuple{t})

	p.forcedValues = append(p.forcedValues, t)
	p.providers = append(p.providers, kc)

	if err := p.Koanf.Load(kc, nil, []koanf.Option{}...); err != nil {
		return err
	}

	return nil
}

func (p *Provider) Set(key string, value interface{}) error {
	p.l.Lock()
	defer p.l.Unlock()

	p.forcedValues = append(p.forcedValues, tuple{Key: key, Value: value})
	p.providers = append(p.providers, NewKoanfConfmap([]tuple{{Key: key, Value: value}}))

	k, err := p.newKoanf()
	if err != nil {
		return err
	}

	p.replaceKoanf(k)
	return nil
}

func (p *Provider) BoolF(key string, fallback bool) bool {
	p.l.RLock()
	defer p.l.RUnlock()

	if !p.Koanf.Exists(key) {
		return fallback
	}

	return p.Bool(key)
}

func (p *Provider) StringF(key, fallback string) string {
	p.l.RLock()
	defer p.l.RUnlock()

	if !p.Koanf.Exists(key) {
		return fallback
	}

	return p.String(key)
}

func (p *Provider) StringsF(key string, fallback []string) (val []string) {
	p.l.RLock()
	defer p.l.RUnlock()

	if !p.Koanf.Exists(key) {
		return fallback
	}

	return p.Strings(key)
}

func (p *Provider) IntF(key string, fallback int) (val int) {
	p.l.RLock()
	defer p.l.RUnlock()

	if !p.Koanf.Exists(key) {
		return fallback
	}

	return p.Int(key)
}

func (p *Provider) Float64F(key string, fallback float64) (val float64) {
	p.l.RLock()
	defer p.l.RUnlock()

	if !p.Koanf.Exists(key) {
		return fallback
	}

	return p.Float64(key)
}

func (p *Provider) DurationF(key string, fallback time.Duration) (val time.Duration) {
	p.l.RLock()
	defer p.l.RUnlock()

	if !p.Koanf.Exists(key) {
		return fallback
	}

	return p.Duration(key)
}

func (p *Provider) ByteSizeF(key string, fallback bytesize.ByteSize) bytesize.ByteSize {
	p.l.RLock()
	defer p.l.RUnlock()

	if !p.Koanf.Exists(key) {
		return fallback
	}

	switch v := p.Koanf.Get(key).(type) {
	case string:
		// this type usually comes from user input
		dec, err := bytesize.Parse(v)
		if err != nil {
			p.logger.With(
				slog.String("key", key),
				slog.String("raw_value", v),
			).Warn(fmt.Sprintf("error parsing byte size value, using fallback of %s", fallback))
			return fallback
		}
		return dec
	case float64:
		// this type comes from json.Unmarshal
		return bytesize.ByteSize(v)
	case bytesize.ByteSize:
		return v
	default:
		p.logger.With(
			slog.String("key", key),
			slog.Any("raw_value", v),
		).Warn(fmt.Sprintf("error parsing byte size value, using fallback of %s", fallback))
		return fallback
	}
}

func (p *Provider) GetF(key string, fallback interface{}) (val interface{}) {
	p.l.RLock()
	defer p.l.RUnlock()

	if !p.Exists(key) {
		return fallback
	}

	return p.Get(key)
}

//func (p *Provider) CORS(prefix string, defaults cors.Options) (cors.Options, bool) {
//	if len(prefix) > 0 {
//		prefix = strings.TrimRight(prefix, ".") + "."
//	}
//
//	return cors.Options{
//		AllowedOrigins:     p.StringsF(prefix+"cors.allowed_origins", defaults.AllowedOrigins),
//		AllowedMethods:     p.StringsF(prefix+"cors.allowed_methods", defaults.AllowedMethods),
//		AllowedHeaders:     p.StringsF(prefix+"cors.allowed_headers", defaults.AllowedHeaders),
//		ExposedHeaders:     p.StringsF(prefix+"cors.exposed_headers", defaults.ExposedHeaders),
//		AllowCredentials:   p.BoolF(prefix+"cors.allow_credentials", defaults.AllowCredentials),
//		OptionsPassthrough: p.BoolF(prefix+"cors.options_passthrough", defaults.OptionsPassthrough),
//		MaxAge:             p.IntF(prefix+"cors.max_age", defaults.MaxAge),
//		Debug:              p.BoolF(prefix+"cors.debug", defaults.Debug),
//	}, p.Bool(prefix + "cors.enabled")
//}

func (p *Provider) RequestURIF(path string, fallback *url.URL) *url.URL {
	p.l.RLock()
	defer p.l.RUnlock()

	switch t := p.Get(path).(type) {
	case *url.URL:
		return t
	case url.URL:
		return &t
	case string:
		if parsed, err := url.ParseRequestURI(t); err == nil {
			return parsed
		}
	}

	return fallback
}

func (p *Provider) URIF(path string, fallback *url.URL) *url.URL {
	p.l.RLock()
	defer p.l.RUnlock()

	switch t := p.Get(path).(type) {
	case *url.URL:
		return t
	case url.URL:
		return &t
	case string:
		if parsed, err := url.Parse(t); err == nil {
			return parsed
		}
	}

	return fallback
}

func (p *Provider) Unmarshal(path string, o interface{}) error {
	return p.UnmarshalWithConf(path, o, koanf.UnmarshalConf{})
}

func (p *Provider) UnmarshalWithConf(path string, o interface{}, c koanf.UnmarshalConf) error {
	if c.DecoderConfig == nil {
		c.DecoderConfig = &mapstructure.DecoderConfig{
			DecodeHook:       p.decodeHookFunc,
			Metadata:         nil,
			Result:           o,
			WeaklyTypedInput: true,
		}
	}
	if c.Tag == "" {
		c.Tag = "json"
	}

	return p.Koanf.UnmarshalWithConf(path, o, c)
}

func getFiles(paths []string, filter func(path string) bool) []string {
	var files []string
	for _, path := range paths {
		if fileutil.IsDir(path) {
			_ = filepath.Walk(path, func(p string, info fs.FileInfo, err error) error {
				if !info.IsDir() && filter(p) {
					files = append(files, p)
				}
				return nil
			})
		} else if filter(path) {
			files = append(files, path)
		}
	}

	return files
}

func isConfigFile(path string) bool {
	return strings.HasSuffix(strings.ToLower(filepath.Base(path)), ".json") ||
		strings.HasSuffix(strings.ToLower(filepath.Base(path)), ".yaml") ||
		strings.HasSuffix(strings.ToLower(filepath.Base(path)), ".yml") ||
		strings.HasSuffix(strings.ToLower(filepath.Base(path)), ".toml")
}
