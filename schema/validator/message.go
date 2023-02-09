package validator

import (
	"sort"

	"github.com/alinz/rpc.go/schema/ast"
)

func validateMessage(message *ast.Message, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) error {
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
