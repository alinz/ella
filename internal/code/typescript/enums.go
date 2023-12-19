package typescript

import (
	"compiler.ella.to/internal/ast"
	"compiler.ella.to/internal/ast/astutil"
	"compiler.ella.to/pkg/sliceutil"
	"compiler.ella.to/pkg/strcase"
)

type EnumKeyValue struct {
	Name  string
	Value string
}

type Enum struct {
	Name string
	Keys []EnumKeyValue
}

type Enums []Enum

func (e *Enums) Parse(prog *ast.Program) error {
	*e = sliceutil.Mapper(astutil.GetEnums(prog), func(enum *ast.Enum) Enum {
		return Enum{
			Name: enum.Name.String(),
			Keys: sliceutil.Mapper(sliceutil.Filter(enum.Sets, func(set *ast.EnumSet) bool {
				return set.Name.String() != "_"
			}), func(set *ast.EnumSet) EnumKeyValue {
				return EnumKeyValue{
					Name:  set.Name.String(),
					Value: strcase.ToSnake(set.Name.String()),
				}
			}),
		}
	})

	return nil
}
