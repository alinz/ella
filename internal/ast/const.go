package ast

import "strings"

type Const struct {
	Name    *Identifier
	Value   Value
	Comment *Comment
}

var _ Node = (*Const)(nil)

func (c *Const) nodeLiteral() {}

func (c *Const) String() string {
	var sb strings.Builder

	if c.Comment != nil {
		c.Comment.WriteTops(&sb)
	}

	sb.WriteString(c.Name.String())

	switch v := c.Value.(type) {
	case *ValueBool:
		if v.IsUserSet {
			sb.WriteString(" = ")
			sb.WriteString(v.String())
		}

	default:
		sb.WriteString(" = ")
		sb.WriteString(v.String())
	}

	if c.Comment != nil {
		sb.WriteString(" ")
		c.Comment.WriteRright(&sb)
		sb.WriteString("\n")
	}

	return sb.String()
}
