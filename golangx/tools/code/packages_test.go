package code

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackages_LoadAll(t *testing.T) {
	p := NewPackages()
	pkgs := p.LoadAll("github.com/aesoper101/x/golangx/tools/code")
	require.NotNil(t, pkgs)

	//var mode = packages.NeedName |
	//	packages.NeedFiles |
	//	packages.NeedSyntax |
	//	packages.NeedTypes |
	//	packages.NeedDeps |
	//	packages.NeedImports |
	//	packages.NeedTypesInfo |
	//	packages.NeedModule |
	//	packages.NeedTypesSizes
	//
	//load, err := packages.Load(&packages.Config{
	//	Mode: mode,
	//})
	//require.Nil(t, err)
	//require.NotNil(t, load)

	for _, pkg := range pkgs {
		fmt.Printf("Package: %s\n", pkg.ID)
		fmt.Printf("Name: %s\n", pkg.Name)
		fmt.Printf("PkgPath: %s\n", pkg.PkgPath)
		fmt.Printf("GoFiles: %v\n", pkg.GoFiles)
		fmt.Printf("CompiledGoFiles: %v\n", pkg.CompiledGoFiles)
		fmt.Printf("Syntax: %v\n", pkg.Syntax)
		fmt.Printf("Types: %v\n", pkg.Types)
		fmt.Printf("TypesInfo: %v\n", pkg.TypesInfo)
		fmt.Printf("TypesSizes: %v\n", pkg.TypesSizes)
		fmt.Printf("Imports: %v\n", pkg.Imports)
		fmt.Printf("Deps: %v\n", pkg.OtherFiles)
		fmt.Printf("Module: %v\n", pkg.Module)
	}
}
