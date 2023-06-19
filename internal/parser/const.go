package parser

import (
	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func (p *Parser) parseConst(permitEmpty bool) (*ast.Const, error) {
	nameTok := p.currTok

	p.scanToken() // skip identifier

	if permitEmpty && p.currTok.Type != token.Assign {
		return &ast.Const{
			Name: &ast.Identifier{Token: nameTok},
			Value: &ast.ValueBool{
				Token: p.currTok,
				Value: true,
			},
		}, nil
	} else if !permitEmpty && p.currTok.Type != token.Assign {
		return nil, p.newError(p.currTok, "expected '=' after an identifier for defining a constant")
	}

	p.scanToken() // skip '='

	value, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	return &ast.Const{
		Name:  &ast.Identifier{Token: nameTok},
		Value: value,
	}, nil
}
