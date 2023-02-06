package ast

import (
	"strings"
)

type Constant struct {
	Name  *Identifier
	Value Value
}

func (c *Constant) TokenLiteral() string {
	var sb strings.Builder

	sb.WriteString(c.Name.TokenLiteral())
	if c.Value != nil {
		sb.WriteString(" = ")
		sb.WriteString(c.Value.TokenLiteral())
	}

	return sb.String()
}
