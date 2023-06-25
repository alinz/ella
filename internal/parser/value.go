package parser

import (
	"strconv"

	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func ParseValue(p *Parser) (value ast.Value, err error) {
	peekTok := p.Peek()

	switch peekTok.Type {
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
	default:
		return nil, p.WithError(peekTok, "expected one of the following, 'int', 'float', 'bool', 'null', 'string' values")
	}

	p.Next() // skip value if no error

	return value, nil
}
