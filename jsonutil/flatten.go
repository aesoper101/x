package jsonutil

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

func Flatten(raw json.RawMessage) map[string]interface{} {
	parsed := gjson.ParseBytes(raw)
	if !parsed.IsObject() {
		return nil
	}

	flattened := make(map[string]interface{})
	flatten(parsed, nil, flattened)
	return flattened
}

func flatten(parsed gjson.Result, parents []string, flattened map[string]interface{}) {
	if parsed.IsObject() {
		parsed.ForEach(func(k, v gjson.Result) bool {
			flatten(v, append(parents, strings.ReplaceAll(k.String(), ".", "\\.")), flattened)
			return true
		})
	} else if parsed.IsArray() {
		for kk, vv := range parsed.Array() {
			flatten(vv, append(parents, strconv.Itoa(kk)), flattened)
		}
	} else {
		flattened[strings.Join(parents, ".")] = parsed.Value()
	}
}
