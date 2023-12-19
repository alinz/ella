package golang

import (
	"fmt"

	"compiler.ella.to/internal/ast"
	"compiler.ella.to/internal/ast/astutil"
	"compiler.ella.to/pkg/sliceutil"
)

type EnumKeyValue struct {
	Name  string
	Value string
}

type Enum struct {
	Name string
	Type string // int8, int16, int32, int64
	Keys []EnumKeyValue
}

type Enums []Enum

func (e *Enums) Parse(prog *ast.Program) error {
	*e = sliceutil.Mapper(astutil.GetEnums(prog), func(enum *ast.Enum) Enum {
		return Enum{
			Name: enum.Name.String(),
			Type: fmt.Sprintf("int%d", enum.Size),
			Keys: sliceutil.Mapper(enum.Sets, func(set *ast.EnumSet) EnumKeyValue {
				return EnumKeyValue{
					Name:  set.Name.String(),
					Value: fmt.Sprintf("%d", set.Value.Value),
				}
			}),
		}
	})

	return nil
}
