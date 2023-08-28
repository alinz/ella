package validator

import (
	"ella.to/internal/ast"
	"ella.to/internal/astutil"
)

func validateBases(prog *ast.Program) error {
	bases := astutil.GetBases(prog)
	constantsMap := astutil.CreateConstsMap(prog)

	for _, base := range bases {
		for _, option := range base.Options {
			value, err := getValueRef(option.Value, constantsMap)
			if err != nil {
				return err
			}

			option.Value = value
		}
	}

	return nil
}
