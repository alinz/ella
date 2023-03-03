package parser

import (
	"fmt"

	"github.com/alinz/ella.to/schema/ast"
	"github.com/alinz/ella.to/schema/token"
)

func (p *Parser) parseMessage() (message *ast.Message, err error) {
	if p.nextToken.Kind != token.Message {
		return nil, fmt.Errorf("expected message token but got %s", p.nextToken.Kind)
	}

	p.scanToken() // skip message

	err = mustBeNameFor(p.nextToken, "message", true)
	if err != nil {
		return nil, err
	}

	message = &ast.Message{
		Name: &ast.Identifier{
			Name:  p.nextToken.Val,
			Token: p.nextToken,
		},
	}

	p.scanToken() // skip name

	if p.nextToken.Kind != token.OpenCurly {
		return nil, fmt.Errorf("expected { for message's body start but got %s", p.nextToken.Val)
	}

	p.scanToken() // skip {

	for p.nextToken.Kind != token.CloseCurly {
		if p.nextToken.Kind == token.Dot { // check for extends
			p.scanToken() // skip first .
			if p.nextToken.Kind != token.Dot {
				return nil, fmt.Errorf("expect .. after the first . in message %s", message.Name.Name)
			}
			p.scanToken() // skip second .
			if p.nextToken.Kind != token.Dot {
				return nil, fmt.Errorf("expect . after .. in message %s", message.Name.Name)
			}
			p.scanToken() // skip third .

			if p.nextToken.Kind != token.Word || !isIdentifier(p.nextToken.Val, true) {
				return nil, fmt.Errorf("expected a valid name for extends but got %s", p.nextToken.Val)
			}

			message.Extends = append(message.Extends, &ast.TypeCustom{
				Token: p.nextToken,
				Name:  p.nextToken.Val,
			})

			p.scanToken() // skip name

			continue
		}

		field, err := p.parseField()
		if err != nil {
			return nil, err
		}

		message.Fields = append(message.Fields, field)
	}

	p.scanToken() // skip }

	return message, nil
}

func (p *Parser) parseField() (field *ast.Field, err error) {
	if p.nextToken.Kind != token.Word || !isIdentifier(p.nextToken.Val, true) {
		return nil, fmt.Errorf("expected a valid name for field name but got %s", p.nextToken.Val)
	}

	field = &ast.Field{
		Name: &ast.Identifier{
			Name:  p.nextToken.Val,
			Token: p.nextToken,
		},
		Options: make([]*ast.Constant, 0),
	}

	p.scanToken() // skip name

	if p.nextToken.Kind != token.Colon {
		return nil, fmt.Errorf("expected : after the field's name but got %s", p.nextToken.Val)
	}

	p.scanToken() // skip :

	field.Type, err = p.parseType()
	if err != nil {
		return nil, err
	}

	if p.nextToken.Kind == token.OpenCurly {
		field.Options, err = p.parseFieldOptions()
		if err != nil {
			return nil, err
		}
	}

	return field, nil
}

func (p *Parser) parseFieldOptions() ([]*ast.Constant, error) {
	if p.nextToken.Kind != token.OpenCurly {
		return nil, fmt.Errorf("expected { for defining field options but got %s", p.nextToken.Val)
	}

	p.scanToken() // skip {

	options := make([]*ast.Constant, 0)

	for p.nextToken.Kind != token.CloseCurly {
		option, err := p.parseConstant(true)
		if err != nil {
			return nil, err
		}

		options = append(options, option)
	}

	p.scanToken() // skip }

	return options, nil
}
