package parser

import (
	"compiler.ella.to/internal/ast"
	"compiler.ella.to/internal/token"
	"compiler.ella.to/pkg/strcase"
)

func ParseModel(p *Parser) (*ast.Model, error) {
	if p.Peek().Type != token.Model {
		return nil, p.WithError(p.Peek(), "expected 'message' keyword")
	}

	message := &ast.Model{Token: p.Next()}

	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a message")
	}

	nameTok := p.Next()

	if !strcase.IsPascal(nameTok.Literal) {
		return nil, p.WithError(nameTok, "message name must be in PascalCase format")
	}

	message.Name = &ast.Identifier{Token: nameTok}

	if p.Peek().Type != token.OpenCurly {
		return nil, p.WithError(p.Peek(), "expected '{' after message declaration")
	}

	p.Next() // skip '{'

	for p.Peek().Type != token.CloseCurly {
		if p.Peek().Type == token.Extend {
			extend, err := ParseExtend(p)
			if err != nil {
				return nil, err
			}
			message.Extends = append(message.Extends, extend)
		} else {
			field, err := ParseModelField(p)
			if err != nil {
				return nil, err
			}
			message.Fields = append(message.Fields, field)
		}
	}

	p.Next() // skip '}'

	return message, nil
}

func ParseExtend(p *Parser) (*ast.Identifier, error) {
	if p.Peek().Type != token.Extend {
		return nil, p.WithError(p.Peek(), "expected '...' keyword")
	}

	p.Next() // skip '...'

	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for extending a message")
	}

	nameTok := p.Next()

	if !strcase.IsPascal(nameTok.Literal) {
		return nil, p.WithError(nameTok, "extend message name must be in PascalCase format")
	}

	return &ast.Identifier{Token: nameTok}, nil
}

func ParseModelField(p *Parser) (field *ast.Field, err error) {
	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a message field")
	}

	nameTok := p.Next()

	if !strcase.IsPascal(nameTok.Literal) {
		return nil, p.WithError(nameTok, "message field name must be in PascalCase format")
	}

	field = &ast.Field{
		Name:    &ast.Identifier{Token: nameTok},
		Options: make([]*ast.Option, 0),
	}

	if p.Peek().Type != token.Colon {
		return nil, p.WithError(p.Peek(), "expected ':' after message field name")
	}

	p.Next() // skip ':'

	field.Type, err = ParseType(p)
	if err != nil {
		return nil, err
	}

	if p.Peek().Type != token.OpenCurly {
		return field, nil
	}

	field.Options, err = ParseOptions(p)
	if err != nil {
		return nil, err
	}

	return field, nil
}
