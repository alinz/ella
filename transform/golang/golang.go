package golang

import (
	"github.com/alinz/rpc.go/schema/ast"
	"github.com/alinz/rpc.go/transform"
)

func Package(name string) transform.Func {
	return func(out transform.Writer) {
		out.
			String("package ").String(name).
			NewLines(2)
	}
}

func Constants(constants []*ast.Constant) transform.Func {
	return func(out transform.Writer) {
		if len(constants) == 0 {
			return
		}

		out.String("const (").Indents(1).NewLines(1)
		for _, constant := range constants {
			out.
				String(constant.Name.Name).
				String(" = ").
				String(constant.Value.TokenLiteral()).
				NewLines(1)
		}
		out.Indents(-1).String(")").NewLines(2)
	}
}

func Enum(enum *ast.Enum) transform.Func {
	return func(out transform.Writer) {
		out.
			String("type ").
			String(enum.Name.Name).
			String(" ").
			String(enum.Type.TokenLiteral()).
			NewLines(2)

		Constants(enum.Constants)(out)
	}
}

func Enums(enums []*ast.Enum) transform.Func {
	return func(out transform.Writer) {
		if len(enums) == 0 {
			return
		}

		for i, enum := range enums {
			if i != 0 {
				out.NewLines(2)
			}
			Enum(enum)(out)
		}
	}
}
