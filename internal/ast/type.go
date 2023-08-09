package ast

import (
	"strings"

	"ella.to/internal/token"
)

type Type interface {
	Node
	typeLiteral()
}

// FILE TYPE

type File struct {
	Token *token.Token
}

var _ Type = (*File)(nil)

func (t *File) nodeLiteral() {}
func (t *File) typeLiteral() {}
func (t *File) String() string {
	return t.Token.Val
}

// CUSTOM TYPE

type CustomType struct {
	Token *token.Token
}

var _ Type = (*CustomType)(nil)

func (t *CustomType) nodeLiteral() {}
func (t *CustomType) typeLiteral() {}
func (t *CustomType) String() string {
	return t.Token.Val
}

// BYTE

type Byte struct {
	Token *token.Token
}

var _ Type = (*Byte)(nil)

func (t *Byte) nodeLiteral() {}
func (t *Byte) typeLiteral() {}
func (t *Byte) String() string {
	return t.Token.Val
}

// UNSIGNED INTEGER

type Uint struct {
	Token *token.Token
	Size  int // 8, 16, 32, 64
}

var _ Type = (*Uint)(nil)

func (t *Uint) nodeLiteral() {}
func (t *Uint) typeLiteral() {}
func (t *Uint) String() string {
	return t.Token.Val
}

// SIGNED INTEGER

type Int struct {
	Token *token.Token
	Size  int // 8, 16, 32, 64
}

var _ Type = (*Int)(nil)

func (t *Int) nodeLiteral() {}
func (t *Int) typeLiteral() {}
func (t *Int) String() string {
	return t.Token.Val
}

// FLOAT

type Float struct {
	Token *token.Token
	Size  int // 32, 64
}

var _ Type = (*Float)(nil)

func (t *Float) nodeLiteral() {}
func (t *Float) typeLiteral() {}
func (t *Float) String() string {
	return t.Token.Val
}

// STRING

type String struct {
	Token *token.Token
}

var _ Type = (*String)(nil)

func (t *String) nodeLiteral() {}
func (t *String) typeLiteral() {}
func (t *String) String() string {
	return t.Token.Val
}

// BOOL

type Bool struct {
	Token *token.Token
}

var _ Type = (*Bool)(nil)

func (t *Bool) nodeLiteral() {}
func (t *Bool) typeLiteral() {}
func (t *Bool) String() string {
	return t.Token.Val
}

// ANY

type Any struct {
	Token *token.Token
}

var _ Type = (*Any)(nil)

func (t *Any) nodeLiteral() {}
func (t *Any) typeLiteral() {}
func (t *Any) String() string {
	return t.Token.Val
}

// ARRAY

type Array struct {
	Token *token.Token // this is the '[' token
	Type  Type
}

var _ Type = (*Array)(nil)

func (t *Array) nodeLiteral() {}
func (t *Array) typeLiteral() {}
func (t *Array) String() string {
	var sb strings.Builder

	sb.WriteString("[]")
	sb.WriteString(t.Type.String())

	return sb.String()
}

// MAP

type Map struct {
	Token *token.Token
	Key   Type
	Value Type
}

var _ Type = (*Map)(nil)

func (t *Map) nodeLiteral() {}
func (t *Map) typeLiteral() {}
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
	Token *token.Token
}

var _ Type = (*Timestamp)(nil)

func (t *Timestamp) nodeLiteral() {}
func (t *Timestamp) typeLiteral() {}
func (t *Timestamp) String() string {
	return t.Token.Val
}
