package jsonschemaext

import (
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRead(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat()
	compiler.AssertContent()
	compiler.AssertVocabs()

	schema, err := compiler.Compile("./stub/config.schema.json")
	require.Nil(t, err)
	require.NotNil(t, schema)
}
