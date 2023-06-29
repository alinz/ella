package parser

import (
	"fmt"
	"strings"

	"ella.to/internal/scanner"
	"ella.to/internal/token"
)

type Parser struct {
	input   string
	tokens  token.Iterator
	nextTok *token.Token
	currTok *token.Token
}

func (p *Parser) Current() *token.Token {
	return p.currTok
}

func (p *Parser) Next() *token.Token {
	if p.nextTok != nil {
		p.currTok = p.nextTok
		p.nextTok = nil
	} else {
		p.currTok = p.tokens.NextToken()
	}

	return p.currTok
}

func (p *Parser) Peek() *token.Token {
	if p.nextTok == nil {
		p.nextTok = p.tokens.NextToken()
	}

	return p.nextTok
}

func (p *Parser) WithError(token *token.Token, args ...any) error {
	var sb strings.Builder
	for i, arg := range args {
		sb.WriteString(fmt.Sprintf("%v", arg))
		if i != len(args)-1 {
			sb.WriteString(": ")
		}
	}

	return fmt.Errorf("%s\n%s", sb.String(), p.showContext(token, 5))
}

func (p *Parser) showContext(token *token.Token, lines int) string {
	start := token.Start
	end := token.End

	for i := 0; i < lines; i++ {
		if start > 0 {
			start = strings.LastIndex(p.input[:start], "\n") - 1
		}

		if end < len(p.input)-1 {
			end = strings.Index(p.input[end:], "\n") + end + 1
			if end == -1 {
				end = len(p.input)
			}
		}
	}

	if start < 0 {
		start = 0
	}

	return p.input[start:end]
}

func New(input string) *Parser {
	tokenEmitter := token.NewEmitterIterator()
	go scanner.Start(input, tokenEmitter, scanner.Lex)
	return &Parser{input: input, tokens: tokenEmitter}
}