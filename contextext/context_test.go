package contextext

import (
	"context"
	"reflect"
	"testing"
)

func TestNewContext(t *testing.T) {
	type args[K comparable, V any] struct {
		ctx   context.Context
		key   K
		value V
	}
	type Key struct{}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want context.Context
	}
	tests := []testCase[Key, int]{
		{
			name: "nil",
			want: context.WithValue(context.Background(), Key{}, 11),
			args: args[Key, int]{
				ctx:   context.Background(),
				key:   Key{},
				value: 11,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gotCtx := NewContext(tt.args.ctx, tt.args.key, tt.args.value)
				if !reflect.DeepEqual(gotCtx, tt.want) {
					t.Errorf("NewContext() = %v, want %v", gotCtx, tt.want)
				}

				if v, _ := FromContext[Key, int](gotCtx, tt.args.key); !reflect.DeepEqual(v, tt.args.value) {
					t.Errorf("FromContext() = %v, want %v", v, tt.want)
				}
			},
		)
	}
}
