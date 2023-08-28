package ast

import (
	"strings"

	"ella.to/internal/token"
)

type Base struct {
	Token   *token.Token `json:"token"`
	Name    *Identifier  `json:"name"`
	Type    Type         `json:"type"`
	Options Options      `json:"options"`
}

var _ Statement = (*Base)(nil)

func (b *Base) statementLiteral() {}

func (b *Base) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Base) String() string {
	var sb strings.Builder

	sb.WriteString(b.TokenLiteral())
	sb.WriteString(" ")
	sb.WriteString(b.Name.String())
	sb.WriteString(" ")
	sb.WriteString(b.Type.String())
	sb.WriteString(b.Options.String(1))

	return sb.String()
}
