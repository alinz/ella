package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Program struct {
	Statements []Statement `json:"statements"`
}

var _ Node = (*Program)(nil)

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var sb strings.Builder

	for i, stmt := range p.Statements {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(stmt.String())
	}

	return sb.String()
}

func (p *Program) MarshalText() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("[")
	for i, stmt := range p.Statements {
		if i != 0 {
			buf.WriteString(",")
		}

		buf.WriteString(`{`)

		buf.WriteString(`"type":"`)
		buf.WriteString(stmt.TokenLiteral())
		buf.WriteString(`","data":`)

		b, err := json.Marshal(stmt)
		if err != nil {
			return nil, err
		}
		buf.Write(b)

		buf.WriteString(`}`)
	}
	buf.WriteString("]")

	return buf.Bytes(), nil
}

func (p *Program) UnmarshalText(text []byte) error {
	results := []struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}{}

	if err := json.Unmarshal(text, &results); err != nil {
		return err
	}

	statements := make([]Statement, 0)

	for _, result := range results {
		var stmt Statement

		switch result.Type {
		case "const":
			stmt = &Const{}
		case "model":
			stmt = &Model{}
		case "service":
			stmt = &Service{}
		default:
			return fmt.Errorf("unknown type: %s", result.Type)
		}

		if err := json.Unmarshal(result.Data, stmt); err != nil {
			return err
		}
		statements = append(statements, stmt)
	}

	p.Statements = statements
	return nil
}
