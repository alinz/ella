package validator

import (
	"fmt"
	"sort"

	"github.com/alinz/rpc.go/schema/ast"
)

const (
	rpcVersion = "0.0.1"
)

func validateConstants(constantsMap map[string]*ast.Constant) ([]*ast.Constant, error) {
	constants := make([]*ast.Constant, 0, len(constantsMap))

	rpcConstant, ok := constantsMap["rpc"]
	if !ok {
		return nil, fmt.Errorf("rpc constant is not defined")
	}

	if rpcConstant.Value.TokenLiteral() != rpcVersion {
		return nil, fmt.Errorf("rpc version is not supported")
	}

	for _, constant := range constantsMap {
		if constant.Name.Name == "rpc" {
			continue
		}

		constants = append(constants, constant)
	}

	sort.Slice(constants, func(i, j int) bool {
		return constants[i].Name.Name < constants[j].Name.Name
	})

	// add rpcConstant to the beginning of the list
	constants = append([]*ast.Constant{rpcConstant}, constants...)

	return constants, nil
}
