package parser

import (
	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func ParseProgram(p *Parser) (prog *ast.Program, err error) {
	prog = &ast.Program{
		Statements: make([]ast.Statement, 0),
	}

	for p.Peek().Type != token.EOF {
		var stmt ast.Statement

		switch p.Peek().Type {
		case token.Identifier:
			stmt, err = ParseConst(p)
		case token.Enum:
			stmt, err = ParseEnum(p)
		case token.Message:
			stmt, err = ParseMessage(p)
		case token.Service:
			stmt, err = ParseService(p)
		default:
			return nil, p.WithError(p.Peek(), "unexpected token")
		}

		if err != nil {
			return nil, err
		}

		prog.Statements = append(prog.Statements, stmt)
	}

	return prog, nil
}
