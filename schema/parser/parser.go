package parser

import (
	"fmt"
	"strings"

	"ella.to/schema/ast"
	"ella.to/schema/scanner"
	"ella.to/schema/token"
)

type Parser struct {
	input     string
	tokens    token.Iterator
	nextToken *token.Token
}

func (p *Parser) scanToken() {
	// for now all comment token will be ignored and skipped
	for {
		p.nextToken = p.tokens.NextToken()
		if p.nextToken.Kind == token.Comment {
			continue
		}
		break
	}
}

func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{
		Nodes: []ast.Node{},
	}

	for p.nextToken.Kind != token.EOF {
		var node ast.Node
		var err error

		switch p.nextToken.Kind {
		case token.Enum:
			node, err = p.parseEnum()
		case token.Message:
			node, err = p.parseMessage()
		case token.Service:
			node, err = p.parseService()
		case token.Word:
			node, err = p.parseConstant(false)
		}

		if err != nil {
			return nil, err
		}
		program.Nodes = append(program.Nodes, node)
	}

	return program, nil
}

func (p *Parser) ShowContext(token *token.Token, lines int) string {
	start := token.Start
	end := token.End

	for i := 0; i < lines; i++ {
		if start > 0 {
			start = strings.LastIndex(p.input[:start], "\n") - 1
		}

		if end < len(p.input)-1 {
			end = strings.Index(p.input[end+1:], "\n") + end + 1
			if end == -1 {
				end = len(p.input)
			}
		}
	}

	return p.input[start:end]

}

func New(input string) *Parser {
	tokenEmitter := token.NewEmitterIterator()
	go scanner.Start(input, tokenEmitter, scanner.Lex)
	parser := &Parser{input: input, tokens: tokenEmitter}
	parser.scanToken()
	return parser
}

func isIdentifier(word string, isFirstCharCapital bool) bool {
	if word == "" {
		return false
	}
	if isFirstCharCapital && !(word[0] >= 'A' && word[0] <= 'Z') {
		return false
	}
	if !isFirstCharCapital && !(word[0] >= 'a' && word[0] <= 'z') {
		return false
	}
	return true
}

func mustBeNameFor(tok *token.Token, name string, firstChatCaptial bool) error {
	if tok.Kind != token.Word {
		return fmt.Errorf("expected a name for %s but got %s", name, tok.Kind)
	}

	if !(isIdentifier(tok.Val, firstChatCaptial)) {
		return fmt.Errorf("invalid %s name %s", name, tok.Val)
	}

	return nil
}
