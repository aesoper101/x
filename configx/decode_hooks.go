package configx

import (
	"encoding"
	"encoding/json"
	"fmt"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"

	"github.com/go-viper/mapstructure/v2"
	"github.com/pkg/errors"
)

const (
	errFmtDecodeHookCouldNotParse           = "could not decode '%s' to a %s%s: %w"
	errFmtDecodeHookCouldNotParseBasic      = "could not decode to a %s%s: %w"
	errFmtDecodeHookCouldNotParseEmptyValue = "could not decode an empty value to a %s%s: %w"
)

var errDecodeNonPtrMustHaveValue = errors.New("must have a non-empty value")

func TextUnmarshalerHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		result := reflect.New(t).Interface()
		unmarshaller, ok := result.(encoding.TextUnmarshaler)
		if !ok {
			return data, nil
		}

		// default text representation is the actual value of the `from` string
		var (
			dataVal = reflect.ValueOf(data)
			text    = []byte(dataVal.String())
		)
		if f.Kind() == t.Kind() {
			// source and target are of underlying type string
			var (
				err    error
				ptrVal = reflect.New(dataVal.Type())
			)
			if !ptrVal.Elem().CanSet() {
				// cannot set, skip, this should not happen
				if err := unmarshaller.UnmarshalText(text); err != nil {
					return nil, err
				}
				return result, nil
			}
			ptrVal.Elem().Set(dataVal)

			// We need to assert that both, the value type and the pointer type
			// do (not) implement the TextMarshaller interface before proceeding and simply
			// using the string value of the string type.
			// it might be the case that the internal string representation differs from
			// the (un)marshalled string.

			for _, v := range []reflect.Value{dataVal, ptrVal} {
				if marshaller, ok := v.Interface().(encoding.TextMarshaler); ok {
					text, err = marshaller.MarshalText()
					if err != nil {
						return nil, err
					}
					break
				}
			}
		}

		// text is either the source string's value or the source string type's marshaled value
		// which may differ from its internal string value.
		if err := unmarshaller.UnmarshalText(text); err != nil {
			return nil, err
		}
		return result, nil
	}
}

func JsonUnmarshalerHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}

		if t.Kind() == reflect.Pointer && t.Elem().Kind() != reflect.Struct {
			return data, nil
		} else if t.Kind() != reflect.Struct {
			return data, nil
		}

		result := reflect.New(t).Interface()
		unmarshaller, ok := result.(json.Unmarshaler)
		if !ok {
			return data, nil
		}

		bytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		if err := unmarshaller.UnmarshalJSON(bytes); err != nil {
			return nil, err
		}

		if t.Kind() == reflect.Ptr {
			return result, nil
		}

		return reflect.Indirect(reflect.ValueOf(result)).Interface(), nil
	}
}

func StringToURLHookFunc() mapstructure.DecodeHookFuncType {
	return func(f, t reflect.Type, data any) (value any, err error) {
		var ptr bool

		if f.Kind() != reflect.String {
			return data, nil
		}

		prefixType := ""

		if t.Kind() == reflect.Ptr {
			ptr = true
			prefixType = "*"
		}

		expectedType := reflect.TypeOf(url.URL{})

		if ptr && t.Elem() != expectedType {
			return data, nil
		} else if !ptr && t != expectedType {
			return data, nil
		}

		dataStr := data.(string)

		var result *url.URL

		if dataStr != "" {
			if result, err = url.Parse(dataStr); err != nil {
				return nil, fmt.Errorf(errFmtDecodeHookCouldNotParse, dataStr, prefixType, expectedType, err)
			}
		}

		if ptr {
			return result, nil
		}

		if result == nil {
			return url.URL{}, nil
		}

		return *result, nil
	}
}

func StringToMailAddressHookFunc() mapstructure.DecodeHookFuncType {
	return func(f, t reflect.Type, data any) (value any, err error) {
		var ptr bool

		if f.Kind() != reflect.String {
			return data, nil
		}

		prefixType := ""

		if t.Kind() == reflect.Ptr {
			ptr = true
			prefixType = "*"
		}

		expectedType := reflect.TypeOf(mail.Address{})

		if ptr && t.Elem() != expectedType {
			return data, nil
		} else if !ptr && t != expectedType {
			return data, nil
		}

		dataStr := data.(string)

		var result *mail.Address

		if dataStr != "" {
			if result, err = mail.ParseAddress(dataStr); err != nil {
				return nil, fmt.Errorf(errFmtDecodeHookCouldNotParse, dataStr, prefixType, expectedType.String()+" (RFC5322)", err)
			}
		}

		if ptr {
			return result, nil
		}

		if result == nil {
			return mail.Address{}, nil
		}

		return *result, nil
	}
}

func StringToRegexpHookFunc() mapstructure.DecodeHookFuncType {
	return func(f, t reflect.Type, data any) (value any, err error) {
		var ptr bool

		if f.Kind() != reflect.String {
			return data, nil
		}

		prefixType := ""

		if t.Kind() == reflect.Ptr {
			ptr = true
			prefixType = "*"
		}

		expectedType := reflect.TypeOf(regexp.Regexp{})

		if ptr && t.Elem() != expectedType {
			return data, nil
		} else if !ptr && t != expectedType {
			return data, nil
		}

		dataStr := data.(string)

		var result *regexp.Regexp

		if dataStr != "" {
			if result, err = regexp.Compile(dataStr); err != nil {
				return nil, fmt.Errorf(errFmtDecodeHookCouldNotParse, dataStr, prefixType, expectedType, err)
			}
		}

		if ptr {
			return result, nil
		}

		if result == nil {
			return nil, fmt.Errorf(errFmtDecodeHookCouldNotParseEmptyValue, prefixType, expectedType, errDecodeNonPtrMustHaveValue)
		}

		return *result, nil
	}
}
