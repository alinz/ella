package ast

import (
	"strings"

	"ella.to/internal/token"
)

type Value interface {
	Node
	valueLiteral()
}

// SIGNED INTEGER

type ValueInt struct {
	Token   *token.Token
	Value   int64
	Size    int // 8, 16, 32, 64
	Defined bool
}

var _ Value = (*ValueInt)(nil)

func (v *ValueInt) nodeLiteral()  {}
func (v *ValueInt) valueLiteral() {}
func (v *ValueInt) String() string {
	return v.Token.Val
}

// UNSIGNED INTEGER

type ValueUint struct {
	Token   *token.Token
	Value   uint64
	Size    int // 8, 16, 32, 64
	Defined bool
}

var _ Value = (*ValueUint)(nil)

func (v *ValueUint) nodeLiteral()  {}
func (v *ValueUint) valueLiteral() {}
func (v *ValueUint) String() string {
	return v.Token.Val
}

// FLOAT

type ValueFloat struct {
	Token *token.Token
	Value float64
}

var _ Value = (*ValueFloat)(nil)
var _ Node = (*ValueFloat)(nil)

func (v *ValueFloat) nodeLiteral()  {}
func (v *ValueFloat) valueLiteral() {}
func (v *ValueFloat) String() string {
	return v.Token.Val
}

// STRING

type ValueString struct {
	Token *token.Token
	Value string
}

var _ Value = (*ValueString)(nil)
var _ Node = (*ValueString)(nil)

func (v *ValueString) nodeLiteral()  {}
func (v *ValueString) valueLiteral() {}
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
	Token   *token.Token
	Value   bool
	Defined bool // means if user explicitly set it
}

var _ Value = (*ValueBool)(nil)
var _ Node = (*ValueBool)(nil)

func (v *ValueBool) nodeLiteral()  {}
func (v *ValueBool) valueLiteral() {}
func (v *ValueBool) String() string {
	// NOTE: when a constant defines as a flag, technically it doesn't have a
	// token, so we need to check for nil, if it's nil, the default
	// value is true
	if v.Token == nil {
		return "true"
	}
	return v.Token.Val
}

// NULL

type ValueNull struct {
	Token *token.Token
}

var _ Value = (*ValueNull)(nil)
var _ Node = (*ValueNull)(nil)

func (v *ValueNull) nodeLiteral()  {}
func (v *ValueNull) valueLiteral() {}
func (v *ValueNull) String() string {
	return v.Token.Val
}

// VARIABLE

type ValueVariable struct {
	Token *token.Token
}

var _ Value = (*ValueVariable)(nil)
var _ Node = (*ValueVariable)(nil)

func (v *ValueVariable) nodeLiteral()  {}
func (v *ValueVariable) valueLiteral() {}
func (v *ValueVariable) String() string {
	return v.Token.Val
}
