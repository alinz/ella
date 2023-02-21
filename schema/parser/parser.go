package parser

import (
	"github.com/alinz/rpc.go/schema/ast"
	"github.com/alinz/rpc.go/schema/lexer"
	"github.com/alinz/rpc.go/schema/token"
)

type Parser struct {
	tokens    token.Iterator
	nextToken *token.Token
}

func (p *Parser) scanToken() {
	// for now all comment token will be ignored and skipped
	for {
		p.nextToken = p.tokens.NextToken()
		if p.nextToken.Kind == token.Comment {
			// the next token must be a value
			p.tokens.NextToken()
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
		case token.Identifier:
			node, err = p.parseConstant(false)
		}

		if err != nil {
			return nil, err
		}
		program.Nodes = append(program.Nodes, node)
	}

	return program, nil
}

func New(input string) *Parser {
	tokenEmitter := token.NewEmitterIterator()
	go lexer.Start(input, tokenEmitter, lexer.Stmt(nil))
	parser := &Parser{tokens: tokenEmitter}
	parser.scanToken()
	return parser
}
