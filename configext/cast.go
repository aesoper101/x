package configext

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cast"
	"reflect"
	"strings"
)

// ToFloatSlice casts an interface to a []float64 type.
func toFloatSlice(i interface{}) []float64 {
	f, _ := toFloatSliceE(i)
	return f
}

// ToFloatSliceE casts an interface to a []float64 type.
func toFloatSliceE(i interface{}) ([]float64, error) {
	if i == nil {
		return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
	}

	switch v := i.(type) {
	case []float64:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]float64, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToFloat64E(s.Index(j).Interface())
			if err != nil {
				return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
	}
}

// ToStringSlice casts an interface to a []string type and respects comma-separated values.
func toStringSlice(i interface{}) []string {
	s, _ := toStringSliceE(i)
	return s
}

// ToStringSliceE casts an interface to a []string type and respects comma-separated values.
func toStringSliceE(i interface{}) ([]string, error) {
	switch s := i.(type) {
	case string:
		return parseCSV(s)
	}

	return cast.ToStringSliceE(i)
}

func parseCSV(v string) ([]string, error) {
	stringReader := strings.NewReader(v)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}
