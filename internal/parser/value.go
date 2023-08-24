package parser

import (
	"strconv"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/token"
)

func parseBytesNumber(value string) (number string, scale ast.ByteSize) {
	switch value[len(value)-2] {
	case 'k':
		scale = ast.ByteSizeKB
	case 'm':
		scale = ast.ByteSizeMB
	case 'g':
		scale = ast.ByteSizeGB
	case 't':
		scale = ast.ByteSizeTB
	case 'p':
		scale = ast.ByteSizePB
	case 'e':
		scale = ast.ByteSizeEB
	default:
		return value[:len(value)-1], 1
	}

	return value[:len(value)-2], scale
}

func parseDurationNumber(value string) (number string, scale ast.DurationScale) {
	switch value[len(value)-2] {
	case 'n':
		scale = ast.DurationScaleNanosecond
		return value[:len(value)-2], scale
	case 'u':
		scale = ast.DurationScaleMicrosecond
		return value[:len(value)-2], scale
	case 'm':
		scale = ast.DurationScaleMillisecond
		return value[:len(value)-2], scale
	default:
		switch value[len(value)-1] {
		case 's':
			scale = ast.DurationScaleSecond
		case 'm':
			scale = ast.DurationScaleMinute
		case 'h':
			scale = ast.DurationScaleHour
		}
		return value[:len(value)-1], scale
	}
}

func ParseValue(p *Parser) (value ast.Value, err error) {
	peekTok := p.Peek()

	switch peekTok.Type {
	case token.ConstBytes:
		literal := strings.ReplaceAll(peekTok.Literal, "_", "")
		num, scale := parseBytesNumber(literal)
		integer, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			return nil, p.WithError(peekTok, "failed to parse int value for bytes size", err)
		}
		value = &ast.ValueByteSize{
			Token: peekTok,
			Value: integer,
			Scale: scale,
		}
	case token.ConstDuration:
		literal := strings.ReplaceAll(peekTok.Literal, "_", "")
		num, scale := parseDurationNumber(literal)
		integer, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			return nil, p.WithError(peekTok, "failed to parse int value for duration size", err)
		}
		value = &ast.ValueDuration{
			Token: peekTok,
			Value: integer,
			Scale: scale,
		}
	case token.ConstFloat:
		literal := strings.ReplaceAll(peekTok.Literal, "_", "")
		float, err := strconv.ParseFloat(literal, 64)
		if err != nil {
			return nil, p.WithError(peekTok, "failed to parse float value", err)
		}
		value = &ast.ValueFloat{
			Token: peekTok,
			Value: float,
		}
	case token.ConstInt:
		literal := strings.ReplaceAll(peekTok.Literal, "_", "")
		integer, err := strconv.ParseInt(literal, 10, 64)
		if err != nil {
			return nil, p.WithError(peekTok, "failed to parse int value", err)
		}
		value = &ast.ValueInt{
			Token:   peekTok,
			Value:   integer,
			Defined: true,
		}
	case token.ConstBool:
		boolean, err := strconv.ParseBool(peekTok.Literal)
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
			Value: peekTok.Literal,
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
