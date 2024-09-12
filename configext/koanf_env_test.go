package configext

import (
	_ "embed"
	"github.com/dgraph-io/ristretto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

//go:embed stub/watch/config.schema.json
var _xConfigSchema []byte

func TestNewKoanfEnvCache(t *testing.T) {
	ref, compiler, err := newCompiler(_xConfigSchema)
	require.NoError(t, err)
	require.NotNil(t, compiler)
	require.NotEmpty(t, ref)

	schema, err := compiler.Compile(ref)
	require.NoError(t, err)

	c := *schemaPathCacheConfig
	c.Metrics = true
	schemaPathCache, _ = ristretto.NewCache(&c)
	_, _ = NewKoanfEnv("", _xConfigSchema, schema)
	_, _ = NewKoanfEnv("", _xConfigSchema, schema)
	assert.EqualValues(t, 1, schemaPathCache.Metrics.Hits())
}
