package configext

import (
	"bytes"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/santhosh-tekuri/jsonschema/v6"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

func newCompiler(schemaByte []byte) (string, *jsonschema.Compiler, error) {
	id := gjson.GetBytes(schemaByte, "$id").String()
	if id == "" {
		id = fmt.Sprintf("%s.json", uuid.Must(uuid.NewV4()).String())
	}

	schema, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaByte))
	if err != nil {
		return "", nil, err
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(id, schema); err != nil {
		return "", nil, errors.WithStack(err)
	}

	compiler.AssertContent()
	compiler.AssertFormat()
	compiler.AssertVocabs()

	return id, compiler, nil
}
