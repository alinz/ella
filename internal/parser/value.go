package parser

import (
	"strconv"

	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func (p *Parser) parseValue() (ast.Value, error) {
	switch p.currTok.Type {
	case token.ConstFloat:
		value, err := strconv.ParseFloat(p.currTok.Val, 64)
		if err != nil {
			return nil, p.newError(p.currTok, "failed to parse float value", err)
		}
		return &ast.ValueFloat{
			Token: p.currTok,
			Value: value,
		}, nil
	case token.ConstInt:
		value, err := strconv.ParseInt(p.currTok.Val, 10, 64)
		if err != nil {
			return nil, p.newError(p.currTok, "failed to parse int value", err)
		}
		return &ast.ValueInt{
			Token: p.currTok,
			Value: value,
		}, nil
	case token.ConstBool:
		value, err := strconv.ParseBool(p.currTok.Val)
		if err != nil {
			return nil, p.newError(p.currTok, "failed to parse bool value", err)
		}
		return &ast.ValueBool{
			Token:     p.currTok,
			Value:     value,
			IsUserSet: true,
		}, nil
	case token.ConstNull:
		return &ast.ValueNull{
			Token: p.currTok,
		}, nil
	case token.ConstStringSingleQuote, token.ConstStringDoubleQuote, token.ConstStringBacktickQoute:
		return &ast.ValueString{
			Token: p.currTok,
			Value: p.currTok.Val,
		}, nil
	}
	return nil, p.newError(p.currTok, "expected one of the following, 'int', 'float', 'bool', 'null', 'string' values")
}
