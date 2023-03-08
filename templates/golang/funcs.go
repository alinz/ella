package golang

import (
	"net/http"
	"strings"
	"text/template"

	"ella.to/pkg/stringcase"
)

func getSliceLength[T any]() func(arr []T) int {
	return func(arr []T) int {
		return len(arr)
	}
}

var tmplFuncs = template.FuncMap{
	"ToUpper":          strings.ToUpper,
	"ToLower":          strings.ToLower,
	"ToCamel":          stringcase.ToCamel,
	"ToPascal":         stringcase.ToPascal,
	"GetArgsLength":    getSliceLength[Arg](),
	"GetReturnsLength": getSliceLength[Return](),

	"MethodArgs": func(args []Arg) string {
		var sb strings.Builder

		sb.WriteString("ctx context.Context")

		for _, arg := range args {
			sb.WriteString(", ")
			sb.WriteString(arg.Name)
			sb.WriteString(" ")
			sb.WriteString(arg.Type)
		}

		return sb.String()
	},

	"MethodReturns":           MethodReturns(false),
	"MethodReturnsIgnoreName": MethodReturns(true),

	"IsStream": func(returns []Return) bool {
		for _, ret := range returns {
			if ret.Stream {
				return true
			}
		}

		return false
	},

	"ToStructArgs": toStructArgs,

	"ToStructReturns": toStructReturns,

	"ToExtractArgs": func(args []Arg) string {
		var sb strings.Builder

		for i, arg := range args {
			if i > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString("args.")
			sb.WriteString(stringcase.ToPascal(arg.Name))
			sb.WriteString(",")
		}

		return sb.String()
	},
	"ToStreamReturnsName": func(returns []Return) string {
		var sb strings.Builder

		for i, ret := range returns {
			if i > 0 {
				sb.WriteString(", ")
			}

			sb.WriteString(ret.Name)
			sb.WriteString("Stream")
		}

		return sb.String()
	},

	"ReturnsExtract": func(returns []Return) string {
		var sb strings.Builder

		for i, ret := range returns {
			if i > 0 {
				sb.WriteString(", ")
			}

			sb.WriteString("ret.")
			sb.WriteString(stringcase.ToPascal(ret.Name))
		}

		if sb.Len() > 0 {
			sb.WriteString(", err")
		} else {
			sb.WriteString("err")
		}
		return sb.String()
	},

	"MethodArgsStructClient": func(args []Arg) string {
		var sb strings.Builder

		if len(args) > 0 {
			sb.WriteString("struct {\n")
			sb.WriteString(toStructArgs(args))
			sb.WriteString("}{\n")
			for _, arg := range args {
				sb.WriteString(stringcase.ToPascal(arg.Name))
				sb.WriteString(": ")
				sb.WriteString(arg.Name)
				sb.WriteString(",\n")
			}
			sb.WriteString("}\n")
		} else {
			sb.WriteString("emptyStruct{}")
		}

		return sb.String()
	},

	"MethodReturnsStructClient": func(returns []Return) string {
		var sb strings.Builder

		if len(returns) > 0 {
			sb.WriteString("struct {\n")
			sb.WriteString(toStructReturns(returns))
			sb.WriteString("}{}")
		} else {
			sb.WriteString("emptyStruct{}")
		}

		return sb.String()
	},

	"ServiceMethodHttpMethod": func(method Method) string {
		for _, option := range method.Options {
			if option.Key == "HttpMethod" {
				value := option.Value.(string)
				value = strings.ToUpper(value)
				value = strings.ReplaceAll(value, "\"", "")
				return value
			}
		}

		return http.MethodPost
	},

	"ReturnsOut": func(returns []Return) string {
		var sb strings.Builder

		for i, ret := range returns {
			if i > 0 {
				sb.WriteString(", ")
			}

			sb.WriteString("out.")
			sb.WriteString(stringcase.ToPascal(ret.Name))
		}

		if sb.Len() > 0 {
			sb.WriteString(", err")
		} else {
			sb.WriteString("err")
		}

		return sb.String()
	},
}
