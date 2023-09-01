package typescript

import (
	"ella.to/internal/ast"
	"ella.to/internal/ast/astutil"
	"ella.to/pkg/sliceutil"
)

type Constant struct {
	Name  string
	Value string
}

type Constants []Constant

func (c *Constants) Parse(prog *ast.Program) error {
	*c = sliceutil.Mapper(astutil.GetConstants(prog), func(constant *ast.Const) Constant {
		return Constant{
			Name:  constant.Name.String(),
			Value: getValue(constant.Value),
		}
	})

	return nil
}
