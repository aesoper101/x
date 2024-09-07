package args

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Pack(c interface{}, argNameGetters ...func(f reflect.StructField) string) (res []string, err error) {
	t := reflect.TypeOf(c)
	v := reflect.ValueOf(c)
	if reflect.TypeOf(c).Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("passed c must be struct or pointer of struct")
	}

	argNameGetter := func(f reflect.StructField) string {
		for _, g := range argNameGetters {
			if res := g(f); len(res) > 0 {
				return res
			}
		}
		return f.Name
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		x := v.Field(i)
		n := argNameGetter(f)

		if x.IsZero() {
			continue
		}

		switch x.Kind() {
		case reflect.Bool:
			if x.Bool() == false {
				continue
			}
			res = append(res, n+"="+fmt.Sprint(x.Bool()))
		case reflect.String:
			if x.String() == "" {
				continue
			}
			res = append(res, n+"="+x.String())
		case reflect.Slice:
			if x.Len() == 0 {
				continue
			}
			ft := f.Type.Elem()
			if ft.Kind() != reflect.String {
				return nil, fmt.Errorf("slice field %v must be '[]string', err: %v", f.Name, err)
			}
			var ss []string
			for i := 0; i < x.Len(); i++ {
				ss = append(ss, x.Index(i).String())
			}
			res = append(res, n+"="+strings.Join(ss, ";"))
		case reflect.Map:
			if x.Len() == 0 {
				continue
			}
			fk := f.Type.Key()
			if fk.Kind() != reflect.String {
				return nil, fmt.Errorf("map field %v must be 'map[string]string', err: %v", f.Name, err)
			}
			fv := f.Type.Elem()
			if fv.Kind() != reflect.String {
				return nil, fmt.Errorf("map field %v must be 'map[string]string', err: %v", f.Name, err)
			}
			var sk []string
			it := x.MapRange()
			for it.Next() {
				sk = append(sk, it.Key().String()+"="+it.Value().String())
			}
			res = append(res, n+"="+strings.Join(sk, ";"))
		default:
			return nil, fmt.Errorf("unsupported field type: %+v, err: %v", f, err)
		}
	}
	return res, nil
}

func Unpack(args []string, c interface{}, argNameGetters ...func(f reflect.StructField) string) error {
	m, err := MapForm(args)
	if err != nil {
		return fmt.Errorf("unmarshal args failed, err: %v", err.Error())
	}

	t := reflect.TypeOf(c).Elem()
	v := reflect.ValueOf(c).Elem()
	if t.Kind() != reflect.Struct {
		return errors.New("passed c must be struct or pointer of struct")
	}

	argNameGetter := func(f reflect.StructField) string {
		for _, g := range argNameGetters {
			if res := g(f); len(res) > 0 {
				return res
			}
		}
		return f.Name
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		x := v.Field(i)
		n := argNameGetter(f)
		values, ok := m[n]
		if !ok || len(values) == 0 || values[0] == "" {
			continue
		}
		switch x.Kind() {
		case reflect.Bool:
			if len(values) != 1 {
				return fmt.Errorf("field %s can't be assigned multi values: %v", n, values)
			}
			x.SetBool(values[0] == "true")
		case reflect.String:
			if len(values) != 1 {
				return fmt.Errorf("field %s can't be assigned multi values: %v", n, values)
			}
			x.SetString(values[0])
		case reflect.Slice:
			if len(values) != 1 {
				return fmt.Errorf("field %s can't be assigned multi values: %v", n, values)
			}
			ss := strings.Split(values[0], ";")
			if x.Type().Elem().Kind() == reflect.Int {
				n := reflect.MakeSlice(x.Type(), len(ss), len(ss))
				for i, s := range ss {
					val, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						return err
					}
					n.Index(i).SetInt(val)
				}
				x.Set(n)
			} else {
				for _, s := range ss {
					val := reflect.Append(x, reflect.ValueOf(s))
					x.Set(val)
				}
			}
		case reflect.Map:
			if len(values) != 1 {
				return fmt.Errorf("field %s can't be assigned multi values: %v", n, values)
			}
			ss := strings.Split(values[0], ";")
			out := make(map[string]string, len(ss))
			for _, s := range ss {
				sk := strings.SplitN(s, "=", 2)
				if len(sk) != 2 {
					return fmt.Errorf("map filed %v invalid key-value pair '%v'", n, s)
				}
				out[sk[0]] = sk[1]
			}
			x.Set(reflect.ValueOf(out))
		default:
			return fmt.Errorf("field %s has unsupported type %+v", n, f.Type)
		}
	}
	return nil
}

func MapForm(input []string) (map[string][]string, error) {
	out := make(map[string][]string, len(input))

	for _, str := range input {
		parts := strings.SplitN(str, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid argument: '%s'", str)
		}
		key, val := parts[0], parts[1]
		out[key] = append(out[key], val)
	}

	return out, nil
}
