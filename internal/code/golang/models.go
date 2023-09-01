package golang

import (
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/ast/astutil"
	"ella.to/pkg/sliceutil"
	"ella.to/pkg/strcase"
)

type ModelField struct {
	Name string
	Type string
	Tags string
}

type ModelFields []ModelField

func (m *ModelFields) Parse(message *ast.Model) error {
	*m = sliceutil.Mapper(message.Fields, func(field *ast.Field) ModelField {
		typ := parseType(field.Type)
		return ModelField{
			Name: field.Name.String(),
			Type: typ,
			Tags: parseModelFieldOptions(field),
		}
	})
	return nil
}

type Model struct {
	Name   string
	Fields ModelFields
}

type Models []Model

func (m *Models) Parse(prog *ast.Program) error {
	*m = sliceutil.Mapper(astutil.GetModels(prog), func(message *ast.Model) Model {
		msg := Model{
			Name: message.Name.String(),
		}

		msg.Fields.Parse(message)

		return msg
	})

	return nil
}

func parseModelFieldOptions(field *ast.Field) string {
	var sb strings.Builder

	mapper := make(map[string]ast.Value)
	for _, opt := range field.Options {
		mapper[strings.ToLower(opt.Name.Token.Literal)] = opt.Value
	}

	jsonTagValue := strings.ToLower(strcase.ToCamel(field.Name.String()))

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

	yamlTagValue := strings.ToLower(strcase.ToCamel(field.Name.String()))

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
