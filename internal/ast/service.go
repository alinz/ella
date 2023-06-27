package ast

import (
	"strings"

	"ella.to/internal/token"
)

type Arg struct {
	Name *Identifier
	Type Type
}

var _ Node = (*Arg)(nil)

func (a *Arg) nodeLiteral() {}
func (a *Arg) String() string {
	var sb strings.Builder

	sb.WriteString(a.Name.String())
	sb.WriteString(": ")
	sb.WriteString(a.Type.String())

	return sb.String()
}

type Return struct {
	Name *Identifier
	Type Type
}

var _ Node = (*Return)(nil)

func (r *Return) nodeLiteral() {}
func (r *Return) String() string {
	var sb strings.Builder

	sb.WriteString(r.Name.String())
	sb.WriteString(": ")
	sb.WriteString(r.Type.String())

	return sb.String()
}

type Method struct {
	Type    *token.Token // rpc or http token
	Name    *Identifier
	Args    []*Arg
	Returns []*Return
	Options []*Const
}

var _ Node = (*Method)(nil)

func (m *Method) nodeLiteral() {}
func (m *Method) String() string {
	var sb strings.Builder

	sb.WriteString(m.Type.Val)
	sb.WriteString(" ")
	sb.WriteString(m.Name.String())
	sb.WriteString("(")

	for i, a := range m.Args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(a.String())
	}

	sb.WriteString(")")

	if len(m.Returns) > 0 {
		sb.WriteString(" => (")
		for i, r := range m.Returns {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(r.String())
		}
		sb.WriteString(")")
	}

	if len(m.Options) > 0 {
		sb.WriteString(" {")
		for _, o := range m.Options {
			sb.WriteString("\n\t\t")
			sb.WriteString(o.String())
		}
		sb.WriteString("\n\t}")
	}

	return sb.String()
}

type Service struct {
	Token   *token.Token
	Name    *Identifier
	Methods []*Method
	Comment *Comment
}

var _ Node = (*Service)(nil)

func (r *Service) nodeLiteral() {}
func (r *Service) String() string {
	var sb strings.Builder

	sb.WriteString(r.Token.Val)
	sb.WriteString(" ")
	sb.WriteString(r.Name.String())
	sb.WriteString(" {")
	for _, m := range r.Methods {
		sb.WriteString("\n\t")
		sb.WriteString(m.String())
	}

	if len(r.Methods) > 0 {
		sb.WriteString("\n")
	}

	sb.WriteString("}")

	return sb.String()
}
