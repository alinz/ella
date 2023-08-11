package parser

import (
	"strconv"
	"time"

	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func parseBytesNumber(value string) (number string, mul int64) {
	switch value[len(value)-2] {
	case 'k':
		mul = 1024
	case 'm':
		mul = 1024 * 1024
	case 'g':
		mul = 1024 * 1024 * 1024
	case 't':
		mul = 1024 * 1024 * 1024 * 1024
	case 'p':
		mul = 1024 * 1024 * 1024 * 1024 * 1024
	case 'e':
		mul = 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	default:
		return value[:len(value)-1], 1
	}

	return value[:len(value)-2], mul
}

func parseDurationNumber(value string) (number string, mul int64) {
	switch value[len(value)-2] {
	case 'n':
		mul = int64(time.Nanosecond)
		return value[:len(value)-2], mul
	case 'u':
		mul = int64(time.Microsecond)
		return value[:len(value)-2], mul
	case 'm':
		mul = int64(time.Millisecond)
		return value[:len(value)-2], mul
	default:
		switch value[len(value)-1] {
		case 's':
			mul = int64(time.Second)
		case 'm':
			mul = int64(time.Minute)
		case 'h':
			mul = int64(time.Hour)
		}
		return value[:len(value)-1], mul
	}
}

func ParseValue(p *Parser) (value ast.Value, err error) {
	peekTok := p.Peek()

	switch peekTok.Type {
	case token.ConstBytes:
		num, mul := parseBytesNumber(peekTok.Val)
		integer, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			return nil, p.WithError(peekTok, "failed to parse int value for bytes size", err)
		}
		value = &ast.ValueInt{
			Token:   peekTok,
			Value:   integer * mul,
			Defined: true,
		}
	case token.ConstDuration:
		num, mul := parseDurationNumber(peekTok.Val)
		integer, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			return nil, p.WithError(peekTok, "failed to parse int value for duration size", err)
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
