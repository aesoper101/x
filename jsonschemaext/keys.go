package jsonschemaext

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"math/big"
	"slices"
	"sort"
	"strings"
)

type (
	byName       []Path
	PathEnhancer interface {
		EnhancePath(Path) map[string]interface{}
	}
	TypeHint int
)

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name < s[j].Name }

type Enum = jsonschema.Enum
type Regexp = jsonschema.Regexp
type Format = jsonschema.Format

const (
	String TypeHint = iota + 1
	Float
	Int
	Bool
	JSON
	Nil

	BoolSlice
	StringSlice
	IntSlice
	FloatSlice
)

// Path represents a JSON Schema Path.
type Path struct {
	// Title of the path.
	Title string

	// Description of the path.
	Description string

	// Examples of the path.
	Examples []interface{}

	// Name is the JSON path name.
	Name string

	// Default is the default value of that path.
	Default interface{}

	// Type is a prototype (e.g. float64(0)) of the path type.
	//Type interface{}

	TypeHint

	// Format is the format of the path if defined
	Format *Format

	// Pattern is the pattern of the path if defined
	Pattern Regexp

	// Enum are the allowed enum values
	Enum *Enum

	// first element in slice is constant value. note: slice is used to capture nil constant.
	Constant interface{}

	// ReadOnly is whether the value is readonly
	ReadOnly bool

	// -1 if not specified
	MinLength *int
	MaxLength *int

	// Required if set indicates this field is required.
	Required bool

	Minimum *big.Rat
	Maximum *big.Rat

	MultipleOf *big.Rat

	CustomProperties map[string]interface{}
}

// ListPathsBytes works like ListPathsWithRecursion but prepares the JSON Schema itself.
func ListPathsBytes(ctx context.Context, raw json.RawMessage, maxRecursion int16) ([]Path, error) {
	compiler := jsonschema.NewCompiler()
	compiler.AssertContent()
	compiler.AssertVocabs()
	compiler.AssertFormat()

	id := fmt.Sprintf("%x.json", sha256.Sum256(raw))
	if err := compiler.AddResource(id, bytes.NewReader(raw)); err != nil {
		return nil, err
	}
	compiler.AssertContent()
	compiler.AssertVocabs()
	compiler.AssertFormat()
	return runPathsFromCompiler(ctx, id, compiler, maxRecursion, false)
}

// ListPathsWithRecursion will follow circular references until maxRecursion is reached, without
// returning an error.
func ListPathsWithRecursion(ctx context.Context, ref string, compiler *jsonschema.Compiler, maxRecursion uint8) (
	[]Path,
	error,
) {
	return runPathsFromCompiler(ctx, ref, compiler, int16(maxRecursion), false)
}

// ListPaths lists all paths of a JSON Schema. Will return an error
// if circular references are found.
func ListPaths(ctx context.Context, ref string, compiler *jsonschema.Compiler) ([]Path, error) {
	return runPathsFromCompiler(ctx, ref, compiler, -1, false)
}

// ListPathsWithArraysIncluded lists all paths of a JSON Schema. Will return an error
// if circular references are found.
// Includes arrays with `#`.
func ListPathsWithArraysIncluded(ctx context.Context, ref string, compiler *jsonschema.Compiler) ([]Path, error) {
	return runPathsFromCompiler(ctx, ref, compiler, -1, true)
}

// ListPathsWithInitializedSchema loads the paths from the schema without compiling it.
func ListPathsWithInitializedSchema(schema *jsonschema.Schema) ([]Path, error) {
	return runPaths(schema, -1, false)
}

// ListPathsWithInitializedSchemaAndArraysIncluded loads the paths from the schema without compiling it.
func ListPathsWithInitializedSchemaAndArraysIncluded(schema *jsonschema.Schema) ([]Path, error) {
	return runPaths(schema, -1, true)
}

func runPathsFromCompiler(
	ctx context.Context,
	ref string,
	compiler *jsonschema.Compiler,
	maxRecursion int16,
	includeArrays bool,
) ([]Path, error) {
	if compiler == nil {
		compiler = jsonschema.NewCompiler()
	}

	compiler.AssertContent()
	compiler.AssertVocabs()
	compiler.AssertFormat()

	schema, err := compiler.Compile(ref)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return runPaths(schema, maxRecursion, includeArrays)
}

func runPaths(schema *jsonschema.Schema, maxRecursion int16, includeArrays bool) ([]Path, error) {
	pointers := map[string]bool{}
	paths, err := listPaths(schema, nil, nil, pointers, 0, maxRecursion, includeArrays)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sort.Stable(paths)
	return makeUnique(paths)
}

func makeUnique(in byName) (byName, error) {
	cache := make(map[string]Path)
	for _, p := range in {
		vc, ok := cache[p.Name]
		if !ok {
			cache[p.Name] = p
			continue
		}

		//if fmt.Sprintf("%T", p.Type) != fmt.Sprintf("%T", p.Type) {
		//	return nil, errors.Errorf(
		//		"multiple types %+v are not supported for path: %s",
		//		[]interface{}{p.Type, vc.Type},
		//		p.Name,
		//	)
		//}

		if vc.Default == nil {
			cache[p.Name] = p
		}
	}

	k := 0
	out := make([]Path, len(cache))
	for _, v := range cache {
		out[k] = v
		k++
	}

	paths := byName(out)
	sort.Sort(paths)
	return paths, nil
}

func appendPointer(in map[string]bool, pointer *jsonschema.Schema) map[string]bool {
	out := make(map[string]bool)
	for k, v := range in {
		out[k] = v
	}
	out[fmt.Sprintf("%p", pointer)] = true
	return out
}

func listPaths(
	schema *jsonschema.Schema,
	parent *jsonschema.Schema,
	parents []string,
	pointers map[string]bool,
	currentRecursion int16,
	maxRecursion int16,
	includeArrays bool,
) (byName, error) {
	//var pathTypeHint TypeHint
	var paths []Path
	pathTypeHint := String
	_, isCircular := pointers[fmt.Sprintf("%p", schema)]

	if schema.Const != nil {
		v := *schema.Const
		switch v.(type) {
		case float64, json.Number:
			pathTypeHint = Float
		case int8, int16, int, int64:
			pathTypeHint = Int
		case string:
			pathTypeHint = String
		case bool:
			pathTypeHint = Bool
		default:
			pathTypeHint = JSON
		}
	} else if schema.Types != nil {
		if len(schema.Types.ToStrings()) == 1 {
			switch schema.Types.ToStrings()[0] {
			case "null":
				pathTypeHint = Nil
			case "boolean":
				pathTypeHint = Bool
			case "number":
				pathTypeHint = Float
			case "integer":
				pathTypeHint = Int
			case "string":
				pathTypeHint = String
			case "array":
				if schema.Items != nil {
					var itemSchemas []*jsonschema.Schema
					switch t := schema.Items.(type) {
					case []*jsonschema.Schema:
						itemSchemas = t
					case *jsonschema.Schema:
						itemSchemas = []*jsonschema.Schema{t}
					}
					var types []string
					for _, is := range itemSchemas {
						if is.Types != nil {
							types = append(types, is.Types.ToStrings()...)
						}
						if is.Ref != nil && is.Ref.Types != nil {
							types = append(types, is.Ref.Types.ToStrings()...)
						}
					}
					types = slices.Compact(types)
					if len(types) == 1 {
						switch types[0] {
						case "boolean":
							pathTypeHint = BoolSlice
						case "number":
							pathTypeHint = FloatSlice
						case "integer":
							pathTypeHint = IntSlice
						case "string":
							pathTypeHint = StringSlice
						default:
							pathTypeHint = JSON
						}
					}
				}
			case "object":
				pathTypeHint = JSON
			}
		}
	}

	var def interface{} = schema.Default
	if v, ok := def.(json.Number); ok {
		def, _ = v.Float64()
	}

	if len(parents) > 0 {
		name := parents[len(parents)-1]
		var required bool
		if parent != nil {
			for _, r := range parent.Required {
				if r == name {
					required = true
					break
				}
			}
		}

		path := Path{
			Name:        strings.Join(parents, "."),
			Default:     def,
			TypeHint:    pathTypeHint,
			Format:      schema.Format,
			Pattern:     schema.Pattern,
			Enum:        schema.Enum,
			Constant:    schema.Const,
			MinLength:   schema.MinLength,
			MaxLength:   schema.MaxLength,
			Minimum:     schema.Minimum,
			Maximum:     schema.Maximum,
			MultipleOf:  schema.MultipleOf,
			ReadOnly:    schema.ReadOnly,
			Title:       schema.Title,
			Description: schema.Description,
			Examples:    schema.Examples,
			Required:    required,
		}

		for _, e := range schema.Extensions {
			if enhancer, ok := e.(PathEnhancer); ok {
				path.CustomProperties = enhancer.EnhancePath(path)
			}
		}
		paths = append(paths, path)
	}

	if isCircular {
		if maxRecursion == -1 {
			return nil, errors.Errorf("detected circular dependency in schema path: %s", strings.Join(parents, "."))
		} else if currentRecursion > maxRecursion {
			return paths, nil
		}
		currentRecursion++
	}

	if schema.Ref != nil {
		path, err := listPaths(
			schema.Ref,
			schema,
			parents,
			appendPointer(pointers, schema),
			currentRecursion,
			maxRecursion,
			includeArrays,
		)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	if schema.Not != nil {
		path, err := listPaths(
			schema.Not,
			schema,
			parents,
			appendPointer(pointers, schema),
			currentRecursion,
			maxRecursion,
			includeArrays,
		)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	if schema.If != nil {
		path, err := listPaths(
			schema.If,
			schema,
			parents,
			appendPointer(pointers, schema),
			currentRecursion,
			maxRecursion,
			includeArrays,
		)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	if schema.Then != nil {
		path, err := listPaths(
			schema.Then,
			schema,
			parents,
			appendPointer(pointers, schema),
			currentRecursion,
			maxRecursion,
			includeArrays,
		)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	if schema.Else != nil {
		path, err := listPaths(
			schema.Else,
			schema,
			parents,
			appendPointer(pointers, schema),
			currentRecursion,
			maxRecursion,
			includeArrays,
		)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	for _, sub := range schema.AllOf {
		path, err := listPaths(
			sub,
			schema,
			parents,
			appendPointer(pointers, schema),
			currentRecursion,
			maxRecursion,
			includeArrays,
		)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	for _, sub := range schema.AnyOf {
		path, err := listPaths(
			sub,
			schema,
			parents,
			appendPointer(pointers, schema),
			currentRecursion,
			maxRecursion,
			includeArrays,
		)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	for _, sub := range schema.OneOf {
		path, err := listPaths(
			sub,
			schema,
			parents,
			appendPointer(pointers, schema),
			currentRecursion,
			maxRecursion,
			includeArrays,
		)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	for name, sub := range schema.Properties {
		path, err := listPaths(
			sub,
			schema,
			append(parents, name),
			appendPointer(pointers, schema),
			currentRecursion,
			maxRecursion,
			includeArrays,
		)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	if schema.Items != nil && includeArrays {
		switch t := schema.Items.(type) {
		case []*jsonschema.Schema:
			for _, sub := range t {
				path, err := listPaths(
					sub,
					schema,
					append(parents, "#"),
					appendPointer(pointers, schema),
					currentRecursion,
					maxRecursion,
					includeArrays,
				)
				if err != nil {
					return nil, err
				}
				paths = append(paths, path...)
			}
		case *jsonschema.Schema:
			path, err := listPaths(
				t,
				schema,
				append(parents, "#"),
				appendPointer(pointers, schema),
				currentRecursion,
				maxRecursion,
				includeArrays,
			)
			if err != nil {
				return nil, err
			}
			paths = append(paths, path...)
		}
	}

	return paths, nil
}
