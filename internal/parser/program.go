package parser

import (
	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func ParseProgram(p *Parser) (prog *ast.Program, err error) {
	prog = &ast.Program{
		Nodes: make([]ast.Node, 0),
	}

	for p.Peek().Type != token.EOF {
		var node ast.Node

		switch p.Peek().Type {
		case token.Identifier:
			node, err = ParseConst(p)
		case token.Enum:
			node, err = ParseEnum(p)
		case token.Message:
			node, err = ParseMessage(p)
		case token.Service:
			node, err = ParseService(p)
		default:
			return nil, p.WithError(p.Peek(), "unexpected token")
		}

		if err != nil {
			return nil, err
		}

		prog.Nodes = append(prog.Nodes, node)
	}

	return prog, nil
}
