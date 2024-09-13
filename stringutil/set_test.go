package stringutil

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"reflect"
	"testing"
)

func TestStringSet_MarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		set     StringSet
		want    interface{}
		wantErr bool
	}{
		{
			name: "empty set",
			set:  StringSet{},
			want: []string{},
		}, {
			name: "non-empty set",
			set:  NewStringSet("a", "b"),
			want: []string{"a", "b"},
		},
		{
			name: "duplicate elements",
			set:  NewStringSet("a", "a", "b"),
			want: []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := yaml.Marshal(tt.set)
				if (err != nil) != tt.wantErr {
					t.Errorf("MarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				var sl []string
				if err := yaml.Unmarshal(got, &sl); err != nil {
					t.Errorf("UnmarshalYAML() error = %v", err)
				}

				if !reflect.DeepEqual(sl, tt.want) {
					t.Errorf("MarshalYAML() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestStringSet_UnmarshalYAML(t *testing.T) {
	schema := []byte(`---
set:
 - a
 - b
`)

	type TestSet struct {
		Set StringSet `yaml:"set"`
	}

	var data TestSet
	err := yaml.Unmarshal(schema, &data)
	require.NoError(t, err)
}
