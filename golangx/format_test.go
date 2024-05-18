package golangx

import (
	"bytes"
	"go/parser"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatCodeFromSource(t *testing.T) {
	p, err := NewParser("", []byte(testSrc), parser.ParseComments)
	require.Zero(t, err)

	buf := bytes.NewBuffer(nil)
	err = FormatCodeFromParser(p, buf)
	require.Zero(t, err)

	t.Logf("%s", buf.String())
}
