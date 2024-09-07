package glob

import (
	"strings"
	"testing"
)

func TestEffectiveCriteria(t *testing.T) {
	result, _ := EffectiveCriteria("test/*.txt", "test2/**/*.html", "test/*.js", "!test/*.html")
	t.Log(result.Roots())
	if len(result.Roots()) != 2 {
		t.Error("expected 2 items in result set")
	}

	success := 0
	for _, c := range result.Items {
		if strings.HasSuffix(c.Root, "/test") {
			success++
		}
		if strings.HasSuffix(c.Root, "/test2") {
			success++
		}
	}

	if success != 2 {
		t.Error("should calc effective criteria")
	}

}
