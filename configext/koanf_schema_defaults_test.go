package configext

import (
	"bytes"
	"context"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestKoanfSchemaDefaults(t *testing.T) {
	schemaPath := path.Join("stub", "domain-aliases", "config.schema.json")

	rawSchema, err := os.ReadFile(schemaPath)
	require.NoError(t, err)

	schemaDoc, err := jsonschema.UnmarshalJSON(bytes.NewReader(rawSchema))
	require.Nil(t, err)
	require.NotNil(t, schemaDoc)

	c := jsonschema.NewCompiler()
	require.NoError(t, c.AddResource(schemaPath, schemaDoc))

	schema, err := c.Compile(schemaPath)
	require.NoError(t, err)

	k, err := newKoanf(context.Background(), schemaPath, nil)
	require.NoError(t, err)

	def, err := NewKoanfSchemaDefaults(rawSchema, schema)
	require.NoError(t, err)

	require.NoError(t, k.Load(def, nil))
}
