package parser

import (
	"fmt"
	"os"
	"strings"

	"compiler.ella.to/internal/scanner"
	"compiler.ella.to/internal/token"
)

type Parser struct {
	input   string
	tokens  token.Iterator
	nextTok *token.Token
	currTok *token.Token

	errorCodesMap  map[int64]struct{}
	errorCodeValue int64
}

func (p *Parser) setErrorCodeValue(value int64) error {
	if _, ok := p.errorCodesMap[value]; ok {
		return fmt.Errorf("error code %d is already defined", value)
	}
	p.errorCodesMap[value] = struct{}{}
	return nil
}

func (p *Parser) getNextErrorCode() int64 {
	for {
		if _, ok := p.errorCodesMap[p.errorCodeValue]; !ok {
			p.errorCodesMap[p.errorCodeValue] = struct{}{}
			return p.errorCodeValue
		}

		p.errorCodeValue++
	}
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

	return fmt.Errorf("%s: ->%s<-\n%s", sb.String(), token.Literal, p.showContext(token, 5))
}

func (p *Parser) showContext(token *token.Token, lines int) string {
	start := token.Start
	end := token.End

	if p.input == "" {
		content, err := os.ReadFile(token.Filename)
		if err != nil {
			return err.Error()
		}
		p.input = string(content)
	}

	// going backwards n lines
	for i := 0; i < lines; i++ {
		if start == -1 {
			break
		}

		start = strings.LastIndex(p.input[:start], "\n")
	}

	for i := 0; i < lines; i++ {
		if start > 0 {
			start = strings.LastIndex(p.input[:start], "\n")
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

	return p.input[start:end] + "\n"
}

func New(input string) *Parser {
	tokenEmitter := token.NewEmitterIterator()
	go scanner.Start(tokenEmitter, scanner.Lex, input)
	return &Parser{input: input, tokens: tokenEmitter, errorCodesMap: make(map[int64]struct{}), errorCodeValue: 1000}
}

func NewFilenames(filenames ...string) *Parser {
	tokenEmitter := token.NewEmitterIterator()
	go scanner.StartWithFilenames(tokenEmitter, scanner.Lex, filenames...)
	return &Parser{tokens: tokenEmitter, errorCodesMap: make(map[int64]struct{}), errorCodeValue: 1000}
}
