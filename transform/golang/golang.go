package golang

import (
	"github.com/alinz/ella.to/schema/ast"
	"github.com/alinz/ella.to/transform"
)

func Package(name string) transform.Func {
	return func(out transform.Writer) error {
		out.
			String("package ").String(name).
			NewLines(2)

		return nil
	}
}

func Constants(constants []*ast.Constant) transform.Func {
	return func(out transform.Writer) error {
		if len(constants) == 0 {
			return nil
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

		return nil
	}
}

func Enum(enum *ast.Enum, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) transform.Func {
	return func(out transform.Writer) error {
		enumType, err := parseType(enum.Type, messagesMap, enumsMap)
		if err != nil {
			return err
		}

		out.
			String("type ").
			String(enum.Name.Name).
			String(" ").
			String(enumType).
			NewLines(2)

		err = Constants(enum.Constants)(out)
		if err != nil {
			return err
		}

		return nil
	}
}

func Enums(enums []*ast.Enum, messagesMap map[string]*ast.Message, enumsMap map[string]*ast.Enum) transform.Func {
	return func(out transform.Writer) error {
		if len(enums) == 0 {
			return nil
		}

		for i, enum := range enums {
			if i != 0 {
				out.NewLines(2)
			}
			err := Enum(enum, messagesMap, enumsMap)(out)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
