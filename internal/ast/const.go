package ast

import "strings"

type Const struct {
	Name  *Identifier
	Value Value
}

var _ Node = (*Const)(nil)

func (c *Const) nodeLiteral() {}

func (c *Const) String() string {
	var sb strings.Builder

	sb.WriteString(c.Name.String())

	switch v := c.Value.(type) {
	case *ValueBool:
		if v.Defined {
			sb.WriteString(" = ")
			sb.WriteString(v.String())
		}

	case *ValueInt:
		if v.Defined {
			sb.WriteString(" = ")
			sb.WriteString(v.String())
		}

	default:
		sb.WriteString(" = ")
		sb.WriteString(v.String())
	}

	return sb.String()
}
