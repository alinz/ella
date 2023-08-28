package golang

import (
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/astutil"
	"ella.to/pkg/sliceutil"
	"ella.to/pkg/strcase"
)

type MessageField struct {
	Name string
	Type string
	Tags string
}

type MessageFields []MessageField

func (m *MessageFields) Parse(message *ast.Message) error {
	*m = sliceutil.Mapper(message.Fields, func(field *ast.Field) MessageField {
		typ := parseType(field.Type)
		return MessageField{
			Name: field.Name.String(),
			Type: typ,
			Tags: parseMessageFieldOptions(field),
		}
	})
	return nil
}

type Message struct {
	Name   string
	Fields MessageFields
}

type Messages []Message

func (m *Messages) Parse(prog *ast.Program) error {
	*m = sliceutil.Mapper(astutil.GetMessages(prog), func(message *ast.Message) Message {
		msg := Message{
			Name: message.Name.String(),
		}

		msg.Fields.Parse(message)

		return msg
	})

	return nil
}

func parseMessageFieldOptions(field *ast.Field) string {
	var sb strings.Builder

	mapper := make(map[string]ast.Value)
	for _, opt := range field.Options {
		mapper[strings.ToLower(opt.Name.Token.Literal)] = opt.Value
	}

	jsonTagValue := strings.ToLower(strcase.ToSnake(field.Name.String()))

	jsonValue, ok := mapper["json"]
	if ok {
		switch jsonValue := jsonValue.(type) {
		case *ast.ValueString:
			jsonTagValue = jsonValue.Token.Literal
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
			yamlTagValue = yamlValue.Token.Literal
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
