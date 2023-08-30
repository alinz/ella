package typescript

import (
	"fmt"
	"strconv"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func getValue(value ast.Value) string {
	switch v := value.(type) {
	case *ast.ValueString:
		if v.Token.Type == token.ConstStringSingleQuote {
			return fmt.Sprintf(`"%s"`, strings.ReplaceAll(v.Token.Literal, `"`, `\"`))
		} else {
			return value.String()
		}
	case *ast.ValueInt:
		return strconv.FormatInt(v.Value, 10)
	case *ast.ValueByteSize:
		return fmt.Sprintf(`%d`, v.Value*int64(v.Scale))
	case *ast.ValueDuration:
		return fmt.Sprintf(`%d`, v.Value*int64(v.Scale))
	default:
		return value.String()
	}
}
