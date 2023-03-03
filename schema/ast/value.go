package ast

import (
	"strconv"
	"strings"

	"github.com/alinz/ella.to/schema/token"
)

type Value interface {
	valueNode()
	Node
}

type ValueString struct {
	Token   *token.Token
	Content string
}

func (v ValueString) valueNode() {}
func (v *ValueString) TokenLiteral() string {
	var sb strings.Builder

	if strings.Contains(v.Token.Val, `"`) {
		sb.WriteString(`'`)
		sb.WriteString(v.Token.Val)
		sb.WriteString(`'`)
	} else {
		sb.WriteString(`"`)
		sb.WriteString(v.Token.Val)
		sb.WriteString(`"`)
	}

	return sb.String()
}

type ValueInt struct {
	Token   *token.Token
	Size    int // 8 | 16 | 32 | 64
	Content int64
}

func (v ValueInt) valueNode() {}
func (v *ValueInt) TokenLiteral() string {
	return strconv.FormatInt(v.Content, 10)
}

type ValueUint struct {
	Token   *token.Token
	Size    int // 8 | 16 | 32 | 64
	Content uint64
}

func (v ValueUint) valueNode() {}
func (v *ValueUint) TokenLiteral() string {
	return strconv.FormatUint(v.Content, 10)
}

type ValueFloat struct {
	Token   *token.Token
	Size    int // 32 | 64
	Content float64
}

func (v ValueFloat) valueNode() {}
func (v *ValueFloat) TokenLiteral() string {
	return v.Token.Val
}

type ValueBool struct {
	Token   *token.Token
	Content bool
}

func (v ValueBool) valueNode() {}
func (v *ValueBool) TokenLiteral() string {
	return v.Token.Val
}

func ParseValue(token *token.Token) Value {
	{
		value, err := strconv.ParseInt(token.Val, 10, 64)
		if err == nil {
			return &ValueInt{
				Token:   token,
				Content: value,
			}
		}
	}

	{
		value, err := strconv.ParseFloat(token.Val, 64)
		if err == nil {
			return &ValueFloat{
				Token:   token,
				Content: value,
			}
		}
	}

	{
		value, err := strconv.ParseBool(token.Val)
		if err == nil {
			return &ValueBool{
				Token:   token,
				Content: value,
			}
		}
	}

	return &ValueString{
		Token:   token,
		Content: token.Val,
	}
}
