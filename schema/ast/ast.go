package ast

import (
	"github.com/alinz/ella.to/schema/token"
)

type Node interface {
	TokenLiteral() string
}

type Identifier struct {
	Name  string
	Token *token.Token
}

var _ Node = (*Identifier)(nil)

func (i *Identifier) TokenLiteral() string {
	return i.Token.Val
}
