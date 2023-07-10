package golang

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/code"
	"ella.to/pkg/strcase"
)

func generateMessages(out io.Writer, messages []*ast.Message, isMessage func(string) bool) error {
	tmpl, err := code.LoadTemplate(files, "messages")
	if err != nil {
		return err
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Name.String() < messages[j].Name.String()
	})

	templMessages := code.Mapper(messages, func(message *ast.Message) any {
		return struct {
			Name   string
			Fields []any
		}{
			Name: message.Name.String(),
			Fields: code.Mapper(message.Fields, func(field *ast.Field) any {
				typ := parseType(field.Type)
				if isMessage(field.Type.String()) {
					typ = fmt.Sprintf("*%s", typ)
				}

				return struct {
					Name string
					Type string
					Tags string
				}{
					Name: field.Name.String(),
					Type: typ,
					Tags: parseFieldOptions(field),
				}
			}),
		}
	})

	return tmpl.Execute(out, templMessages)
}

func parseFieldOptions(field *ast.Field) string {
	var sb strings.Builder

	mapper := make(map[string]ast.Value)
	for _, opt := range field.Options {
		mapper[strings.ToLower(opt.Name.Token.Val)] = opt.Value
	}

	jsonTagValue := strings.ToLower(strcase.ToSnake(field.Name.String()))

	jsonValue, ok := mapper["json"]
	if ok {
		switch jsonValue := jsonValue.(type) {
		case *ast.ValueString:
			jsonTagValue = jsonValue.Token.Val
		case *ast.ValueBool:
			if !jsonValue.Value {
				jsonTagValue = "-"
			}
		}
	}

	sb.WriteString(`json:"`)
	sb.WriteString(jsonTagValue)
	sb.WriteString(`"`)

	yamlTagValue := strings.ToLower(strcase.ToSnake(field.Name.String()))

	yamlValue, ok := mapper["yaml"]
	if ok {
		switch yamlValue := yamlValue.(type) {
		case *ast.ValueString:
			yamlTagValue = yamlValue.Token.Val
		case *ast.ValueBool:
			if !yamlValue.Value {
				yamlTagValue = "-"
			}
		}
	}

	sb.WriteString(` yaml:"`)
	sb.WriteString(yamlTagValue)
	sb.WriteString(`"`)

	return sb.String()
}

func parseType(typ ast.Type) string {
	switch typ := typ.(type) {
	case *ast.CustomType:
		return typ.String()
	case *ast.Any:
		return "any"
	case *ast.Int:
		return fmt.Sprintf("int%d", typ.Size)
	case *ast.Uint:
		return fmt.Sprintf("uint%d", typ.Size)
	case *ast.Byte:
		return "byte"
	case *ast.Float:
		return fmt.Sprintf("float%d", typ.Size)
	case *ast.String:
		return "string"
	case *ast.Bool:
		return "bool"
	case *ast.Timestamp:
		return "time.Time"
	case *ast.Map:
		return fmt.Sprintf("map[%s]%s", parseType(typ.Key), parseType(typ.Value))
	case *ast.Array:
		return fmt.Sprintf("[]%s", parseType(typ.Type))
	}

	// This shouldn't happen as the validator should catch this any errors
	panic(fmt.Sprintf("unknown type: %T", typ))
}
