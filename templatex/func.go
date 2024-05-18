package templatex

import (
	"text/template"
	ttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/iancoleman/strcase"
)

var defaultFuncMaps = template.FuncMap{
	"ToSnake":              strcase.ToSnake,
	"ToSnakeWithIgnore":    strcase.ToSnakeWithIgnore,
	"ToScreamingSnake":     strcase.ToScreamingSnake,
	"ToKebab":              strcase.ToKebab,
	"ToScreamingKebab":     strcase.ToScreamingKebab,
	"ToDelimited":          strcase.ToDelimited,
	"ToScreamingDelimited": strcase.ToScreamingDelimited,
	"ToCamel":              strcase.ToCamel,
	"ToLowerCamel":         strcase.ToLowerCamel,
}

func TxtFuncMap() ttemplate.FuncMap {
	funcMap := sprig.TxtFuncMap()
	for k, v := range defaultFuncMaps {
		funcMap[k] = v
	}
	return funcMap
}

func HtmlFuncMap() ttemplate.FuncMap {
	funcMap := sprig.HtmlFuncMap()
	for k, v := range defaultFuncMaps {
		funcMap[k] = v
	}
	return funcMap
}

func GenericFuncMap() ttemplate.FuncMap {
	funcMap := sprig.GenericFuncMap()
	for k, v := range defaultFuncMaps {
		funcMap[k] = v
	}
	return funcMap
}
