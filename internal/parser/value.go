package parser

import (
	"math"
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
			Size:  getFloatSize(float),
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
			Size:    getIntSize(integer, integer),
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

// find out about the min size for integer based on min and max values
// 8, –128, 127
// 16, –32768, 32767
// 32, -2147483648, 2147483647
// 64, -9223372036854775808, 9223372036854775807
func getIntSize(min, max int64) int {
	if min >= -128 && max <= 127 {
		return 8
	} else if min >= -32768 && max <= 32767 {
		return 16
	} else if min >= -2147483648 && max <= 2147483647 {
		return 32
	} else {
		return 64
	}
}

func getFloatSize(value float64) int {
	if value >= math.SmallestNonzeroFloat32 && value <= math.MaxFloat32 {
		return 32
	}
	return 64
}
