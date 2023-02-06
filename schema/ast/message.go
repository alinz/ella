package ast

import (
	"strings"
)

type Field struct {
	Name    *Identifier
	Type    Type
	Options []*Constant
}

var _ Node = (*Field)(nil)

func (f *Field) TokenLiteral() string {
	var sb strings.Builder

	sb.WriteString(f.Name.TokenLiteral())
	sb.WriteString(" ")
	sb.WriteString(f.Type.TokenLiteral())

	if len(f.Options) > 0 {
		sb.WriteString(" {")
		for _, o := range f.Options {
			sb.WriteString("\n")
			sb.WriteString(o.TokenLiteral())
		}
		sb.WriteString("\n}")
	}

	return sb.String()
}

type Message struct {
	Name   *Identifier
	Fields []*Field
}

var _ Node = (*Message)(nil)

func (m *Message) TokenLiteral() string {
	var sb strings.Builder

	sb.WriteString("message ")
	sb.WriteString(m.Name.TokenLiteral())
	sb.WriteString(" {\n")
	for _, f := range m.Fields {
		sb.WriteString(f.TokenLiteral())
		sb.WriteString("\n")
	}
	sb.WriteString("}")

	return sb.String()
}
