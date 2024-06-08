package golangx

import (
	"github.com/aesoper101/x/filex"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSearchGoMod(t *testing.T) {
	module, _, ok := SearchGoMod(filex.Getwd())
	require.Equal(t, true, ok)
	require.Equal(t, "github.com/aesoper101/x", module)
}
