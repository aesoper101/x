package code

import (
	"testing"

	"github.com/aesoper101/x/filex"
	"github.com/stretchr/testify/require"
)

func TestFindGoModuleRoot(t *testing.T) {
	module, ok := goModuleRoot(filex.Getwd())
	require.Equal(t, true, ok)
	require.Equal(t, "github.com/aesoper101/x/golangx/code", module)
}
