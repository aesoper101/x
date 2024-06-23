package configext

import (
	"context"
	"log/slog"

	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

type ProviderOptions func(p *Provider)

func WithLogger(l *slog.Logger) ProviderOptions {
	return func(p *Provider) {
		p.logger = l
	}
}

func WithConfigFiles(files ...string) ProviderOptions {
	return func(p *Provider) {
		p.files = append(p.files, files...)
	}
}

func WithFlags(flags *pflag.FlagSet) ProviderOptions {
	return func(p *Provider) {
		p.flags = flags
	}
}

func WithUserProviders(providers ...koanf.Provider) ProviderOptions {
	return func(p *Provider) {
		p.userProviders = providers
	}
}

func WithValue(key string, value interface{}) ProviderOptions {
	return func(p *Provider) {
		p.forcedValues = append(p.forcedValues, tuple{Key: key, Value: value})
	}
}

func WithValues(values map[string]interface{}) ProviderOptions {
	return func(p *Provider) {
		for key, value := range values {
			p.forcedValues = append(p.forcedValues, tuple{Key: key, Value: value})
		}
	}
}

func WithBaseValues(values map[string]interface{}) ProviderOptions {
	return func(p *Provider) {
		for key, value := range values {
			p.baseValues = append(p.baseValues, tuple{Key: key, Value: value})
		}
	}
}

func WithImmutables(immutables ...string) ProviderOptions {
	return func(p *Provider) {
		p.immutables = append(p.immutables, immutables...)
	}
}

func WithExceptImmutables(exceptImmutables ...string) ProviderOptions {
	return func(p *Provider) {
		p.exceptImmutables = append(p.exceptImmutables, exceptImmutables...)
	}
}

func WithContext(ctx context.Context) ProviderOptions {
	return func(p *Provider) {
		for _, o := range ConfigOptionsFromContext(ctx) {
			o(p)
		}
	}
}

func WithChangeNotifier(notifier func(event interface{}, err error)) ProviderOptions {
	return func(p *Provider) {
		p.onChanges = append(p.onChanges, notifier)
	}
}

// EnableEnvLoading  enables env loading.
// Notes: cant be used with DisabledEnvLoading at the same time
func EnableEnvLoading(prefix string) ProviderOptions {
	return func(p *Provider) {
		p.enableEnvLoading = true
		p.envPrefix = prefix
	}
}

// DisabledEnvLoading disables env loading.
// Notes: cant be used with EnableEnvLoading at the same time
func DisabledEnvLoading() ProviderOptions {
	return func(p *Provider) {
		p.enableEnvLoading = false
	}
}

func WithDecodeHookFunc(decodeHookFunc mapstructure.DecodeHookFunc) ProviderOptions {
	return func(p *Provider) {
		p.decodeHookFunc = decodeHookFunc
	}
}
