package code

import (
	"strings"
	"text/template"

	"compiler.ella.to/pkg/strcase"
)

var DefaultFuncsMap = template.FuncMap{
	"ToLower":      strings.ToLower,
	"ToUpper":      strings.ToUpper,
	"ToPascalCase": strcase.ToPascal,
	"ToCamelCase":  strcase.ToCamel,
	"ToSnakeCase":  strcase.ToSnake,
}
