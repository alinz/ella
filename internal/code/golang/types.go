package golang

import (
	"fmt"

	"compiler.ella.to/internal/ast"
)

func parseType(typ ast.Type, isModelType func(value string) bool) string {
	switch typ := typ.(type) {
	case *ast.CustomType:
		val := typ.String()
		if isModelType(val) {
			return "*" + val
		}
		return val
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
		return fmt.Sprintf("map[%s]%s", parseType(typ.Key, isModelType), parseType(typ.Value, isModelType))
	case *ast.Array:
		return fmt.Sprintf("[]%s", parseType(typ.Type, isModelType))
	}

	// This shouldn't happen as the validator should catch this any errors
	panic(fmt.Sprintf("unknown type: %T", typ))
}
