package parser

import (
	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func ParseAlias(p *Parser) (*ast.Alias, error) {
	if p.Peek().Type != token.Alias {
		return nil, p.WithError(p.Peek(), "expected 'alias' keyword")
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
		return &ast.Alias{
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

	return &ast.Alias{
		Token:   tok,
		Name:    name,
		Type:    typ,
		Options: options,
	}, nil
}
