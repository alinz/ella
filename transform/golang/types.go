package golang

import (
	"fmt"

	"github.com/alinz/ella.to/schema/ast"
)

func parseType(typ ast.Type, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) (string, error) {
	switch typ := typ.(type) {
	case *ast.TypeCustom:
		// all the Message type should be pointer
		// all the Enum type should be value
		if _, ok := messagesMap[typ.Name]; ok {
			return fmt.Sprintf("*%s", typ.Name), nil
		}
		if _, ok := enumsMap[typ.Name]; ok {
			return typ.Name, nil
		}

		return "", fmt.Errorf("it's neither Enum Nor Message type: %s", typ.Name)
	case *ast.TypeAny:
		return "any", nil
	case *ast.TypeInt:
		return fmt.Sprintf("int%d", typ.Size), nil
	case *ast.TypeUint:
		return fmt.Sprintf("uint%d", typ.Size), nil
	case *ast.TypeByte:
		return "byte", nil
	case *ast.TypeFloat:
		return fmt.Sprintf("float%d", typ.Size), nil
	case *ast.TypeString:
		return "string", nil
	case *ast.TypeBool:
		return "bool", nil
	case *ast.TypeTimestamp:
		return "time.Time", nil
	case *ast.TypeMap:
		key, err := parseType(typ.Key, messagesMap, enumsMap)
		if err != nil {
			return "", err
		}

		value, err := parseType(typ.Value, messagesMap, enumsMap)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("map[%s]%s", key, value), nil
	case *ast.TypeArray:
		name, err := parseType(typ.Type, messagesMap, enumsMap)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[]%s", name), nil
	}

	return "", fmt.Errorf("unknown type: %T", typ)
}
