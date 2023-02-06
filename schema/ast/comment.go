package ast

import (
	"strings"

	"github.com/alinz/rpc.go/schema/token"
)

type CommentSide int

const (
	_ CommentSide = iota
	CommentSideRight
	CommentSideTop
)

type Comment struct {
	Token  *token.Token
	Values []string
	Side   CommentSide
}

func (c *Comment) TokenLiteral() string {
	var sb strings.Builder

	for _, v := range c.Values {
		sb.WriteString("#")
		sb.WriteString(v)
		sb.WriteString("\n")
	}

	return sb.String()
}
