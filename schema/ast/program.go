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

	previous := "constant"

	for i, n := range p.Nodes {
		switch n.(type) {
		case *Constant:
			if previous != "constant" {
				sb.WriteString("\n")
			}
			previous = "constant"
		case *Enum:
			if previous != "enum" {
				sb.WriteString("\n")
			}
			previous = "enum"
		case *Message:
			if previous != "message" {
				sb.WriteString("\n")
			}
			previous = "message"
		case *Service:
		}

		sb.WriteString(n.TokenLiteral())
		if i < len(p.Nodes)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
