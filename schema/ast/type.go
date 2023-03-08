package ast

import (
	"strings"

	"ella.to/schema/token"
)

type Type interface {
	Node
	kindNode()
}

type TypeCustom struct {
	Token *token.Token
	Name  string
}

var _ Type = (*TypeCustom)(nil)

func (v TypeCustom) kindNode() {}
func (v TypeCustom) TokenLiteral() string {
	return v.Token.Val
}

type TypeAny struct {
	Token *token.Token
}

var _ Type = (*TypeAny)(nil)

func (v TypeAny) kindNode() {}
func (v TypeAny) TokenLiteral() string {
	return v.Token.Val
}

type TypeInt struct {
	Token *token.Token
	Size  int // 8 | 16 | 32 | 64
}

var _ Type = (*TypeInt)(nil)

func (v TypeInt) kindNode() {}
func (v TypeInt) TokenLiteral() string {
	return v.Token.Val
}

type TypeUint struct {
	Token *token.Token
	Size  int // 8 | 16 | 32 | 64
}

var _ Type = (*TypeUint)(nil)

func (v TypeUint) kindNode() {}
func (v TypeUint) TokenLiteral() string {
	return v.Token.Val
}

type TypeByte struct {
	Token *token.Token
}

var _ Type = (*TypeByte)(nil)

func (v TypeByte) kindNode() {}
func (v TypeByte) TokenLiteral() string {
	return v.Token.Val
}

type TypeFloat struct {
	Token *token.Token
	Size  int // 32 | 64
}

var _ Type = (*TypeFloat)(nil)

func (v TypeFloat) kindNode() {}
func (v TypeFloat) TokenLiteral() string {
	return v.Token.Val
}

type TypeString struct {
	Token *token.Token
}

var _ Type = (*TypeString)(nil)

func (v TypeString) kindNode() {}
func (v TypeString) TokenLiteral() string {
	return v.Token.Val
}

type TypeBool struct {
	Token *token.Token
}

var _ Type = (*TypeBool)(nil)

func (v TypeBool) kindNode() {}
func (v TypeBool) TokenLiteral() string {
	return v.Token.Val
}

type TypeTimestamp struct {
	Token *token.Token
}

func (v TypeTimestamp) kindNode() {}
func (v TypeTimestamp) TokenLiteral() string {
	return v.Token.Val
}

type TypeMap struct {
	Token *token.Token
	Key   Type
	Value Type
}

var _ Type = (*TypeMap)(nil)

func (v TypeMap) kindNode() {}
func (v TypeMap) TokenLiteral() string {
	var sb strings.Builder

	sb.WriteString(v.Token.Val)
	sb.WriteString("<")
	sb.WriteString(v.Key.TokenLiteral())
	sb.WriteString(", ")
	sb.WriteString(v.Value.TokenLiteral())
	sb.WriteString(">")

	return sb.String()
}

type TypeArray struct {
	Token *token.Token
	Type  Type
}

var _ Type = (*TypeArray)(nil)

func (v TypeArray) kindNode() {}
func (v TypeArray) TokenLiteral() string {
	return v.Token.Val
}
