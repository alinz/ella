package ast

import (
	"bytes"
	"encoding/json"
	"strings"
)

type Option struct {
	Name  *Identifier `json:"name"`
	Value Value       `json:"value"`
}

var _ Node = (*Option)(nil)

func (o *Option) MarshalText() ([]byte, error) {
	var buff bytes.Buffer

	buff.WriteString(`{"name":"`)
	buff.WriteString(o.Name.TokenLiteral())
	buff.WriteString(`","value":`)
	val, err := marshalTextValue(o.Value)
	if err != nil {
		return nil, err
	}
	buff.Write(val)
	buff.WriteString(`}`)

	return buff.Bytes(), nil
}

func (o *Option) UnmarshalText(text []byte) error {
	results := struct {
		Name  *Identifier     `json:"name"`
		Value json.RawMessage `json:"value"`
	}{}

	if err := json.Unmarshal(text, &results); err != nil {
		return err
	}

	o.Name = results.Name

	return unmarshalTextValue(results.Value, &o.Value)
}

func (o *Option) String() string {
	var sb strings.Builder

	sb.WriteString(o.Name.String())

	showValue := true

	switch v := o.Value.(type) {
	case *ValueBool:
		showValue = v.Defined
	}

	if showValue {
		sb.WriteString(" = ")
		sb.WriteString(o.Value.String())
	}

	return sb.String()
}

func (o *Option) TokenLiteral() string {
	return o.Name.TokenLiteral()
}

type Options []*Option

func (o Options) String(tabs int) string {
	if len(o) == 0 {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(" {")

	for _, opt := range o {
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("\t", tabs))
		sb.WriteString(opt.String())
	}

	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("\t", tabs-1))
	sb.WriteString("}")

	return sb.String()
}
