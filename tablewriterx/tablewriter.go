package tablewriterx

import (
	"errors"
	"fmt"
	cond "github.com/aesoper101/x/condition"
	"io"
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
)

type TableWriter struct {
	*tablewriter.Table
}

type TableWriterOptions struct {
	writer io.Writer
}

type TableWriterOption func(*TableWriterOptions)

func WithWriter(writer io.Writer) TableWriterOption {
	return func(options *TableWriterOptions) {
		options.writer = writer
	}
}

func NewWriter(options ...TableWriterOption) *TableWriter {
	opts := &TableWriterOptions{
		writer: os.Stdout,
	}
	for _, o := range options {
		o(opts)
	}
	return &TableWriter{tablewriter.NewWriter(opts.writer)}
}

func (t *TableWriter) SetStructs(v interface{}) error {
	if v == nil {
		return errors.New("nil value")
	}
	vt := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)
	switch vt.Kind() {
	case reflect.Slice, reflect.Array:
		if vv.Len() < 1 {
			return errors.New("empty value")
		}

		// check first element to set header
		first := vv.Index(0)
		e := first.Type()
		switch e.Kind() {
		case reflect.Struct:
			// OK
		case reflect.Ptr:
			if first.IsNil() {
				return errors.New("the first element is nil")
			}
			e = first.Elem().Type()
			if e.Kind() != reflect.Struct {
				return fmt.Errorf("invalid kind %s", e.Kind())
			}
		default:
			return fmt.Errorf("invalid kind %s", e.Kind())
		}
		n := e.NumField()
		var headers []string
		ignores := make(map[int]struct{})
		for i := 0; i < n; i++ {
			f := e.Field(i)
			header := f.Tag.Get("tablewriter")
			if header == "-" {
				ignores[i] = struct{}{}
				continue
			}
			if header == "" {
				header = cond.Ternary(f.Tag.Get("json") == "", f.Name, f.Tag.Get("json"))
			}
			headers = append(headers, header)
		}
		t.SetHeader(headers)

		for i := 0; i < vv.Len(); i++ {
			item := reflect.Indirect(vv.Index(i))
			itemType := reflect.TypeOf(item)
			switch itemType.Kind() {
			case reflect.Struct:
				// OK
			default:
				return fmt.Errorf("invalid item type %v", itemType.Kind())
			}
			if !item.IsValid() {
				// skip rendering
				continue
			}

			nf := item.NumField()
			if n != nf {
				return errors.New("invalid num of field")
			}
			var rows []string
			for j := 0; j < nf; j++ {
				f := reflect.Indirect(item.Field(j))
				if f.Kind() == reflect.Ptr {
					f = f.Elem()
				}
				if _, ok := ignores[j]; ok {
					continue
				}
				if f.IsValid() {
					if s, ok := f.Interface().(fmt.Stringer); ok {
						rows = append(rows, s.String())
						continue
					}
					rows = append(rows, fmt.Sprint(f))
				} else {
					rows = append(rows, "nil")
				}
			}
			t.Append(rows)
		}
	default:
		return fmt.Errorf("invalid type %T", v)
	}
	return nil
}
