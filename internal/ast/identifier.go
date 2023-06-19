package ast

import "ella.to/internal/token"

type Identifier struct {
	Token *token.Token
}

var _ Node = (*Identifier)(nil)

func (i *Identifier) nodeLiteral() {}

func (i *Identifier) String() string {
	return i.Token.Val
}
