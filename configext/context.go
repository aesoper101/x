package configext

import "context"

type contextKey int

const configContextKey contextKey = iota + 1

func ContextWithConfigOptions(ctx context.Context, opts ...ProviderOptions) context.Context {
	return context.WithValue(ctx, configContextKey, opts)
}

func ConfigOptionsFromContext(ctx context.Context) []ProviderOptions {
	opts, ok := ctx.Value(configContextKey).([]ProviderOptions)
	if !ok {
		return []ProviderOptions{}
	}
	return opts
}
