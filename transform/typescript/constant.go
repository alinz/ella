package typescript

import (
	"fmt"

	"ella.to/schema/ast"
	"ella.to/transform"
)

func Constant(constant *ast.Constant) transform.Func {
	return func(out transform.Writer) error {
		out.Str(`export const `).Pascal(constant.Name.Name).Str(` = `)

		switch v := constant.Value.(type) {
		case *ast.ValueBool:
			if v.Content {
				out.Str(`true`)
			} else {
				out.Str(`false`)
			}
		case *ast.ValueInt:
			out.Str("%d", v.Content)
		case *ast.ValueFloat:
			out.Str("%f", v.Content)
		case *ast.ValueString:
			out.Str(`"`).Str(v.Content).Str(`"`)
		default:
			return fmt.Errorf("unknown value type for constant value: %T", v)
		}

		out.Str(`;`)

		return nil
	}
}

func Constants(constants []*ast.Constant) transform.Func {
	return func(out transform.Writer) error {
		out.Str("// Constants").Lines(2)

		for _, constant := range constants {
			if err := Constant(constant)(out); err != nil {
				return err
			}
			out.Lines(1)
		}

		return nil
	}
}
