package validator

import (
	"fmt"
	"sort"

	"github.com/alinz/ella.to/schema/ast"
)

const (
	ellaVersion = `"0.0.1"`
)

func validateConstants(constantsMap map[string]*ast.Constant) ([]*ast.Constant, error) {
	constants := make([]*ast.Constant, 0, len(constantsMap))

	value, ok := constantsMap["ella"]
	if !ok {
		return nil, fmt.Errorf("rpc constant is not defined")
	}

	if value.Value.TokenLiteral() != ellaVersion {
		return nil, fmt.Errorf("rpc version is not supported")
	}

	for _, constant := range constantsMap {
		if constant.Name.Name == "ella" {
			continue
		}

		constants = append(constants, constant)
	}

	sort.Slice(constants, func(i, j int) bool {
		return constants[i].Name.Name < constants[j].Name.Name
	})

	// add ellaVersion to the beginning of the list
	constants = append([]*ast.Constant{value}, constants...)

	return constants, nil
}
