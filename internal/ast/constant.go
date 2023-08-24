package ast

import (
	"strings"

	"ella.to/internal/token"
)

type Const struct {
	Token *token.Token `json:"token"`
	Name  *Identifier  `json:"name"`
	Value Value        `json:"value"`
}

var _ Statement = (*Const)(nil)

func (c *Const) statementLiteral() {}

func (c *Const) TokenLiteral() string {
	return c.Token.Literal
}

func (c *Const) String() string {
	var sb strings.Builder

	sb.WriteString(c.TokenLiteral())
	sb.WriteString(" ")
	sb.WriteString(c.Name.String())
	sb.WriteString(" = ")
	sb.WriteString(c.Value.String())

	return sb.String()
}
