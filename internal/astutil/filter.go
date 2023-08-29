package astutil

import "ella.to/internal/ast"

func getContent[T any](node ast.Node) []T {
	var results []T

	prog, ok := node.(*ast.Program)
	if !ok {
		return results
	}

	for _, stmt := range prog.Statements {
		result, ok := stmt.(T)
		if !ok {
			continue
		}

		results = append(results, result)
	}

	return results
}

func GetConstants(node ast.Node) []*ast.Const {
	return getContent[*ast.Const](node)
}

func GetEnums(node ast.Node) []*ast.Enum {
	return getContent[*ast.Enum](node)
}

func GetModels(node ast.Node) []*ast.Model {
	return getContent[*ast.Model](node)
}

func GetServices(node ast.Node) []*ast.Service {
	return getContent[*ast.Service](node)
}
