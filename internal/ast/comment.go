package ast

import (
	"strings"

	"ella.to/internal/token"
)

type Comment struct {
	Tops  []*token.Token
	Right *token.Token
}

var _ Node = (*Comment)(nil)

func (c *Comment) nodeLiteral() {}
func (c *Comment) String() string {
	var sb strings.Builder

	c.WriteTops(&sb)
	c.WriteRright(&sb)

	return sb.String()
}

func (c *Comment) WriteTops(sb *strings.Builder) {
	for _, t := range c.Tops {
		sb.WriteString("# ")
		sb.WriteString(strings.TrimSpace(t.Val))
	}
}

func (c *Comment) WriteRright(sb *strings.Builder) {
	if c.Right != nil {
		sb.WriteString("# ")
		sb.WriteString(strings.TrimSpace(c.Right.Val))
	}
}

func (c *Comment) HasAny() bool {
	return len(c.Tops) > 0 || c.Right != nil
}
