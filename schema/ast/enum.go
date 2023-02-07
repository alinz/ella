package ast

import (
	"strings"
)

type Enum struct {
	Name      *Identifier
	Type      Type
	Constants []*Constant
}

func (e *Enum) TokenLiteral() string {
	var sb strings.Builder

	sb.WriteString("enum ")
	sb.WriteString(e.Name.TokenLiteral())
	sb.WriteString(" ")
	sb.WriteString(e.Type.TokenLiteral())
	sb.WriteString(" {")
	for _, c := range e.Constants {
		sb.WriteString("\n\t")
		sb.WriteString(c.TokenLiteral())
	}

	if len(e.Constants) > 0 {
		sb.WriteString("\n")
	}

	sb.WriteString("}")

	return sb.String()
}
