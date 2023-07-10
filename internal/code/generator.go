package code

import (
	"io"

	"ella.to/internal/ast"
)

type Generator interface {
	Generate(w io.Writer, pkg string, prog *ast.Program) error
}

type GeneratorFunc func(w io.Writer, pkg string, prog *ast.Program) error

func (f GeneratorFunc) Generate(w io.Writer, pkg string, prog *ast.Program) error {
	return f(w, pkg, prog)
}
