package validator

import (
	"fmt"

	"compiler.ella.to/internal/ast"
)

func Validate(prog *ast.Program) error {
	return runValidators(
		prog,
		validateUniqueNames,
		validateModels,
	)
}

func validateUniqueNames(prog *ast.Program) error {
	names := make(map[string]struct{})
	for _, stmt := range prog.Statements {
		var name string
		switch stmt := stmt.(type) {
		case *ast.Const:
			name = stmt.Name.String()
			if _, ok := names[name]; ok {
				return fmt.Errorf("const %s is defined multiple times", stmt.Name)
			}
		case *ast.Enum:
			name = stmt.Name.String()
			if _, ok := names[name]; ok {
				return fmt.Errorf("enum %s is defined multiple times", stmt.Name)
			}
		case *ast.Model:
			name = stmt.Name.String()
			if _, ok := names[name]; ok {
				return fmt.Errorf("message %s is defined multiple times", stmt.Name)
			}
		case *ast.Service:
			name = stmt.Name.String()
			if _, ok := names[name]; ok {
				return fmt.Errorf("service %s is defined multiple times", stmt.Name)
			}
		}
		names[name] = struct{}{}
	}

	return nil
}

type ValidatorFunc func(prog *ast.Program) error

func runValidators(prog *ast.Program, validatorFuncs ...ValidatorFunc) error {
	for _, validatorFunc := range validatorFuncs {
		if err := validatorFunc(prog); err != nil {
			return err
		}
	}

	return nil
}
