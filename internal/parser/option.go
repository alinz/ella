package parser

import (
	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func ParseOption(p *Parser) (option *ast.Option, err error) {
	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a message field option")
	}

	nameTok := p.Next()

	option = &ast.Option{
		Name: &ast.Identifier{Token: nameTok},
	}

	if p.Peek().Type != token.Assign {
		option.Value = &ast.ValueBool{
			Token:   nil,
			Value:   true,
			Defined: false,
		}

		return option, nil
	}

	p.Next() // skip '='

	option.Value, err = ParseValue(p)
	if err != nil {
		return nil, err
	}

	return option, nil
}

func ParseOptions(p *Parser) (ast.Options, error) {
	options := make([]*ast.Option, 0)

	p.Next() // skip '{'

	for p.Peek().Type != token.CloseCurly {
		option, err := ParseOption(p)
		if err != nil {
			return nil, err
		}

		options = append(options, option)
	}

	p.Next() // skip '}'

	return options, nil
}
