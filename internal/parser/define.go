package parser

import (
	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func ParseDefine(p *Parser) (*ast.Define, error) {
	if p.Peek().Type != token.Const {
		return nil, p.WithError(p.Peek(), "expected 'define' keyword")
	}

	tok := p.Next()

	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected an identifier after define")
	}

	name := &ast.Identifier{
		Token: p.Next(),
	}

	// if p.Peek().Type != token.Identifier
	typ, err := ParseType(p)
	if err != nil {
		return nil, err
	}

	if p.Peek().Type != token.OpenCurly {
		return &ast.Define{
			Token:   tok,
			Name:    name,
			Type:    typ,
			Options: make([]*ast.Option, 0),
		}, nil
	}

	options, err := ParseOptions(p)
	if err != nil {
		return nil, err
	}

	return &ast.Define{
		Token:   tok,
		Name:    name,
		Type:    typ,
		Options: options,
	}, nil
}
