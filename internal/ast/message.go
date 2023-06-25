package ast

import (
	"strings"

	"ella.to/internal/token"
)

type Field struct {
	Name    *Identifier
	Type    Type
	Options []*Const
}

var _ Node = (*Field)(nil)

func (f *Field) nodeLiteral() {}
func (f *Field) String() string {
	var sb strings.Builder

	sb.WriteString(f.Name.String())
	sb.WriteString(": ")
	sb.WriteString(f.Type.String())

	if len(f.Options) > 0 {
		sb.WriteString(" {")
		for _, o := range f.Options {
			sb.WriteString("\n\t\t")
			sb.WriteString(o.String())
		}
		sb.WriteString("\n\t}")
	}

	return sb.String()
}

type Message struct {
	Token   *token.Token
	Name    *Identifier
	Extends []*CustomType
	Fields  []*Field
	Comment *Comment
}

var _ Node = (*Message)(nil)

func (m *Message) nodeLiteral() {}
func (m *Message) String() string {
	var sb strings.Builder

	sb.WriteString("message ")
	sb.WriteString(m.Name.String())
	sb.WriteString(" {")

	for _, e := range m.Extends {
		sb.WriteString("\n\t...")
		sb.WriteString(e.String())
	}

	for _, f := range m.Fields {
		sb.WriteString("\n\t")
		sb.WriteString(f.String())
	}

	if len(m.Extends) > 0 || len(m.Fields) > 0 {
		sb.WriteString("\n")
	}

	sb.WriteString("}")

	return sb.String()
}
