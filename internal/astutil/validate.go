package astutil

import "ella.to/internal/ast"

func CreateMessageTypeMap(messages []*ast.Message) map[string]*ast.Message {
	messagesMap := make(map[string]*ast.Message)
	for _, message := range messages {
		messagesMap[message.Name.String()] = message
	}

	return messagesMap
}

func CreateIsMessageTypeFunc(messages []*ast.Message) func(value string) bool {
	messagesMap := make(map[string]struct{})
	for _, message := range messages {
		messagesMap[message.Name.String()] = struct{}{}
	}

	return func(value string) bool {
		_, ok := messagesMap[value]
		return ok
	}
}

func CreateIsEnumTypeFunc(enums []*ast.Enum) func(value string) bool {
	enumsMap := make(map[string]struct{})
	for _, enum := range enums {
		enumsMap[enum.Name.String()] = struct{}{}
	}

	return func(value string) bool {
		_, ok := enumsMap[value]
		return ok
	}
}

func CreateIsValidType(prog *ast.Program) func(typ ast.Type) bool {
	isEnumType := CreateIsEnumTypeFunc(GetEnums(prog))
	isMessageType := CreateIsMessageTypeFunc(GetMessages(prog))

	var isValidType func(typ ast.Type) bool

	isValidType = func(typ ast.Type) bool {
		switch v := typ.(type) {
		case *ast.Map:
			return IsTypeComparable(v.Key) && isValidType(v.Value)
		case *ast.CustomType:
			return isEnumType(v.TokenLiteral()) || isMessageType(v.TokenLiteral())
		case *ast.Array:
			return isValidType(v.Type)
		default:
			return true // the reason this returns true because parser already validate the type
		}
	}

	return isValidType
}

func CreateConstsMap(prog *ast.Program) map[string]*ast.Const {
	constantsMap := make(map[string]*ast.Const)
	for _, constant := range GetConstants(prog) {
		constantsMap[constant.Name.String()] = constant
	}

	return constantsMap
}

func IsTypeComparable(typ ast.Type) bool {
	switch typ.(type) {
	case *ast.Byte, *ast.Uint, *ast.Int, *ast.Float, *ast.String, *ast.Bool, *ast.Timestamp:
		return true
	default:
		return false
	}
}
