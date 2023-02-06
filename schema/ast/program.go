package ast

import (
	"strings"
)

type Program struct {
	Nodes []Node
}

var _ Node = (*Program)(nil)

func (p *Program) TokenLiteral() string {
	var sb strings.Builder

	for _, n := range p.Nodes {
		sb.WriteString(n.TokenLiteral())
		sb.WriteString("\n")
	}

	return sb.String()
}
