package golang

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/code"
	"ella.to/internal/token"
)

func generateConstants(out io.Writer, pkg string, constants []*ast.Const) error {
	tmpl, err := code.LoadTemplate(files, "constants")
	if err != nil {
		return err
	}

	sort.Slice(constants, func(i, j int) bool {
		return constants[i].Name.String() < constants[j].Name.String()
	})

	templConsts := code.Mapper(constants, func(constant *ast.Const) any {
		var value string

		if v, ok := constant.Value.(*ast.ValueString); ok {
			if v.Token.Type == token.ConstStringSingleQuote {
				value = fmt.Sprintf(`"%s"`, strings.ReplaceAll(v.Token.Val, `"`, `\"`))
			} else {
				value = constant.Value.String()
			}
		} else {
			value = constant.Value.String()
		}

		return &struct {
			Name  string
			Value string
		}{
			Name:  constant.Name.String(),
			Value: value,
		}
	})

	return tmpl.Execute(out, templConsts)
}
