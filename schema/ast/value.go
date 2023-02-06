package ast

import (
	"github.com/alinz/rpc.go/schema/token"
)

type Value interface {
	valueNode()
	Node
}

type ValueString struct {
	Token   *token.Token
	Content string
}

func (v ValueString) valueNode() {}
func (v *ValueString) TokenLiteral() string {
	return v.Token.Val
}

type ValueInt struct {
	Token   *token.Token
	Content int64
}

func (v ValueInt) valueNode() {}
func (v *ValueInt) TokenLiteral() string {
	return v.Token.Val
}

type ValueFloat struct {
	Token   *token.Token
	Content float64
}

func (v ValueFloat) valueNode() {}
func (v *ValueFloat) TokenLiteral() string {
	return v.Token.Val
}

type ValueBool struct {
	Token   *token.Token
	Content bool
}

func (v ValueBool) valueNode() {}
func (v *ValueBool) TokenLiteral() string {
	return v.Token.Val
}
