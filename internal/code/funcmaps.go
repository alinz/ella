package code

import (
	"strings"
	"text/template"

	"ella.to/pkg/strcase"
)

var DefaultFuncsMap = template.FuncMap{
	"ToLower":      strings.ToLower,
	"ToUpper":      strings.ToUpper,
	"ToPascalCase": strcase.ToPascal,
	"ToCamelCase":  strcase.ToCamel,
	"ToSnakeCase":  strcase.ToSnake,
	"ToKebabCase":  strcase.ToPascal,

	"AppendPostfix": func(postfix, s string) string {
		return s + postfix
	},

	"ReplaceAll": func(old, new, s string) string {
		return strings.ReplaceAll(s, old, new)
	},
	"Split": func(sep, s string) []string {
		return strings.Split(s, sep)
	},
	"Join": func(sep string, pick int, s []string) string {
		var result []string

		for i, str := range s {
			if i%pick == 0 {
				result = append(result, str)
			}
		}

		return strings.Join(result, sep)
	},
}
