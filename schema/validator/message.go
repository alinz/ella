package validator

import (
	"sort"

	"ella.to/pkg/stringcase"
	"ella.to/schema/ast"
	"ella.to/schema/token"
)

func addStringFieldOption(key, value string) *ast.Constant {
	return &ast.Constant{
		Name: &ast.Identifier{
			Name: key,
			Token: &token.Token{
				Kind: token.Word,
				Val:  key,
			},
		},
		Value: &ast.ValueString{
			Token: &token.Token{
				Kind: token.ConstantString,
				Val:  value,
			},
			Content: value,
		},
	}
}

func validateMessage(message *ast.Message, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) error {
	for _, field := range message.Fields {
		defineJson := false
		for _, option := range field.Options {
			if option.Name.Name == "json" {
				defineJson = true
				break
			}
		}

		if !defineJson {
			field.Options = append(field.Options, addStringFieldOption("json", stringcase.ToSnake(field.Name.Name)))
		}
	}

	return nil
}

func validateMessages(messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) ([]*ast.Message, error) {
	messages := make([]*ast.Message, 0, len(messagesMap))

	for _, message := range messagesMap {
		if err := validateMessage(message, messagesMap, enumsMap); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Name.Name < messages[j].Name.Name
	})

	return messages, nil
}
