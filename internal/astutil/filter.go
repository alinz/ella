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

func GetBases(node ast.Node) []*ast.Base {
	return getContent[*ast.Base](node)
}

func GetConstants(node ast.Node) []*ast.Const {
	return getContent[*ast.Const](node)
}

func GetEnums(node ast.Node) []*ast.Enum {
	return getContent[*ast.Enum](node)
}

func GetMessages(node ast.Node) []*ast.Message {
	return getContent[*ast.Message](node)
}

func GetServices(node ast.Node) []*ast.Service {
	return getContent[*ast.Service](node)
}
