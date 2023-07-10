package code

import (
	"fmt"
	"io/fs"
	"text/template"

	"ella.to/internal/ast"
)

func LoadTemplate(files fs.FS, name string) (*template.Template, error) {
	file := fmt.Sprintf("%s.tmpl", name)
	tmpl := template.New(file).Funcs(DefaultFuncsMap)
	return tmpl.ParseFS(files, fmt.Sprintf("tmpl/%s", file))
}

func getContent[T any](node ast.Node) []T {
	var results []T

	prog, ok := node.(*ast.Program)
	if !ok {
		return results
	}

	for _, node := range prog.Nodes {
		result, ok := node.(T)
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

func GetMessages(node ast.Node) []*ast.Message {
	return getContent[*ast.Message](node)
}

func GetServices(node ast.Node) []*ast.Service {
	return getContent[*ast.Service](node)
}
