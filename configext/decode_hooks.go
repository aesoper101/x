package configext

import (
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
	errFmtDecodeHookCouldNotParseEmptyValue = "could not decode an empty value to a %s%s: %w"
)

var errDecodeNonPtrMustHaveValue = errors.New("must have a non-empty value")

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
				return nil, fmt.Errorf(
					errFmtDecodeHookCouldNotParse,
					dataStr,
					prefixType,
					expectedType.String()+" (RFC5322)",
					err,
				)
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
			return nil, fmt.Errorf(
				errFmtDecodeHookCouldNotParseEmptyValue,
				prefixType,
				expectedType,
				errDecodeNonPtrMustHaveValue,
			)
		}

		return *result, nil
	}
}

func MapToStructHookFunc() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}

		if t.Kind() != reflect.Struct {
			return data, nil
		}

		result := reflect.New(t).Interface()
		bs, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(bs, result)
		return result, err
	}
}
