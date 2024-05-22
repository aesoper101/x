package code

import (
	"go/build"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// PkgAndType take a string in the form github.com/package/blah.Type and split it into package and type
// From: https://github.com/99designs/gqlgen/blob/d5c9f896419142f6378639b6eec93584fbf829ed/internal/code/util.go#L11C1-L20C1
func PkgAndType(name string) (string, string) {
	parts := strings.Split(name, ".")
	if len(parts) == 1 {
		return "", name
	}

	return strings.Join(parts[:len(parts)-1], "."), parts[len(parts)-1]
}

var modsRegex = regexp.MustCompile(`^(\*|\[\])*`)

// NormalizeVendor takes a qualified package path and turns it into normal one.
// eg .
// github.com/foo/vendor/github.com/aesoper101/xxx/graphql becomes
// github.com/aesoper101/xxx/graphql
//
// From: https://github.com/99designs/gqlgen/blob/d5c9f896419142f6378639b6eec93584fbf829ed/internal/code/util.go#L27
func NormalizeVendor(pkg string) string {
	modifiers := modsRegex.FindAllString(pkg, 1)[0]
	pkg = strings.TrimPrefix(pkg, modifiers)
	parts := strings.Split(pkg, "/vendor/")
	return modifiers + parts[len(parts)-1]
}

// QualifyPackagePath takes an import and fully qualifies it with a vendor dir, if one is required.
// eg .
// github.com/aesoper101/xxx/graphql becomes
// github.com/foo/vendor/github.com/aesoper101/xxx/graphql
//
// x/tools/packages only supports 'qualified package paths' so this will need to be done prior to calling it
// See https://github.com/golang/go/issues/30289
// From https://github.com/99designs/gqlgen/blob/d5c9f896419142f6378639b6eec93584fbf829ed/internal/code/util.go#L41
func QualifyPackagePath(importPath string) string {
	wd, _ := os.Getwd()

	// in go module mode, the import path doesn't need fixing
	if _, ok := goModuleRoot(wd); ok {
		return importPath
	}

	pkg, err := build.Import(importPath, wd, 0)
	if err != nil {
		return importPath
	}

	return pkg.ImportPath
}

var invalidPackageNameChar = regexp.MustCompile(`\W`)

// SanitizePackageName replaces invalid characters in a package name with underscores.
// See https://github.com/99designs/gqlgen/blob/d5c9f896419142f6378639b6eec93584fbf829ed/internal/code/util.go#L59
func SanitizePackageName(pkg string) string {
	return invalidPackageNameChar.ReplaceAllLiteralString(filepath.Base(pkg), "_")
}
