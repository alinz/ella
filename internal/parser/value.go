package parser

import (
	"strconv"

	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func parseBytesSizeNumber(value string) (number string, mul int64) {
	switch value[len(value)-2] {
	case 'K':
		mul = 1024
	case 'M':
		mul = 1024 * 1024
	case 'G':
		mul = 1024 * 1024 * 1024
	case 'T':
		mul = 1024 * 1024 * 1024 * 1024
	case 'P':
		mul = 1024 * 1024 * 1024 * 1024 * 1024
	case 'E':
		mul = 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	default:
		return value[:len(value)-1], 1
	}

	return value[:len(value)-2], mul
}

func ParseValue(p *Parser) (value ast.Value, err error) {
	peekTok := p.Peek()

	switch peekTok.Type {
	case token.ConstFloatBytes:
		num, mul := parseBytesSizeNumber(peekTok.Val)
		float, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return nil, p.WithError(peekTok, "failed to parse float value", err)
		}
		// NOTE: because the size of bytes usually normalized to int,
		// we treat float bytes as int bytes
		// also first we need to multiply float to mul then convert it to int
		// to make sure int64(0.5) will become 0 and make the whole multiplication
		// to zero
		value = &ast.ValueInt{
			Token:   peekTok,
			Value:   int64(float * float64(mul)),
			Defined: true,
		}
	case token.ConstIntBytes:
		num, mul := parseBytesSizeNumber(peekTok.Val)
		integer, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			return nil, p.WithError(peekTok, "failed to parse int value", err)
		}
		value = &ast.ValueInt{
			Token:   peekTok,
			Value:   integer * mul,
			Defined: true,
		}
	case token.ConstFloat:
		float, err := strconv.ParseFloat(peekTok.Val, 64)
		if err != nil {
			return nil, p.WithError(peekTok, "failed to parse float value", err)
		}
		value = &ast.ValueFloat{
			Token: peekTok,
			Value: float,
		}
	case token.ConstInt:
		integer, err := strconv.ParseInt(peekTok.Val, 10, 64)
		if err != nil {
			return nil, p.WithError(peekTok, "failed to parse int value", err)
		}
		value = &ast.ValueInt{
			Token:   peekTok,
			Value:   integer,
			Defined: true,
		}
	case token.ConstBool:
		boolean, err := strconv.ParseBool(peekTok.Val)
		if err != nil {
			return nil, p.WithError(peekTok, "failed to parse bool value", err)
		}
		value = &ast.ValueBool{
			Token:   peekTok,
			Value:   boolean,
			Defined: true,
		}
	case token.ConstNull:
		value = &ast.ValueNull{
			Token: peekTok,
		}
	case token.ConstStringSingleQuote, token.ConstStringDoubleQuote, token.ConstStringBacktickQoute:
		value = &ast.ValueString{
			Token: peekTok,
			Value: peekTok.Val,
		}
	case token.Identifier:
		value = &ast.ValueVariable{
			Token: peekTok,
		}
	default:
		return nil, p.WithError(peekTok, "expected one of the following, 'int', 'float', 'bool', 'null', 'string' values")
	}

	p.Next() // skip value if no error

	return value, nil
}
