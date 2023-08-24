package ast

import (
	"bytes"
	"encoding/json"
	"strings"

	"ella.to/internal/token"
)

type Field struct {
	Name    *Identifier `json:"name"`
	Type    Type        `json:"type"`
	Options Options     `json:"options"`
}

var _ Node = (*Field)(nil)

func (f *Field) TokenLiteral() string {
	return f.Name.Token.Literal
}

func (f *Field) String() string {
	var sb strings.Builder

	sb.WriteString(f.Name.String())
	sb.WriteString(": ")
	sb.WriteString(f.Type.String())
	sb.WriteString(f.Options.String(2))

	return sb.String()
}

func (f *Field) MarshalText() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(`{"name":"`)
	buf.WriteString(f.Name.Token.Literal)
	buf.WriteString(`","type":`)
	typ, err := marshalTextType(f.Type)
	if err != nil {
		return nil, err
	}
	buf.Write(typ)
	buf.WriteString(`,"options":`)
	opt, err := json.Marshal(f.Options)
	if err != nil {
		return nil, err
	}
	buf.Write(opt)
	buf.WriteString(`}`)

	return buf.Bytes(), nil
}

func (f *Field) UnmarshalText(text []byte) error {
	results := struct {
		Name    *Identifier     `json:"name"`
		Type    json.RawMessage `json:"type"`
		Options Options         `json:"options"`
	}{}

	if err := json.Unmarshal(text, &results); err != nil {
		return err
	}

	f.Name = results.Name
	f.Options = results.Options

	return unmarshalTextType(results.Type, &f.Type)
}

type Fields []*Field

func (f Fields) String() string {
	if len(f) == 0 {
		return ""
	}

	var sb strings.Builder

	for _, field := range f {
		sb.WriteString("\n\t")
		sb.WriteString(field.String())
	}

	sb.WriteString("\n")

	return sb.String()
}

type Message struct {
	Token   *token.Token  `json:"token"`
	Name    *Identifier   `json:"name"`
	Extends []*Identifier `json:"extends"`
	Fields  Fields        `json:"fields"`
}

var _ Statement = (*Message)(nil)

func (m *Message) statementLiteral() {}

func (m *Message) TokenLiteral() string {
	return m.Token.Literal
}

func (m *Message) String() string {
	var sb strings.Builder

	sb.WriteString("message ")
	sb.WriteString(m.Name.String())
	sb.WriteString(" {")

	for _, extend := range m.Extends {
		sb.WriteString("\n\t...")
		sb.WriteString(extend.String())
	}

	if len(m.Extends) > 0 && len(m.Fields) == 0 {
		sb.WriteString("\n")
	}

	sb.WriteString(m.Fields.String())
	sb.WriteString("}")

	return sb.String()
}
