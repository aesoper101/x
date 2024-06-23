package golangx

import (
	"github.com/aesoper101/x/fileutil"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSearchGoMod(t *testing.T) {
	module, _, ok := SearchGoMod(fileutil.Getwd())
	require.Equal(t, true, ok)
	require.Equal(t, "github.com/aesoper101/x", module)
}
