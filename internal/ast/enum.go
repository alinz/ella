package ast

import (
	"fmt"
	"strings"

	"ella.to/internal/token"
)

type EnumSet struct {
	Name    *Identifier `json:"name"`
	Value   *ValueInt   `json:"value"`
	Defined bool        `json:"defined"`
}

var _ Node = (*EnumSet)(nil)

func (e *EnumSet) TokenLiteral() string {
	return e.Name.TokenLiteral()
}

func (e *EnumSet) String() string {
	var sb strings.Builder

	sb.WriteString(e.Name.String())
	if e.Defined {
		sb.WriteString(" = ")
		sb.WriteString(fmt.Sprintf("%d", e.Value.Value))
	}

	return sb.String()
}

type Enum struct {
	Token *token.Token
	Name  *Identifier
	Size  int // 8, 16, 32, 64 selected by compiler based on the largest and smallest values
	Sets  []*EnumSet
}

var _ Statement = (*Enum)(nil)

func (e *Enum) statementLiteral() {}

func (e *Enum) TokenLiteral() string {
	return e.Token.Literal
}

func (e *Enum) String() string {
	var sb strings.Builder

	sb.WriteString("enum ")
	sb.WriteString(e.Name.String())
	sb.WriteString(" {")

	for _, set := range e.Sets {
		sb.WriteString("\n\t")
		sb.WriteString(set.String())
	}

	if len(e.Sets) > 0 {
		sb.WriteString("\n")
	}

	sb.WriteString("}")

	return sb.String()
}
