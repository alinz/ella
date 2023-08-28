package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"ella.to/internal/token"
)

type Value interface {
	Node
	valueLiteral()
}

func marshalTextValue(value Value) ([]byte, error) {
	var buf bytes.Buffer
	var typ string

	switch value.(type) {
	case *ValueByteSize:
		typ = "byteSize"
	case *ValueDuration:
		typ = "duration"
	case *ValueInt:
		typ = "int"
	case *ValueUint:
		typ = "uint"
	case *ValueFloat:
		typ = "float"
	case *ValueString:
		typ = "string"
	case *ValueBool:
		typ = "bool"
	case *ValueNull:
		typ = "null"
	case *ValueVariable:
		typ = "variable"
	default:
		return nil, fmt.Errorf("unknown type for value: %T", value)
	}

	buf.WriteString(`{"type":"`)
	buf.WriteString(typ)
	buf.WriteString(`","data":`)
	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	buf.Write(b)
	buf.WriteString(`}`)

	return buf.Bytes(), nil
}

func unmarshalTextValue(data []byte, value *Value) error {
	result := struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	switch result.Type {
	case "byteSize":
		*value = &ValueByteSize{}
	case "duration":
		*value = &ValueDuration{}
	case "int":
		*value = &ValueInt{}
	case "uint":
		*value = &ValueUint{}
	case "float":
		*value = &ValueFloat{}
	case "string":
		*value = &ValueString{}
	case "bool":
		*value = &ValueBool{}
	case "null":
		*value = &ValueNull{}
	case "variable":
		*value = &ValueVariable{}
	default:
		return fmt.Errorf("unknown type for value: %s", result.Type)
	}

	return json.Unmarshal(result.Data, *value)
}

// BYTE SIZE

type ByteSize int64

const (
	ByteSizeB  ByteSize = 1
	ByteSizeKB          = ByteSizeB * 1024
	ByteSizeMB          = ByteSizeKB * 1024
	ByteSizeGB          = ByteSizeMB * 1024
	ByteSizeTB          = ByteSizeGB * 1024
	ByteSizePB          = ByteSizeTB * 1024
	ByteSizeEB          = ByteSizePB * 1024
)

func (b ByteSize) String() string {
	switch b {
	case ByteSizeB:
		return "b"
	case ByteSizeKB:
		return "kb"
	case ByteSizeMB:
		return "mb"
	case ByteSizeGB:
		return "gb"
	case ByteSizeTB:
		return "tb"
	case ByteSizePB:
		return "pb"
	case ByteSizeEB:
		return "eb"
	default:
		panic("unknown byte size")
	}
}

type ValueByteSize struct {
	Token *token.Token `json:"token"`
	Value int64        `json:"value"`
	Scale ByteSize     `json:"scale"`
}

var _ Value = (*ValueByteSize)(nil)

func (v *ValueByteSize) valueLiteral() {}
func (v *ValueByteSize) TokenLiteral() string {
	return v.Token.Literal
}
func (v *ValueByteSize) String() string {
	return v.Token.Literal
}

// DURATION

type DurationScale int64

const (
	DurationScaleNanosecond  DurationScale = 1
	DurationScaleMicrosecond               = DurationScaleNanosecond * 1000
	DurationScaleMillisecond               = DurationScaleMicrosecond * 1000
	DurationScaleSecond                    = DurationScaleMillisecond * 1000
	DurationScaleMinute                    = DurationScaleSecond * 60
	DurationScaleHour                      = DurationScaleMinute * 60
)

func (d DurationScale) String() string {
	switch d {
	case DurationScaleNanosecond:
		return "ns"
	case DurationScaleMicrosecond:
		return "us"
	case DurationScaleMillisecond:
		return "ms"
	case DurationScaleSecond:
		return "s"
	case DurationScaleMinute:
		return "m"
	case DurationScaleHour:
		return "h"
	default:
		panic("unknown duration scale")
	}
}

type ValueDuration struct {
	Token *token.Token  `json:"token"`
	Value int64         `json:"value"`
	Scale DurationScale `json:"scale"`
}

var _ Value = (*ValueDuration)(nil)

func (v *ValueDuration) valueLiteral() {}
func (v *ValueDuration) TokenLiteral() string {
	return v.Token.Literal
}
func (v *ValueDuration) String() string {
	return v.Token.Literal
}

// SIGNED INTEGER

type ValueInt struct {
	Token   *token.Token `json:"token"`
	Value   int64        `json:"value"`
	Size    int          `json:"size"`    // 8, 16, 32, 64
	Defined bool         `json:"defined"` // means if user explicitly set it
}

var _ Value = (*ValueInt)(nil)

func (v *ValueInt) valueLiteral() {}
func (v *ValueInt) TokenLiteral() string {
	return v.Token.Literal
}
func (v *ValueInt) String() string {
	return v.Token.Literal
}

// UNSIGNED INTEGER

type ValueUint struct {
	Token *token.Token `json:"token"`
	Value uint64       `json:"value"`
	Size  int          `json:"size"` // 8, 16, 32, 64
}

var _ Value = (*ValueUint)(nil)

func (v *ValueUint) valueLiteral() {}
func (v *ValueUint) TokenLiteral() string {
	return v.Token.Literal
}
func (v *ValueUint) String() string {
	return v.Token.Literal
}

// FLOAT

type ValueFloat struct {
	Token *token.Token `json:"token"`
	Value float64      `json:"value"`
	Size  int          `json:"size"` // 32, 64
}

var _ Value = (*ValueFloat)(nil)

func (v *ValueFloat) valueLiteral() {}
func (v *ValueFloat) TokenLiteral() string {
	return v.Token.Literal
}
func (v *ValueFloat) String() string {
	return v.Token.Literal
}

// STRING

type ValueString struct {
	Token *token.Token `json:"token"`
	Value string       `json:"value"`
}

var _ Value = (*ValueString)(nil)

func (v *ValueString) valueLiteral() {}
func (v *ValueString) TokenLiteral() string {
	return v.Token.Literal
}
func (v *ValueString) String() string {
	var sb strings.Builder

	switch v.Token.Type {
	case token.ConstStringSingleQuote:
		sb.WriteString("'")
		sb.WriteString(v.Value)
		sb.WriteString("'")
	case token.ConstStringDoubleQuote:
		sb.WriteString("\"")
		sb.WriteString(v.Value)
		sb.WriteString("\"")
	case token.ConstStringBacktickQoute:
		sb.WriteString("`")
		sb.WriteString(v.Value)
		sb.WriteString("`")
	}

	return sb.String()
}

// BOOL

type ValueBool struct {
	Token   *token.Token `json:"token"`
	Value   bool         `json:"value"`
	Defined bool         `json:"defined"` // means if user explicitly set it to the values
	/////////////// only used for fmt tools
}

var _ Value = (*ValueBool)(nil)

func (v *ValueBool) valueLiteral() {}
func (v *ValueBool) TokenLiteral() string {
	return v.Token.Literal
}
func (v *ValueBool) String() string {
	// NOTE: when a constant defines as a flag, technically it doesn't have a
	// token, that's why we need to return empty
	if v.Token == nil {
		return ""
	}

	return v.Token.Literal
}

// NULL

type ValueNull struct {
	Token *token.Token `json:"token"`
}

var _ Value = (*ValueNull)(nil)

func (v *ValueNull) valueLiteral() {}
func (v *ValueNull) TokenLiteral() string {
	return v.Token.Literal
}
func (v *ValueNull) String() string {
	return v.Token.Literal
}

// VARIABLE

type ValueVariable struct {
	Token *token.Token `json:"token"`
}

var _ Value = (*ValueVariable)(nil)

func (v *ValueVariable) valueLiteral() {}
func (v *ValueVariable) TokenLiteral() string {
	return v.Token.Literal
}
func (v *ValueVariable) String() string {
	return v.Token.Literal
}
