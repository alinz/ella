package golang

import (
	"fmt"

	"ella.to/internal/ast"
	"ella.to/internal/astutil"
	"ella.to/pkg/sliceutil"
	"ella.to/pkg/strcase"
)

type Option struct {
	Name  string
	Type  string
	Value string
}

type Base struct {
	Name      string
	NameValue string
	Type      string
	Options   []Option
}

type Bases []Base

func (b *Bases) Parse(prog *ast.Program) error {
	*b = sliceutil.Mapper(astutil.GetBases(prog), func(base *ast.Base) Base {
		name := fmt.Sprintf("%sBaseOptions", strcase.ToCamel(base.Name.String()))
		nameValue := fmt.Sprintf("%sValues", name)

		return Base{
			Name:      name,
			NameValue: nameValue,
			Type:      parseType(base.Type),
			Options: sliceutil.Mapper(base.Options, func(opt *ast.Option) Option {
				return Option{
					Name:  opt.Name.String(),
					Type:  parseValueType(opt.Value),
					Value: opt.Value.String(),
				}
			}),
		}
	})

	return nil
}
