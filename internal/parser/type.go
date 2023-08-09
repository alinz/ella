package parser

import (
	"strconv"

	"ella.to/internal/ast"
	"ella.to/internal/token"
	"ella.to/pkg/strcase"
)

func ParseType(p *Parser) (ast.Type, error) {
	switch p.Peek().Type {
	case token.Map:
		return ParseMapType(p)
	case token.Array:
		return ParseArrayType(p)
	case token.Bool:
		return &ast.Bool{Token: p.Next()}, nil
	case token.Byte:
		return &ast.Byte{Token: p.Next()}, nil
	case token.Int8, token.Int16, token.Int32, token.Int64:
		tok := p.Next()
		return &ast.Int{
			Token: tok,
			Size:  extractTypeBits("int", tok.Val),
		}, nil
	case token.Uint8, token.Uint16, token.Uint32, token.Uint64:
		tok := p.Next()
		return &ast.Uint{
			Token: tok,
			Size:  extractTypeBits("uint", tok.Val),
		}, nil
	case token.Float32, token.Float64:
		tok := p.Next()
		return &ast.Float{
			Token: tok,
			Size:  extractTypeBits("float", tok.Val),
		}, nil
	case token.Timestamp:
		return &ast.Timestamp{Token: p.Next()}, nil
	case token.String:
		return &ast.String{Token: p.Next()}, nil
	case token.Any:
		return &ast.Any{Token: p.Next()}, nil
	case token.File:
		return &ast.File{Token: p.Next()}, nil
	case token.Identifier:
		nameTok := p.Next()

		if !strcase.IsPascal(nameTok.Val) {
			return nil, p.WithError(nameTok, "custom type name must be in PascalCase format")
		}

		return &ast.CustomType{Token: nameTok}, nil
	default:
		return nil, p.WithError(p.Peek(), "expected type")
	}
}

func ParseMapType(p *Parser) (*ast.Map, error) {
	if p.Peek().Type != token.Map {
		return nil, p.WithError(p.Peek(), "expected 'map' keyword")
	}

	mapTok := p.Next()

	if p.Peek().Type != token.OpenAngle {
		return nil, p.WithError(p.Peek(), "expected '<' after 'map' keyword")
	}

	p.Next() // skip '<'

	keyType, err := ParseMapKeyType(p)
	if err != nil {
		return nil, err
	}

	if p.Peek().Type != token.Comma {
		return nil, p.WithError(p.Peek(), "expected ',' after map key type")
	}

	p.Next() // skip ','

	valueType, err := ParseType(p)
	if err != nil {
		return nil, err
	}

	if p.Peek().Type != token.CloseAngle {
		return nil, p.WithError(p.Peek(), "expected '>' after map value type")
	}

	p.Next() // skip '>'

	return &ast.Map{
		Token: mapTok,
		Key:   keyType,
		Value: valueType,
	}, nil
}

func ParseMapKeyType(p *Parser) (ast.Type, error) {
	switch p.Peek().Type {
	case token.Int8, token.Int16, token.Int32, token.Int64:
		return ParseType(p)
	case token.Uint8, token.Uint16, token.Uint32, token.Uint64:
		return ParseType(p)
	case token.String:
		return ParseType(p)
	case token.Byte:
		return ParseType(p)
	default:
		return nil, p.WithError(p.Peek(), "expected map key type to be comparable")
	}
}

func ParseArrayType(p *Parser) (*ast.Array, error) {
	if p.Peek().Type != token.Array {
		return nil, p.WithError(p.Peek(), "expected 'array' keyword")
	}

	arrayTok := p.Next()

	arrayType, err := ParseType(p)
	if err != nil {
		return nil, err
	}

	return &ast.Array{
		Token: arrayTok,
		Type:  arrayType,
	}, nil
}

func extractTypeBits(prefix string, value string) int {
	// The resason why we don't return an error here is because
	// scanner already give us int8 ... float64 values and it has already
	// been validated.
	result, _ := strconv.ParseInt(value[len(prefix):], 10, 64)
	return int(result)
}
