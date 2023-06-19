package parser

import (
	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func newComment() *ast.Comment {
	return &ast.Comment{
		Tops: []*token.Token{},
	}
}

func attachComment(node ast.Node, comment *ast.Comment) {
	if comment == nil {
		return
	}

	switch v := node.(type) {
	case *ast.Const:
		v.Comment = comment
	case *ast.Enum:
		v.Comment = comment
	case *ast.Message:
		v.Comment = comment
	case *ast.Service:
		v.Comment = comment
	}
}
