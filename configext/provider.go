package configext

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aesoper101/x/watcherext"
	"github.com/inhies/go-bytesize"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/spf13/pflag"
	"log/slog"
	"net/url"
	"reflect"
	"sync"
	"time"
)

type tuple struct {
	Key   string
	Value interface{}
}

type Provider struct {
	l sync.RWMutex
	*koanf.Koanf
	immutables, exceptImmutables []string

	schema                   []byte
	flags                    *pflag.FlagSet
	validator                *jsonschema.Schema
	onChanges                []func(watcherext.Event, error)
	onValidationError        func(k *koanf.Koanf, err error)
	excludeFieldsFromTracing []string

	forcedValues []tuple
	baseValues   []tuple
	files        []string

	skipValidation    bool
	disableEnvLoading bool

	logger *slog.Logger

	providers     []koanf.Provider
	userProviders []koanf.Provider
}

const (
	FlagConfig = "config"
	Delimiter  = "."
)

// RegisterConfigFlag registers the "--config" flag on pflag.FlagSet.
func RegisterConfigFlag(flags *pflag.FlagSet, fallback []string) {
	flags.StringSliceP(FlagConfig, "c", fallback, "Config files to load, overwriting in the order specified.")
}

// New creates a new provider instance or errors.
// Configuration values are loaded in the following order:
//
// 1. Defaults from the JSON Schema
// 2. Config files (yaml, yml, toml, json)
// 3. Command line flags
// 4. Environment variables
//
// There will also be file-watchers started for all config files. To cancel the
// watchers, cancel the context.
func New(ctx context.Context, schema []byte, modifiers ...OptionModifier) (*Provider, error) {
	validator, err := getSchema(ctx, schema)
	if err != nil {
		return nil, err
	}

	l := slog.Default()

	p := &Provider{
		schema:                   schema,
		validator:                validator,
		onValidationError:        func(k *koanf.Koanf, err error) {},
		excludeFieldsFromTracing: []string{"dsn", "secret", "password", "key"},
		logger:                   l,
		Koanf:                    koanf.NewWithConf(koanf.Conf{Delim: Delimiter, StrictMerge: true}),
	}

	for _, m := range modifiers {
		m(p)
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

func (p *Provider) SkipValidation() bool {
	return p.skipValidation
}

func (p *Provider) createProviders(ctx context.Context) (providers []koanf.Provider, err error) {
	defaultsProvider, err := NewKoanfSchemaDefaults(p.schema, p.validator)
	if err != nil {
		return nil, err
	}
	providers = append(providers, defaultsProvider)

	// Workaround for https://github.com/knadh/koanf/pull/47
	for _, t := range p.baseValues {
		providers = append(providers, NewKoanfConfmap([]tuple{t}))
	}

	paths := p.files
	if p.flags != nil {
		p, _ := p.flags.GetStringSlice(FlagConfig)
		paths = append(paths, p...)
	}

	p.logger.Debug("Adding config file .", slog.Any("files", paths))

	c := make(watcherext.EventChannel)
	defer close(c)

	go p.watchForFileChanges(ctx, c)

	for _, path := range paths {
		fp, err := NewKoanfFile(path)
		if err != nil {
			return nil, err
		}

		if _, err := fp.WatchChannel(ctx, c); err != nil {
			return nil, err
		}

		providers = append(providers, fp)
	}

	providers = append(providers, p.userProviders...)

	if p.flags != nil {
		pp, err := NewPFlagProvider(p.schema, p.validator, p.flags, p.Koanf)
		if err != nil {
			return nil, err
		}
		providers = append(providers, pp)
	}

	if !p.disableEnvLoading {
		envProvider, err := NewKoanfEnv("", p.schema, p.validator)
		if err != nil {
			return nil, err
		}
		providers = append(providers, envProvider)
	}

	// Workaround for https://github.com/knadh/koanf/pull/47
	for _, t := range p.forcedValues {
		providers = append(providers, NewKoanfConfmap([]tuple{t}))
	}

	return providers, nil
}

func (p *Provider) replaceKoanf(k *koanf.Koanf) {
	p.Koanf = k
}

func (p *Provider) validate(k *koanf.Koanf) error {
	if p.skipValidation {
		return nil
	}

	out, err := k.Marshal(json.Parser())
	if err != nil {
		return errors.WithStack(err)
	}

	inst, err := jsonschema.UnmarshalJSON(bytes.NewReader(out))
	if err != nil {
		return errors.WithStack(err)
	}

	if err := p.validator.Validate(inst); err != nil {
		p.onValidationError(k, err)
		return err
	}

	return nil
}

// newKoanf creates a new koanf instance with all the updated config
//
// This is unfortunately required due to several limitations / bugs in koanf:
//
// - https://github.com/knadh/koanf/issues/77
// - https://github.com/knadh/koanf/pull/47
func (p *Provider) newKoanf() (_ *koanf.Koanf, err error) {
	k := koanf.New(Delimiter)

	for _, provider := range p.providers {
		// posflag.Posflag requires access to Koanf instance so we recreate the provider here which is a workaround
		// for posflag.Provider's API.
		if _, ok := provider.(*posflag.Posflag); ok {
			provider = posflag.Provider(p.flags, ".", k)
		}

		var opts []koanf.Option
		if _, ok := provider.(*Env); ok {
			opts = append(opts, koanf.WithMergeFunc(MergeAllTypes))
		}

		if err := k.Load(provider, nil, opts...); err != nil {
			return nil, err
		}
	}

	if err := p.validate(k); err != nil {
		return nil, err
	}

	return k, nil
}

func (p *Provider) runOnChanges(e watcherext.Event, err error) {
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

func (p *Provider) reload(e watcherext.Event) {
	p.l.Lock()

	var err error
	defer func() {
		// we first want to unlock and then runOnChanges, so that the callbacks can actually use the Provider
		p.l.Unlock()
		p.runOnChanges(e, err)
	}()

	nk, err := p.newKoanf()
	if err != nil {
		return // unlocks & runs changes in defer
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

	// unlocks & runs changes in defer
}

func (p *Provider) watchForFileChanges(ctx context.Context, c watcherext.EventChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case e, ok := <-c:
			if !ok {
				return
			}
			p.reload(e)
		}
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

func (p *Provider) StringF(key string, fallback string) string {
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
			p.logger.Warn(
				fmt.Sprintf("error parsing byte size value, using fallback of %s", fallback), slog.Any("key", key),
				slog.Any("raw_value", v),
			)
			return fallback
		}
		return dec
	case float64:
		// this type comes from json.Unmarshal
		return bytesize.ByteSize(v)
	case bytesize.ByteSize:
		return v
	default:
		p.logger.Error(
			fmt.Sprintf(
				"error converting byte size value because of unknown type, using fallback of %s",
				fallback,
			),
			slog.Any("key", key),
			slog.Any("raw_type", fmt.Sprintf("%T", v)),
			slog.Any("raw_value", fmt.Sprintf("%+v", v)),
		)
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
