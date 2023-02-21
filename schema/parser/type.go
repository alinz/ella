package parser

import (
	"fmt"
	"strconv"

	"github.com/alinz/rpc.go/schema/ast"
	"github.com/alinz/rpc.go/schema/token"
)

// T -> map < key , T >
// T -> [ ] T
// T -> int, string, ...

func (p *Parser) parseType() (typ ast.Type, err error) {
	switch p.nextToken.Kind {
	case token.Type:
		switch p.nextToken.Val {
		case "map":
			return p.parseMapType()
		case "Int8", "Int16", "Int32", "Int64":
			typ = &ast.TypeInt{
				Token: p.nextToken,
				Size:  extractBits("Int", p.nextToken.Val),
			}
		case "Uint8", "Uint16", "Uint32", "Uint64":
			typ = &ast.TypeUint{
				Token: p.nextToken,
				Size:  extractBits("Uint", p.nextToken.Val),
			}
		case "Float32", "Float64":
			typ = &ast.TypeFloat{
				Token: p.nextToken,
				Size:  extractBits("Float", p.nextToken.Val),
			}
		case "String":
			typ = &ast.TypeString{
				Token: p.nextToken,
			}
		case "Bool":
			typ = &ast.TypeBool{
				Token: p.nextToken,
			}
		case "Byte":
			typ = &ast.TypeByte{
				Token: p.nextToken,
			}
		case "Timestamp":
			typ = &ast.TypeTimestamp{
				Token: p.nextToken,
			}
		case "Any":
			typ = &ast.TypeAny{
				Token: p.nextToken,
			}
		default:
			typ = &ast.TypeCustom{
				Token: p.nextToken,
				Name:  p.nextToken.Val,
			}
		}
		p.scanToken()
		return typ, nil
	case token.OpenBracket:
		return p.parseArrayType()
	}

	return nil, fmt.Errorf("expected type but got %s", p.nextToken.Val)
}

func (p *Parser) parseMapType() (typ *ast.TypeMap, err error) {
	typ = &ast.TypeMap{}

	if p.nextToken.Kind != token.Type || p.nextToken.Val != "map" {
		return nil, fmt.Errorf("expected map keyword but got %s", p.nextToken.Val)
	}

	typ.Token = p.nextToken

	p.scanToken() // skip map

	if p.nextToken.Kind != token.OpenAngle {
		return nil, fmt.Errorf("expected < but got %s", p.nextToken.Val)
	}

	p.scanToken() // skip <

	typ.Key, err = p.parseKeyType()
	if err != nil {
		return nil, err
	}

	if p.nextToken.Kind != token.Comma {
		return nil, fmt.Errorf("expected , but got %s", p.nextToken.Val)
	}

	p.scanToken() // skip ,

	typ.Value, err = p.parseType()
	if err != nil {
		return nil, err
	}

	if p.nextToken.Kind != token.CloseAngle {
		return nil, fmt.Errorf("expected > but got %s", p.nextToken.Val)
	}

	p.scanToken() // skip >

	return typ, nil
}

func (p *Parser) parseArrayType() (typ *ast.TypeArray, err error) {
	typ = &ast.TypeArray{}

	if p.nextToken.Kind != token.OpenBracket {
		return nil, fmt.Errorf("expected [ but got %s", p.nextToken.Val)
	}

	typ.Token = p.nextToken

	p.scanToken() // skip [

	if p.nextToken.Kind != token.CloseBracket {
		return nil, fmt.Errorf("expected ] but got %s", p.nextToken.Val)
	}

	p.scanToken() // skip ]

	typ.Type, err = p.parseType()
	if err != nil {
		return nil, err
	}

	return typ, nil
}

func (p *Parser) parseKeyType() (typ ast.Type, err error) {
	if p.nextToken.Kind != token.Type {
		return nil, fmt.Errorf("expected a type but got %s", p.nextToken.Val)
	}

	switch p.nextToken.Val {
	case "Int8", "Int16", "Int32", "Int64":
		typ = &ast.TypeInt{
			Token: p.nextToken,
			Size:  extractBits("Int", p.nextToken.Val),
		}
	case "Uint8", "Uint16", "Uint32", "Uint64":
		typ = &ast.TypeUint{
			Token: p.nextToken,
			Size:  extractBits("Uint", p.nextToken.Val),
		}
	case "Float32", "Float64":
		typ = &ast.TypeFloat{
			Token: p.nextToken,
			Size:  extractBits("Float", p.nextToken.Val),
		}
	case "String":
		typ = &ast.TypeString{
			Token: p.nextToken,
		}
	case "Bool":
		typ = &ast.TypeBool{
			Token: p.nextToken,
		}
	case "Byte":
		typ = &ast.TypeByte{
			Token: p.nextToken,
		}
	case "Timestamp":
		typ = &ast.TypeTimestamp{
			Token: p.nextToken,
		}
	default:
		return nil, fmt.Errorf("unknown key type %s", p.nextToken.Val)
	}

	return
}

func extractBits(prefix string, value string) int {
	result, _ := strconv.ParseInt(value[len(prefix):], 10, 64)
	return int(result)
}
