package transportx

import (
	"context"
	"github.com/aesoper101/x/contextext"
)

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s AppInfo) context.Context {
	return contextext.NewContext[appKey, AppInfo](ctx, appKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s AppInfo, ok bool) {
	return contextext.FromContext[appKey, AppInfo](ctx, appKey{})
}
