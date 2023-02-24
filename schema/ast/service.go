package ast

import "strings"

type Arg struct {
	Name *Identifier
	Type Type
}

var _ Node = (*Arg)(nil)

func (a *Arg) TokenLiteral() string {
	var sb strings.Builder

	sb.WriteString(a.Name.TokenLiteral())
	sb.WriteString(": ")
	sb.WriteString(a.Type.TokenLiteral())

	return sb.String()
}

type Return struct {
	Name   *Identifier
	Stream bool
	Type   Type
}

var _ Node = (*Return)(nil)

func (r *Return) TokenLiteral() string {
	var sb strings.Builder

	sb.WriteString(r.Name.TokenLiteral())
	sb.WriteString(": ")
	if r.Stream {
		sb.WriteString("stream ")
	}
	sb.WriteString(r.Type.TokenLiteral())

	return sb.String()
}

type Method struct {
	Name    *Identifier
	Options []*Constant
	Args    []*Arg
	Returns []*Return
}

var _ Node = (*Method)(nil)

func (m *Method) TokenLiteral() string {
	var sb strings.Builder

	sb.WriteString(m.Name.TokenLiteral())
	sb.WriteString("(")
	for i, a := range m.Args {
		sb.WriteString(a.TokenLiteral())
		if i < len(m.Args)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")

	if len(m.Returns) == 0 {
		if len(m.Options) > 0 {
			sb.WriteString(" {")
			for _, c := range m.Options {
				sb.WriteString("\n\t\t")
				sb.WriteString(c.TokenLiteral())
			}

			if len(m.Options) > 0 {
				sb.WriteString("\n\t")
			}

			sb.WriteString("}")
		}

		return sb.String()
	}

	sb.WriteString(" => (")
	for i, r := range m.Returns {
		sb.WriteString(r.TokenLiteral())
		if i < len(m.Returns)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")

	if len(m.Options) > 0 {
		sb.WriteString(" {")
		for _, c := range m.Options {
			sb.WriteString("\n\t\t")
			sb.WriteString(c.TokenLiteral())
		}

		if len(m.Options) > 0 {
			sb.WriteString("\n\t")
		}

		sb.WriteString("}")
	}

	return sb.String()
}

type Service struct {
	Name    *Identifier
	Methods []*Method
}

var _ Node = (*Service)(nil)

func (s *Service) TokenLiteral() string {
	var sb strings.Builder

	sb.WriteString("service ")
	sb.WriteString(s.Name.TokenLiteral())
	sb.WriteString(" {")
	for _, m := range s.Methods {
		sb.WriteString("\n\t")
		sb.WriteString(m.TokenLiteral())
	}

	if len(s.Methods) > 0 {
		sb.WriteString("\n")
	}

	sb.WriteString("}")

	return sb.String()
}
