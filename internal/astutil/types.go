package astutil

import "ella.to/internal/ast"

func GetValueType(val ast.Value) string {
	switch v := val.(type) {
	case *ast.ValueString:
		return "string"
	case *ast.ValueInt:
		return "int"
	case *ast.ValueFloat:
		return "float"
	case *ast.ValueBool:
		return "bool"
	case *ast.ValueByteSize:
		return "int"
	case *ast.ValueDuration:
		return "int"
	case *ast.ValueUint:
		return "uint"
	case *ast.ValueNull:
		return "null"
	case *ast.ValueVariable:
		return v.TokenLiteral()
	}

	// if this happens, our parser has a bug
	panic("unknown value type, ella's paser has a bug")
}
