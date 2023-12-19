package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"compiler.ella.to/internal/token"
)

type Type interface {
	Node
	typeLiteral()
}

func marshalTextType(t Type) ([]byte, error) {
	var buf bytes.Buffer
	var typ string

	switch t.(type) {
	case *File:
		typ = "file"
	case *CustomType:
		typ = "custom"
	case *Byte:
		typ = "byte"
	case *Uint:
		typ = "uint"
	case *Int:
		typ = "int"
	case *Float:
		typ = "float"
	case *String:
		typ = "string"
	case *Bool:
		typ = "bool"
	case *Any:
		typ = "any"
	case *Array:
		typ = "array"
	case *Map:
		typ = "map"
	case *Timestamp:
		typ = "timestamp"
	default:
		return nil, fmt.Errorf("unknown type for type: %T", t)
	}

	buf.WriteString(`{"type":"`)
	buf.WriteString(typ)
	buf.WriteString(`","data":`)
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	buf.Write(b)
	buf.WriteString(`}`)

	return buf.Bytes(), nil
}

func unmarshalTextType(data []byte, t *Type) error {
	result := struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	switch result.Type {
	case "file":
		*t = &File{}
	case "custom":
		*t = &CustomType{}
	case "byte":
		*t = &Byte{}
	case "uint":
		*t = &Uint{}
	case "int":
		*t = &Int{}
	case "float":
		*t = &Float{}
	case "string":
		*t = &String{}
	case "bool":
		*t = &Bool{}
	case "any":
		*t = &Any{}
	case "array":
		*t = &Array{}
	case "map":
		*t = &Map{}
	case "timestamp":
		*t = &Timestamp{}
	default:
		return fmt.Errorf("unknown type for type: %T", result.Type)
	}

	return json.Unmarshal(result.Data, *t)
}

// FILE TYPE

type File struct {
	Token *token.Token `json:"token"`
}

var _ Type = (*File)(nil)

func (t *File) typeLiteral() {}
func (t *File) TokenLiteral() string {
	return t.Token.Literal
}
func (t *File) String() string {
	return t.Token.Literal
}

// CUSTOM TYPE

type CustomType struct {
	Token *token.Token `json:"token"`
}

var _ Type = (*CustomType)(nil)

func (t *CustomType) typeLiteral() {}
func (t *CustomType) TokenLiteral() string {
	return t.Token.Literal
}
func (t *CustomType) String() string {
	return t.Token.Literal
}

// BYTE

type Byte struct {
	Token *token.Token `json:"token"`
}

var _ Type = (*Byte)(nil)

func (t *Byte) typeLiteral() {}
func (t *Byte) TokenLiteral() string {
	return t.Token.Literal
}
func (t *Byte) String() string {
	return t.Token.Literal
}

// UNSIGNED INTEGER

type Uint struct {
	Token *token.Token `json:"token"`
	Size  int          `json:"size"` // 8, 16, 32, 64
}

var _ Type = (*Uint)(nil)

func (t *Uint) typeLiteral() {}
func (t *Uint) TokenLiteral() string {
	return t.Token.Literal
}
func (t *Uint) String() string {
	return t.Token.Literal
}

// SIGNED INTEGER

type Int struct {
	Token *token.Token `json:"token"`
	Size  int          `json:"size"` // 8, 16, 32, 64
}

var _ Type = (*Int)(nil)

func (t *Int) typeLiteral() {}
func (t *Int) TokenLiteral() string {
	return t.Token.Literal
}
func (t *Int) String() string {
	return t.Token.Literal
}

// FLOAT

type Float struct {
	Token *token.Token `json:"token"`
	Size  int          `json:"size"` // 32, 64
}

var _ Type = (*Float)(nil)

func (t *Float) typeLiteral() {}
func (t *Float) TokenLiteral() string {
	return t.Token.Literal
}
func (t *Float) String() string {
	return t.Token.Literal
}

// STRING

type String struct {
	Token *token.Token `json:"token"`
}

var _ Type = (*String)(nil)

func (t *String) typeLiteral() {}
func (t *String) TokenLiteral() string {
	return t.Token.Literal
}
func (t *String) String() string {
	return t.Token.Literal
}

// BOOL

type Bool struct {
	Token *token.Token `json:"token"`
}

var _ Type = (*Bool)(nil)

func (t *Bool) typeLiteral() {}
func (t *Bool) TokenLiteral() string {
	return t.Token.Literal
}
func (t *Bool) String() string {
	return t.Token.Literal
}

// ANY

type Any struct {
	Token *token.Token `json:"token"`
}

var _ Type = (*Any)(nil)

func (t *Any) typeLiteral() {}
func (t *Any) TokenLiteral() string {
	return t.Token.Literal
}
func (t *Any) String() string {
	return t.Token.Literal
}

// ARRAY

type Array struct {
	Token *token.Token `json:"token"` // this is the '[' token
	Type  Type         `json:"type"`
}

var _ Type = (*Array)(nil)

func (t *Array) typeLiteral() {}
func (t *Array) TokenLiteral() string {
	return t.Token.Literal
}
func (t *Array) String() string {
	var sb strings.Builder

	sb.WriteString("[]")
	sb.WriteString(t.Type.String())

	return sb.String()
}

// MAP

type Map struct {
	Token *token.Token `json:"token"`
	Key   Type         `json:"key"`
	Value Type         `json:"value"`
}

var _ Type = (*Map)(nil)

func (t *Map) typeLiteral() {}
func (t *Map) TokenLiteral() string {
	return t.Token.Literal
}
func (t *Map) String() string {
	var sb strings.Builder

	sb.WriteString("map<")
	sb.WriteString(t.Key.String())
	sb.WriteString(", ")
	sb.WriteString(t.Value.String())
	sb.WriteString(">")

	return sb.String()
}

// TIMESTAMP

type Timestamp struct {
	Token *token.Token `json:"token"`
}

var _ Type = (*Timestamp)(nil)

func (t *Timestamp) typeLiteral() {}
func (t *Timestamp) TokenLiteral() string {
	return t.Token.Literal
}
func (t *Timestamp) String() string {
	return t.Token.Literal
}
