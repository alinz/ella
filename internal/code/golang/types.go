package golang

import (
	"fmt"

	"ella.to/internal/ast"
)

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
