package ast

import (
	"strings"

	"ella.to/internal/token"
)

type Enum struct {
	Token     *token.Token // enum token
	Name      *Identifier
	Type      Type
	Constants []*Const
}

var _ Node = (*Enum)(nil)

func (e *Enum) nodeLiteral() {}
func (e *Enum) String() string {
	var sb strings.Builder

	sb.WriteString("enum ")
	sb.WriteString(e.Name.String())
	sb.WriteString(" ")
	sb.WriteString(e.Type.String())
	sb.WriteString(" {")

	for _, c := range e.Constants {
		sb.WriteString("\n\t")
		sb.WriteString(c.String())
	}

	if len(e.Constants) > 0 {
		sb.WriteString("\n")
	}

	sb.WriteString("}")

	return sb.String()
}
