package parser

import (
	"fmt"

	"github.com/alinz/ella.to/schema/ast"
	"github.com/alinz/ella.to/schema/token"
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

	err = mustBeNameFor(p.nextToken, "enum", true)
	if err != nil {
		return nil, err
	}

	enum := &ast.Enum{}

	enum.Name = &ast.Identifier{
		Name:  p.nextToken.Val,
		Token: p.nextToken,
	}

	p.scanToken()

	enum.Type, err = parseEnumType(p.nextToken)
	if err != nil {
		return nil, err
	}

	p.scanToken()

	if p.nextToken.Kind != token.OpenCurly {
		return nil, fmt.Errorf("expected '{' for enum's body start but got %s", p.nextToken.Kind)
	}

	p.scanToken() // skip '{'

	for p.nextToken.Kind != token.CloseCurly {
		constant, err := p.parseConstant(true)
		if err != nil {
			return nil, err
		}

		// we need convert the value from Int to Uint
		// this involves checking if value is not negative
		// this the only place where we can do this check, there is no way for us to
		// detect uint and int values without knowing the type
		if _, ok := enum.Type.(*ast.TypeUint); ok {
			if v, ok := constant.Value.(*ast.ValueInt); ok {
				if v.Content < 0 {
					return nil, fmt.Errorf("enum %s type is unsiged int and one of its value is negative", enum.Name.Name)
				}

				constant.Value = &ast.ValueUint{
					Content: uint64(v.Content),
					Token:   v.Token,
					Size:    v.Size,
				}
			}
		}

		enum.Constants = append(enum.Constants, constant)
	}

	p.scanToken() // skip '}'

	return enum, nil
}

func parseEnumType(token *token.Token) (ast.Type, error) {
	switch token.Val {
	case "int8":
		return &ast.TypeInt{
			Token: token,
			Size:  8,
		}, nil
	case "int16":
		return &ast.TypeInt{
			Token: token,
			Size:  16,
		}, nil
	case "int32":
		return &ast.TypeInt{
			Token: token,
			Size:  32,
		}, nil
	case "int64":
		return &ast.TypeInt{
			Token: token,
			Size:  64,
		}, nil
	case "uint8":
		return &ast.TypeUint{
			Token: token,
			Size:  8,
		}, nil
	case "uint16":
		return &ast.TypeUint{
			Token: token,
			Size:  16,
		}, nil
	case "uint32":
		return &ast.TypeUint{
			Token: token,
			Size:  32,
		}, nil
	case "uint64":
		return &ast.TypeUint{
			Token: token,
			Size:  64,
		}, nil
	default:
		return nil, fmt.Errorf("only int8, int16, int32, int64, uint8, uint16, uint32 and uint64 are supported for enum type but got %s", token.Val)
	}
}
