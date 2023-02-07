package parser

import (
	"fmt"

	"github.com/alinz/rpc.go/schema/ast"
	"github.com/alinz/rpc.go/schema/token"
)

//	enum A int64 {
//		A = 1
//	}
//
// T -> Enum ID Type {}
// T -> Enum ID Type { C }
// T -> Enum ID Type { C = 2 }
//
//	T -> Enum ID Type {
//		C = 2
//		C
//	}
func (p *Parser) parseEnum() (*ast.Enum, error) {
	var err error

	if p.nextToken.Kind != token.Enum {
		return nil, fmt.Errorf("expected 'enum' but got %s", p.nextToken.Kind)
	}

	p.scanToken() // skip 'enum'

	if p.nextToken.Kind != token.Identifier {
		return nil, fmt.Errorf("expected identifier but got %s", p.nextToken.Kind)
	}

	enum := &ast.Enum{}

	enum.Name = &ast.Identifier{
		Name:  p.nextToken.Val,
		Token: p.nextToken,
	}

	p.scanToken()

	if p.nextToken.Kind != token.Type {
		return nil, fmt.Errorf("expected type but got %s", p.nextToken.Kind)
	}

	enum.Type, err = parseEnumType(p.nextToken)
	if err != nil {
		return nil, err
	}

	p.scanToken()

	if p.nextToken.Kind != token.OpenCurl {
		return nil, fmt.Errorf("expected '{' but got %s", p.nextToken.Kind)
	}

	p.scanToken() // skip '{'

	for p.nextToken.Kind != token.CloseCurl {
		constant, err := p.parseConstant(true)
		if err != nil {
			return nil, err
		}

		enum.Constants = append(enum.Constants, constant)
	}

	p.scanToken() // skip '}'

	return enum, nil
}

func parseEnumType(token *token.Token) (ast.Type, error) {
	switch token.Val {
	case "Int8":
		return &ast.TypeInt{
			Token: token,
			Size:  8,
		}, nil
	case "Int16":
		return &ast.TypeInt{
			Token: token,
			Size:  16,
		}, nil
	case "Int32":
		return &ast.TypeInt{
			Token: token,
			Size:  32,
		}, nil
	case "Int64":
		return &ast.TypeInt{
			Token: token,
			Size:  64,
		}, nil
	case "Uint8":
		return &ast.TypeUint{
			Token: token,
			Size:  8,
		}, nil
	case "Uint16":
		return &ast.TypeUint{
			Token: token,
			Size:  16,
		}, nil
	case "Uint32":
		return &ast.TypeUint{
			Token: token,
			Size:  32,
		}, nil
	case "Uint64":
		return &ast.TypeUint{
			Token: token,
			Size:  64,
		}, nil
	default:
		return nil, fmt.Errorf("only Int8, Int16, Int32, Int64, Uint8, Uint16, Uint32 and Uint64 are supported for enum type but got %s", token.Val)
	}
}
