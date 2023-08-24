package ast

import (
	"strings"

	"ella.to/internal/token"
)

type Define struct {
	Token   *token.Token `json:"token"`
	Name    *Identifier  `json:"name"`
	Type    Type         `json:"type"`
	Options Options      `json:"options"`
}

var _ Statement = (*Define)(nil)

func (d *Define) statementLiteral() {}

func (d *Define) TokenLiteral() string {
	return d.Token.Literal
}

func (d *Define) String() string {
	var sb strings.Builder

	sb.WriteString(d.TokenLiteral())
	sb.WriteString(" ")
	sb.WriteString(d.Name.String())
	sb.WriteString(" ")
	sb.WriteString(d.Type.String())
	sb.WriteString(d.Options.String(1))

	return sb.String()
}
