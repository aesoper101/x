package configext

import (
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestPFlagProvider(t *testing.T) {
	const schema = `
{
  "type": "object",
  "properties": {
	"foo": {
	  "type": "string"
	}
  }
}
`

	bs, err := jsonschema.UnmarshalJSON(strings.NewReader(schema))
	require.NoError(t, err)

	c := jsonschema.NewCompiler()
	err = c.AddResource("schema.json", bs)
	require.Nil(t, err)

	s, err := c.Compile("schema.json")
	require.Nil(t, err)
	require.NotNil(t, s)
	t.Run(
		"only parses known flags", func(t *testing.T) {
			flags := pflag.NewFlagSet("", pflag.ContinueOnError)
			flags.String("foo", "", "")
			flags.String("bar", "", "")
			require.NoError(t, flags.Parse([]string{"--foo", "x", "--bar", "y"}))

			p, err := NewPFlagProvider([]byte(schema), s, flags, nil)
			require.NoError(t, err)

			values, err := p.Read()
			require.NoError(t, err)
			assert.Equal(
				t, map[string]interface{}{
					"foo": "x",
				}, values,
			)
		},
	)
}
