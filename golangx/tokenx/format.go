package tokenx

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

func FormatGoCode(filename string, content []byte) ([]byte, error) {
	return formatGoCode(filename, content)
}

func formatGoCode(filename string, content []byte) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("can not parse ast for file: %s, err: %v", filename, err)
	}

	return formatGoCodeByAst(fset, f)
}

func formatGoCodeByAst(fset *token.FileSet, f *ast.File) ([]byte, error) {
	var output []byte
	buffer := bytes.NewBuffer(output)
	err := format.Node(buffer, fset, f)
	if err != nil {
		return nil, fmt.Errorf("can not format file: %s, err: %v", f.Name.Name, err)
	}

	return buffer.Bytes(), nil
}
