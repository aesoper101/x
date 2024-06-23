package code

import (
	"testing"

	"github.com/aesoper101/x/fileutil"
	"github.com/stretchr/testify/require"
)

func TestFindGoModuleRoot(t *testing.T) {
	module, ok := goModuleRoot(fileutil.Getwd())
	require.Equal(t, true, ok)
	require.Equal(t, "github.com/aesoper101/x/golangx/tools/code", module)
}
