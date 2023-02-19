package ast

import (
	"strings"

	"github.com/alinz/rpc.go/schema/token"
)

type Comment struct {
	Token *token.Token
	Value string
}

func (c *Comment) TokenLiteral() string {
	var sb strings.Builder

	sb.WriteString("#")
	sb.WriteString(c.Value)
	sb.WriteString("\n")

	return sb.String()
}
