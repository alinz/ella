package typescript

import (
	"fmt"

	"ella.to/internal/ast"
)

func parseType(typ ast.Type) string {
	switch t := typ.(type) {
	case *ast.Bool:
		return `boolean`
	case *ast.Int:
		return `number`
	case *ast.Float:
		return `number`
	case *ast.Uint:
		return `number`
	case *ast.String:
		return `string`
	case *ast.Any:
		return `any`
	case *ast.Timestamp:
		return `string`
	case *ast.Array:
		typ := parseType(t.Type)
		return typ + "[]"
	case *ast.Map:
		key := parseType(t.Key)
		value := parseType(t.Value)
		return `{ [key: ` + key + `]: ` + value + ` }`
	case *ast.CustomType:
		return t.TokenLiteral()
	case *ast.Byte:
		return "byte"
	case *ast.File:
		return "fileupload"
	default:
		panic(fmt.Errorf("unknown type: %T", t))
	}
}
