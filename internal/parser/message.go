package parser

import (
	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func ParseMessage(p *Parser) (*ast.Message, error) {
	if p.Peek().Type != token.Message {
		return nil, p.WithError(p.Peek(), "expected 'message' keyword")
	}

	message := &ast.Message{Token: p.Next()}

	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a message")
	}

	message.Name = &ast.Identifier{Token: p.Next()}

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

	field = &ast.Field{
		Name:    &ast.Identifier{Token: p.Next()},
		Options: make([]*ast.Const, 0),
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

	p.Next() // skip '{'

	for p.Peek().Type != token.CloseCurly {
		constant, err := ParseMessageFieldConstant(p)
		if err != nil {
			return nil, err
		}

		field.Options = append(field.Options, constant)
	}

	p.Next() // skip '}'

	return field, nil
}

func ParseMessageFieldConstant(p *Parser) (constant *ast.Const, err error) {
	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining a message field constant")
	}

	constant = &ast.Const{
		Name: &ast.Identifier{Token: p.Next()},
	}

	if p.Peek().Type != token.Colon {
		constant.Value = &ast.ValueBool{
			Token:   nil,
			Value:   true,
			Defined: false,
		}

		return constant, nil
	}

	p.Next() // skip ':'

	constant.Value, err = ParseValue(p)
	if err != nil {
		return nil, err
	}

	return constant, nil
}
