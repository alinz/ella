package parser

import (
	"strconv"

	"ella.to/internal/ast"
	"ella.to/internal/token"
)

var enumTypes = []token.Type{
	token.Int8,
	token.Int16,
	token.Int32,
	token.Int64,
	token.Uint8,
	token.Uint16,
	token.Uint32,
	token.Uint64,
}

func ParseEnum(p *Parser) (enum *ast.Enum, err error) {
	if p.Peek().Type != token.Enum {
		return nil, p.WithError(p.Peek(), "expected 'enum' keyword")
	}

	enum = &ast.Enum{Token: p.Next()}

	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining an enum")
	}

	enum.Name = &ast.Identifier{Token: p.Next()}

	enum.Type, err = ParseEnumType(p)
	if err != nil {
		return nil, err
	}

	if p.Peek().Type != token.OpenCurly {
		return nil, p.WithError(p.Peek(), "expected '{' after enum declaration")
	}

	p.Next() // skip '{'

	for p.Peek().Type != token.CloseCurly {
		constant, err := ParseEnumConstant(p)
		if err != nil {
			return nil, err
		}

		enum.Constants = append(enum.Constants, constant)
	}

	p.Next() // skip '}'

	return enum, nil
}

func ParseEnumType(p *Parser) (ast.Type, error) {
	peek := p.Peek()

	if !token.OneOfTypes(peek, enumTypes...) {
		return nil, p.WithError(peek, "expected enum type")
	}

	switch peek.Type {
	case token.Int8:
		return &ast.Int{Token: p.Next(), Size: 8}, nil
	case token.Int16:
		return &ast.Int{Token: p.Next(), Size: 16}, nil
	case token.Int32:
		return &ast.Int{Token: p.Next(), Size: 32}, nil
	case token.Int64:
		return &ast.Int{Token: p.Next(), Size: 64}, nil
	case token.Uint8:
		return &ast.Uint{Token: p.Next(), Size: 8}, nil
	case token.Uint16:
		return &ast.Uint{Token: p.Next(), Size: 16}, nil
	case token.Uint32:
		return &ast.Uint{Token: p.Next(), Size: 32}, nil
	case token.Uint64:
		return &ast.Uint{Token: p.Next(), Size: 64}, nil
	default:
		panic("unreachable") // this should not happen as we already checked for the types
	}
}

func ParseEnumConstant(p *Parser) (*ast.Const, error) {
	if p.Peek().Type != token.Identifier {
		return nil, p.WithError(p.Peek(), "expected identifier for defining an enum constant")
	}

	nameTok := p.Next()

	if p.Peek().Type != token.Assign {
		return &ast.Const{
			Name:  &ast.Identifier{Token: nameTok},
			Value: &ast.ValueInt{},
		}, nil
	}

	p.Next() // skip '='

	if p.Peek().Type != token.ConstInt {
		return nil, p.WithError(p.Peek(), "expected constant integer value for defining an enum constant value")
	}

	valueTok := p.Next()

	value, err := strconv.ParseInt(valueTok.Val, 10, 64)
	if err != nil {
		return nil, p.WithError(valueTok, "invalid integer value for defining an enum constant value", err)
	}

	return &ast.Const{
		Name: &ast.Identifier{Token: nameTok},
		Value: &ast.ValueInt{
			Token:   valueTok,
			Value:   value,
			Defined: true,
		},
	}, nil
}
