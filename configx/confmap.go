package configx

import (
	"bytes"
	"encoding/json"

	"github.com/knadh/koanf/maps"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type KoanfConfmap struct {
	tuples []tuple
}

func NewKoanfConfmap(tuples []tuple) *KoanfConfmap {
	return &KoanfConfmap{tuples: jsonify(tuples)}
}

func jsonify(tuples []tuple) []tuple {
	for k, t := range tuples {
		var parsed interface{}
		switch vt := t.Value.(type) {
		case string:
			if gjson.Valid(vt) && json.NewDecoder(bytes.NewBufferString(vt)).Decode(&parsed) == nil {
				tuples[k].Value = parsed
			}
			continue
		case []byte:
			if gjson.ValidBytes(vt) && json.NewDecoder(bytes.NewBuffer(vt)).Decode(&parsed) == nil {
				tuples[k].Value = parsed
			}
			continue
		case json.RawMessage:
			if gjson.ValidBytes(vt) && json.NewDecoder(bytes.NewBuffer(vt)).Decode(&parsed) == nil {
				tuples[k].Value = parsed
			}
			continue
		}
	}
	return tuples
}

// ReadBytes is not supported by the env provider.
func (e *KoanfConfmap) ReadBytes() ([]byte, error) {
	return nil, errors.New("confmap provider does not support this method")
}

// Read returns the loaded map[string]interface{}.
func (e *KoanfConfmap) Read() (map[string]interface{}, error) {
	values := map[string]interface{}{}
	for _, t := range e.tuples {
		values[t.Key] = t.Value
	}

	// Ensure any nested values are properly converted as well
	cp := maps.Copy(values)
	maps.IntfaceKeysToStrings(cp)
	cp = maps.Unflatten(cp, Delimiter)

	return cp, nil
}
