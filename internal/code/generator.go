package code

import (
	"ella.to/internal/ast"
)

type Generator interface {
	Generate(outFilename string, prog *ast.Program) error
}

type GeneratorFunc func(outFilename string, prog *ast.Program) error

func (f GeneratorFunc) Generate(outFilename string, prog *ast.Program) error {
	return f(outFilename, prog)
}

func RunParsers(prog *ast.Program, fns ...func(prog *ast.Program) error) error {
	for _, fn := range fns {
		if err := fn(prog); err != nil {
			return err
		}
	}
	return nil
}
