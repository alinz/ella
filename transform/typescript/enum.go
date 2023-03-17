package typescript

import (
	"ella.to/schema/ast"
	"ella.to/transform"
)

func Enum(enum *ast.Enum) transform.Func {
	return func(out transform.Writer) error {
		out.Str(`export enum `).Pascal(enum.Name.Name).Str(` {`).Lines(1)

		for _, v := range enum.Constants {
			out.
				Tabs(1).
				Pascal(v.Name.Name).Str(` = `).Str(`"`).Pascal(v.Name.Name).Str(`",`).
				Lines(1)
		}

		out.Str(`}`)

		return nil
	}
}

func Enums(enums []*ast.Enum) transform.Func {
	return func(out transform.Writer) error {
		out.Str("// Enums").Lines(2)

		for _, enum := range enums {
			if err := Enum(enum)(out); err != nil {
				return err
			}
			out.Lines(2)
		}

		return nil
	}
}
