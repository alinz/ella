package typescript

import (
	"io"

	"ella.to/internal/ast"
	"ella.to/internal/code"
)

func New() code.Generator {
	return code.GeneratorFunc(func(w io.Writer, pkg string, prog *ast.Program) error {
		return nil
	})
}
