package parser

import (
	"ella.to/internal/ast"
	"ella.to/internal/token"
	"ella.to/pkg/strcase"
)

func ParseConst(p *Parser) (*ast.Const, error) {
	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a constant")
	}

	nameTok := p.Next()

	if !strcase.IsPascal(nameTok.Val) {
		return nil, p.WithError(nameTok, "constant name must be in PascalCase format")
	}

	if p.Peek().Type != token.Assign {
		return nil, p.WithError(p.Peek(), "expected '=' after an identifier for defining a constant")
	}

	p.Next() // skip '='

	value, err := ParseValue(p)
	if err != nil {
		return nil, err
	}

	return &ast.Const{
		Name:  &ast.Identifier{Token: nameTok},
		Value: value,
	}, nil
}
