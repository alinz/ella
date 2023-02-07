package parser

import (
	"fmt"

	"github.com/alinz/rpc.go/schema/ast"
	"github.com/alinz/rpc.go/schema/token"
)

func (p *Parser) parseMessage() (message *ast.Message, err error) {
	if p.nextToken.Kind != token.Message {
		return nil, fmt.Errorf("expected message keyword")
	}

	p.scanToken() // skip message

	if p.nextToken.Kind != token.Identifier {
		return nil, fmt.Errorf("expected a name for message but got %s", p.nextToken.Val)
	}

	message = &ast.Message{
		Name: &ast.Identifier{
			Name:  p.nextToken.Val,
			Token: p.nextToken,
		},
	}

	p.scanToken() // skip name

	if p.nextToken.Kind != token.OpenCurl {
		return nil, fmt.Errorf("expected { for message but got %s", p.nextToken.Val)
	}

	p.scanToken() // skip {

	for p.nextToken.Kind != token.CloseCurl {

		if p.nextToken.Kind == token.Ellipsis {
			p.scanToken()

			if p.nextToken.Kind != token.Type {
				return nil, fmt.Errorf("expected a name for extends but got %s", p.nextToken.Val)
			}

			message.Extends = append(message.Extends, &ast.TypeCustom{
				Token: p.nextToken,
				Name:  p.nextToken.Val,
			})

			p.scanToken()

			continue
		}

		field, err := p.parseField()
		if err != nil {
			return nil, err
		}

		message.Fields = append(message.Fields, field)
	}

	p.scanToken()

	return message, nil
}

func (p *Parser) parseField() (field *ast.Field, err error) {
	if p.nextToken.Kind != token.Identifier {
		return nil, fmt.Errorf("expected a name for field but got %s", p.nextToken.Val)
	}

	field = &ast.Field{
		Name: &ast.Identifier{
			Name:  p.nextToken.Val,
			Token: p.nextToken,
		},
		Options: make([]*ast.Constant, 0),
	}

	p.scanToken()

	if p.nextToken.Kind != token.Colon {
		return nil, fmt.Errorf("expected : for field but got %s", p.nextToken.Val)
	}

	p.scanToken()

	field.Type, err = p.parseType()
	if err != nil {
		return nil, err
	}

	if p.nextToken.Kind == token.OpenCurl {
		field.Options, err = p.parseFieldOptions()
		if err != nil {
			return nil, err
		}
	}

	return field, nil
}

func (p *Parser) parseFieldOptions() ([]*ast.Constant, error) {
	if p.nextToken.Kind != token.OpenCurl {
		return nil, fmt.Errorf("expected { for field options but got %s", p.nextToken.Val)
	}

	p.scanToken() // skip {

	options := make([]*ast.Constant, 0)

	for p.nextToken.Kind != token.CloseCurl {
		option, err := p.parseConstant(true)
		if err != nil {
			return nil, err
		}

		options = append(options, option)
	}

	return nil, nil
}
