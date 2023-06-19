package parser

import (
	"fmt"
	"strings"

	"ella.to/internal/ast"
	"ella.to/internal/scanner"
	"ella.to/internal/token"
)

type Parser struct {
	input   string
	tokens  token.Iterator
	currTok *token.Token
}

func (p *Parser) scanToken() {
	p.currTok = p.tokens.NextToken()
}

func (p *Parser) Parse() (*ast.Program, error) {
	comment := newComment()
	program := &ast.Program{}

	var prevNode ast.Node

	for p.currTok.Type != token.EOF {
		var err error
		var node ast.Node

		switch p.currTok.Type {
		case token.Identifier:
			node, err = p.parseConst(false)
			attachComment(node, comment)
			prevNode = node
		case token.Enum:
		case token.Message:
		case token.Service:
		case token.Error:
			err = p.newError(p.currTok, p.currTok.Val)
		case token.RightComment:
			comment.Right = p.currTok
			attachComment(prevNode, comment)
			comment = newComment()
		case token.TopComment:
			comment.Tops = append(comment.Tops, p.currTok)
		default:
			err = p.newError(p.currTok, "unexpected token")
		}

		if err != nil {
			return nil, err
		}

		if node != nil {
			program.Nodes = append(program.Nodes, node)
		}

		if comment.HasAny() {
			program.Nodes = append(program.Nodes, comment)
			comment = newComment()
		}

		p.scanToken()
	}

	return program, nil
}

func (p *Parser) newError(token *token.Token, args ...any) error {
	var sb strings.Builder
	for i, arg := range args {
		sb.WriteString(fmt.Sprintf("%v", arg))
		if i != len(args)-1 {
			sb.WriteString(": ")
		}
	}

	return fmt.Errorf("%s\n%s", sb.String(), p.ShowContext(token, 5))
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
