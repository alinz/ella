package parser

import (
	"ella.to/internal/ast"
	"ella.to/internal/token"
	"ella.to/pkg/strcase"
)

func ParseMessage(p *Parser) (*ast.Message, error) {
	if p.Peek().Type != token.Message {
		return nil, p.WithError(p.Peek(), "expected 'message' keyword")
	}

	message := &ast.Message{Token: p.Next()}

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
		field, err := ParseMessageField(p)
		if err != nil {
			return nil, err
		}

		message.Fields = append(message.Fields, field)
	}

	p.Next() // skip '}'

	return message, nil
}

func ParseMessageField(p *Parser) (field *ast.Field, err error) {
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
