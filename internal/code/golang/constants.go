package golang

import (
	"fmt"
	"strconv"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/astutil"
	"ella.to/internal/token"
	"ella.to/pkg/sliceutil"
)

type Constant struct {
	Name  string
	Value string
}

type Constants []Constant

func (c *Constants) Parse(prog *ast.Program) error {
	*c = sliceutil.Mapper(astutil.GetConstants(prog), func(constant *ast.Const) Constant {
		var value string

		switch v := constant.Value.(type) {
		case *ast.ValueString:
			if v.Token.Type == token.ConstStringSingleQuote {
				value = fmt.Sprintf(`"%s"`, strings.ReplaceAll(v.Token.Literal, `"`, `\"`))
			} else {
				value = constant.Value.String()
			}
		case *ast.ValueInt:
			value = strconv.FormatInt(v.Value, 10)
		case *ast.ValueByteSize:
			value = fmt.Sprintf(`%d // %s`, v.Value*int64(v.Scale), v.String())
		case *ast.ValueDuration:
			value = fmt.Sprintf(`%d // %s`, v.Value*int64(v.Scale), v.String())
		default:
			value = constant.Value.String()
		}

		return Constant{
			Name:  constant.Name.String(),
			Value: value,
		}
	})

	return nil
}
