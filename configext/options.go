package configext

import (
	"errors"
	"fmt"
	"github.com/aesoper101/x/watcherext"
	"github.com/knadh/koanf/v2"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/spf13/pflag"
	"log/slog"
)

type (
	OptionModifier func(p *Provider)
)

//func WithContext(ctx context.Context) OptionModifier {
//	return func(p *Provider) {
//		p.ctx = ctx
//	}
//}

func WithConfigFiles(files ...string) OptionModifier {
	return func(p *Provider) {
		p.files = append(p.files, files...)
	}
}

func WithImmutables(immutables ...string) OptionModifier {
	return func(p *Provider) {
		p.immutables = append(p.immutables, immutables...)
	}
}

func WithExceptImmutables(exceptImmutables ...string) OptionModifier {
	return func(p *Provider) {
		p.exceptImmutables = append(p.exceptImmutables, exceptImmutables...)
	}
}

func WithFlags(flags *pflag.FlagSet) OptionModifier {
	return func(p *Provider) {
		p.flags = flags
	}
}

func WithLogger(l *slog.Logger) OptionModifier {
	return func(p *Provider) {
		p.logger = l
	}
}

func SkipValidation() OptionModifier {
	return func(p *Provider) {
		p.skipValidation = true
	}
}

func DisableEnvLoading() OptionModifier {
	return func(p *Provider) {
		p.disableEnvLoading = true
	}
}

func WithValue(key string, value interface{}) OptionModifier {
	return func(p *Provider) {
		p.forcedValues = append(p.forcedValues, tuple{Key: key, Value: value})
	}
}

func WithValues(values map[string]interface{}) OptionModifier {
	return func(p *Provider) {
		for key, value := range values {
			p.forcedValues = append(p.forcedValues, tuple{Key: key, Value: value})
		}
	}
}

func WithBaseValues(values map[string]interface{}) OptionModifier {
	return func(p *Provider) {
		for key, value := range values {
			p.baseValues = append(p.baseValues, tuple{Key: key, Value: value})
		}
	}
}

func WithUserProviders(providers ...koanf.Provider) OptionModifier {
	return func(p *Provider) {
		p.userProviders = providers
	}
}

func OmitKeysFromTracing(keys ...string) OptionModifier {
	return func(p *Provider) {
		p.excludeFieldsFromTracing = keys
	}
}

func AttachWatcher(watcher func(event watcherext.Event, err error)) OptionModifier {
	return func(p *Provider) {
		p.onChanges = append(p.onChanges, watcher)
	}
}

func WithLoggerWatcher(l *slog.Logger) OptionModifier {
	return AttachWatcher(LogWatcher(l))
}

func LogWatcher(l *slog.Logger) func(e watcherext.Event, err error) {
	return func(e watcherext.Event, err error) {
		l.Info(
			"A change to a configuration file was detected.", slog.Any("file", e.Source()), slog.Any(
				"event_type",
				fmt.Sprintf("%T", e),
			),
		)

		if et := new(jsonschema.ValidationError); errors.As(err, &et) {
			l.Error(
				"The changed configuration is invalid and could not be loaded. "+
					"Rolling back to the last working configuration revision. "+
					"Please address the validation errors before restarting the process.",
				"event",
				fmt.Sprintf("%#v", et),
			)
		} else if et := new(ImmutableError); errors.As(err, &et) {
			l.Error(
				"A configuration value marked as immutable has changed. "+
					"Rolling back to the last working configuration revision. "+
					"To reload the values please restart the process.", slog.Any("key", et.Key), slog.Any(
					"old_value",
					fmt.Sprintf("%v", et.From),
				), slog.Any("new_value", fmt.Sprintf("%v", et.To)),
			)
		} else if err != nil {
			l.Error(
				"An error occurred while watching config file.", slog.Any("error", err), slog.Any(
					"file",
					e.Source(),
				),
			)
		} else {
			l.Info(
				"Configuration change processed successfully.",
				slog.Any("file", e.Source()),
				slog.Any("event_type", fmt.Sprintf("%T", e)),
			)
		}
	}
}
