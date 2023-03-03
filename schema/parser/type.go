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
	case token.Map:
		return p.parseMapType()
	case token.OpenSquare:
		return p.parseArrayType()
	case token.Int8, token.Int16, token.Int32, token.Int64:
		typ = &ast.TypeInt{
			Token: p.nextToken,
			Size:  extractBits("int", p.nextToken.Val),
		}
	case token.Uint8, token.Uint16, token.Uint32, token.Uint64:
		typ = &ast.TypeUint{
			Token: p.nextToken,
			Size:  extractBits("uint", p.nextToken.Val),
		}
	case token.Float32, token.Float64:
		typ = &ast.TypeFloat{
			Token: p.nextToken,
			Size:  extractBits("float", p.nextToken.Val),
		}
	case token.String:
		typ = &ast.TypeString{
			Token: p.nextToken,
		}
	case token.Bool:
		typ = &ast.TypeBool{
			Token: p.nextToken,
		}
	case token.Byte:
		typ = &ast.TypeByte{
			Token: p.nextToken,
		}
	case token.Timestamp:
		typ = &ast.TypeTimestamp{
			Token: p.nextToken,
		}
	case token.Any:
		typ = &ast.TypeAny{
			Token: p.nextToken,
		}
	case token.Word:
		if !isIdentifier(p.nextToken.Val, true) {
			return nil, fmt.Errorf("expected a valid name for custom type but got %s", p.nextToken.Val)
		}
		typ = &ast.TypeCustom{
			Token: p.nextToken,
			Name:  p.nextToken.Val,
		}
	default:
		return nil, fmt.Errorf("expected type but got %s", p.nextToken.Val)
	}

	p.scanToken() // skip type

	return typ, nil
}

func (p *Parser) parseMapType() (typ *ast.TypeMap, err error) {
	typ = &ast.TypeMap{}

	if p.nextToken.Kind != token.Map || p.nextToken.Val != "map" {
		return nil, fmt.Errorf("expected map keyword but got %s", p.nextToken.Val)
	}

	typ.Token = p.nextToken

	p.scanToken() // skip map

	if p.nextToken.Kind != token.OpenAngle {
		return nil, fmt.Errorf("expected < after map keyword but got %s", p.nextToken.Val)
	}

	p.scanToken() // skip <

	typ.Key, err = p.parseKeyType()
	if err != nil {
		return nil, err
	}

	if p.nextToken.Kind != token.Comma {
		return nil, fmt.Errorf("expected , after defining key type for map but got %s", p.nextToken.Val)
	}

	p.scanToken() // skip ,

	typ.Value, err = p.parseType()
	if err != nil {
		return nil, err
	}

	if p.nextToken.Kind != token.CloseAngle {
		return nil, fmt.Errorf("expected > after defining key and value types for map but got %s", p.nextToken.Val)
	}

	p.scanToken() // skip >

	return typ, nil
}

func (p *Parser) parseArrayType() (typ *ast.TypeArray, err error) {
	typ = &ast.TypeArray{}

	if p.nextToken.Kind != token.OpenSquare {
		return nil, fmt.Errorf("expected [ for defining array but got %s", p.nextToken.Val)
	}

	typ.Token = p.nextToken

	p.scanToken() // skip [

	if p.nextToken.Kind != token.CloseSquare {
		return nil, fmt.Errorf("expected ] for defining array but got %s", p.nextToken.Val)
	}

	p.scanToken() // skip ]

	typ.Type, err = p.parseType()
	if err != nil {
		return nil, err
	}

	return typ, nil
}

func (p *Parser) parseKeyType() (typ ast.Type, err error) {
	p.nextToken.OneOf(token.Uint64)

	switch p.nextToken.Val {
	case "int8", "int16", "int32", "int64":
		typ = &ast.TypeInt{
			Token: p.nextToken,
			Size:  extractBits("int", p.nextToken.Val),
		}
	case "uint8", "uint16", "uint32", "uint64":
		typ = &ast.TypeUint{
			Token: p.nextToken,
			Size:  extractBits("uint", p.nextToken.Val),
		}
	case "float32", "float64":
		typ = &ast.TypeFloat{
			Token: p.nextToken,
			Size:  extractBits("float", p.nextToken.Val),
		}
	case "string":
		typ = &ast.TypeString{
			Token: p.nextToken,
		}
	case "bool":
		typ = &ast.TypeBool{
			Token: p.nextToken,
		}
	case "byte":
		typ = &ast.TypeByte{
			Token: p.nextToken,
		}
	case "timestamp":
		typ = &ast.TypeTimestamp{
			Token: p.nextToken,
		}
	default:
		return nil, fmt.Errorf("unknown key type %s", p.nextToken.Val)
	}

	p.scanToken()

	return
}

func extractBits(prefix string, value string) int {
	result, _ := strconv.ParseInt(value[len(prefix):], 10, 64)
	return int(result)
}
