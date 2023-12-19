package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"compiler.ella.to/internal/token"
)

type Arg struct {
	Name     *Identifier `json:"name"`
	Type     Type        `json:"type"`
	Optional bool        `json:"optional"`
}

var _ Node = (*Arg)(nil)

func (a *Arg) MarshalText() ([]byte, error) {
	var buff bytes.Buffer

	buff.WriteString(`{"name":"`)
	buff.WriteString(a.Name.TokenLiteral())
	buff.WriteString(`","type":`)
	typ, err := marshalTextType(a.Type)
	if err != nil {
		return nil, err
	}
	buff.Write(typ)
	buff.WriteString(`,"optional":`)
	if a.Optional {
		buff.WriteString(`true`)
	} else {
		buff.WriteString(`false`)
	}
	buff.WriteString(`}`)

	return buff.Bytes(), nil
}

func (a *Arg) UnmarshalText(text []byte) error {
	results := struct {
		Name     *Identifier     `json:"name"`
		Type     json.RawMessage `json:"type"`
		Optional bool            `json:"optional"`
	}{}

	if err := json.Unmarshal(text, &results); err != nil {
		return err
	}

	a.Name = results.Name
	a.Optional = results.Optional

	return unmarshalTextType(results.Type, &a.Type)
}

func (a *Arg) TokenLiteral() string {
	return a.Name.TokenLiteral()
}

func (a *Arg) String() string {
	var sb strings.Builder

	sb.WriteString(a.Name.String())
	if a.Optional {
		sb.WriteString("?")
	}
	sb.WriteString(": ")
	sb.WriteString(a.Type.String())

	return sb.String()
}

type Args []*Arg

func (a Args) String() string {
	if len(a) == 0 {
		return ""
	}

	var sb strings.Builder

	for i, arg := range a {
		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(arg.String())
	}

	return sb.String()
}

type Return struct {
	Name   *Identifier `json:"name"`
	Type   Type        `json:"type"`
	Stream bool        `json:"stream"`
}

var _ Node = (*Return)(nil)

func (r *Return) MarshalText() ([]byte, error) {
	var buff bytes.Buffer

	buff.WriteString(`{"name":"`)
	buff.WriteString(r.Name.TokenLiteral())
	buff.WriteString(`","type":`)
	typ, err := marshalTextType(r.Type)
	if err != nil {
		return nil, err
	}
	buff.Write(typ)
	buff.WriteString(`,"stream":`)
	if r.Stream {
		buff.WriteString(`true`)
	} else {
		buff.WriteString(`false`)
	}
	buff.WriteString(`}`)

	return buff.Bytes(), nil
}

func (r *Return) UnmarshalText(text []byte) error {
	results := struct {
		Name   *Identifier     `json:"name"`
		Type   json.RawMessage `json:"type"`
		Stream bool            `json:"optional"`
	}{}

	if err := json.Unmarshal(text, &results); err != nil {
		return err
	}

	r.Name = results.Name
	r.Stream = results.Stream

	return unmarshalTextType(results.Type, &r.Type)
}

func (r *Return) TokenLiteral() string {
	return r.Name.TokenLiteral()
}

func (r *Return) String() string {
	var sb strings.Builder

	sb.WriteString(r.Name.String())
	sb.WriteString(": ")
	if r.Stream {
		sb.WriteString("stream ")
	}
	sb.WriteString(r.Type.String())

	return sb.String()
}

type Returns []*Return

func (r Returns) String() string {
	if len(r) == 0 {
		return ""
	}

	var sb strings.Builder

	for i, ret := range r {
		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(ret.String())
	}

	return sb.String()
}

type MethodType int

const (
	_          MethodType = iota
	MethodRPC             // rpc
	MethodHTTP            // http
)

func (m MethodType) String() string {
	switch m {
	case MethodRPC:
		return "rpc"
	case MethodHTTP:
		return "http"
	}

	panic("unknown method type")
}

func (m MethodType) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *MethodType) UnmarshalText(text []byte) error {
	switch string(text) {
	case "rpc":
		*m = MethodRPC
	case "http":
		*m = MethodHTTP
	default:
		return fmt.Errorf("unknown type: %s", text)
	}

	return nil
}

type Method struct {
	Type    MethodType  `json:"type"` // rpc, http
	Name    *Identifier `json:"name"`
	Args    Args        `json:"args"`
	Returns Returns     `json:"returns"`
	Options Options     `json:"options"`
}

type Methods []*Method

func (m Methods) String() string {
	if len(m) == 0 {
		return ""
	}

	var sb strings.Builder

	for _, method := range m {
		sb.WriteString("\n\t")
		sb.WriteString(method.Type.String())
		sb.WriteString(" ")
		sb.WriteString(method.Name.String())
		sb.WriteString("(")
		sb.WriteString(method.Args.String())
		sb.WriteString(")")

		if len(method.Returns) == 0 {
			sb.WriteString("\n")
			continue
		}

		sb.WriteString(" => (")
		sb.WriteString(method.Returns.String())
		sb.WriteString(")")
		sb.WriteString(method.Options.String(2))

		sb.WriteString("\n")
	}

	return sb.String()
}

type Service struct {
	Token   *token.Token `json:"token"`
	Name    *Identifier  `json:"name"`
	Methods Methods      `json:"methods"`
}

var _ Statement = (*Service)(nil)

func (s *Service) statementLiteral() {}

func (s *Service) TokenLiteral() string {
	return s.Token.Literal
}

func (s *Service) String() string {
	var sb strings.Builder

	sb.WriteString("service ")
	sb.WriteString(s.Name.String())
	sb.WriteString(" {")
	sb.WriteString(s.Methods.String())
	sb.WriteString("}")

	return sb.String()
}
