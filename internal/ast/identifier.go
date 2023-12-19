package ast

import "compiler.ella.to/internal/token"

type Identifier struct {
	Token *token.Token `json:"token"`
}

var _ Node = (*Identifier)(nil)

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Token.Literal
}
