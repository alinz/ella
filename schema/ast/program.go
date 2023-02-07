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

	for i, n := range p.Nodes {
		sb.WriteString(n.TokenLiteral())
		if i < len(p.Nodes)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
