package ast

import "strings"

type Program struct {
	Nodes []Node
}

var _ Node = (*Program)(nil)

func (p *Program) nodeLiteral() {}
func (p *Program) String() string {
	var sb strings.Builder

	for i, n := range p.Nodes {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(n.String())
	}

	return sb.String()
}
