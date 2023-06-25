package ast

import (
	"strings"

	"ella.to/internal/token"
)

type Comment struct {
	Tops  []*token.Token
	Right *token.Token
	// Bottom is used to collect all comments
	// when they are inside of a block and at the end
	// for example:
	//
	// enum foo int8 {
	//  a = 1
	// 	# comment
	// }
	Bottom []*token.Token
}

var _ Node = (*Comment)(nil)

func (c *Comment) nodeLiteral() {}
func (c *Comment) String() string {
	var sb strings.Builder

	c.WriteTops(&sb)
	c.WriteRright(&sb)
	c.WriteBottom(&sb)

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

func (c *Comment) WriteBottom(sb *strings.Builder) {
	for _, t := range c.Bottom {
		sb.WriteString("# ")
		sb.WriteString(strings.TrimSpace(t.Val))
	}
}

func (c *Comment) HasAny() bool {
	return len(c.Tops) > 0 || c.Right != nil
}
