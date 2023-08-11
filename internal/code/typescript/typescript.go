package typescript

import (
	"ella.to/internal/ast"
	"ella.to/internal/code"
)

func New() code.Generator {
	return code.GeneratorFunc(func(out string, prog *ast.Program) error {
		return nil
	})
}
