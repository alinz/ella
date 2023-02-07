package parser

import (
	"fmt"

	"github.com/alinz/rpc.go/schema/ast"
	"github.com/alinz/rpc.go/schema/token"
)

// T -> ID = V
// OR
// T -> ID = V
// T -> ID
func (p *Parser) parseConstant(permitEmpty bool) (*ast.Constant, error) {
	constant := &ast.Constant{}

	if p.nextToken.Kind != token.Identifier {
		return nil, fmt.Errorf("expected identifier but got %s", p.nextToken.Kind)
	}

	constant.Name = &ast.Identifier{
		Name:  p.nextToken.Val,
		Token: p.nextToken,
	}

	p.scanToken()

	if permitEmpty && p.nextToken.Kind != token.Assign {
		return constant, nil
	} else if p.nextToken.Kind != token.Assign {
		return nil, fmt.Errorf("expected '=' but got %s", p.nextToken.Kind)
	}

	p.scanToken() // skip '='

	if p.nextToken.Kind != token.Value {
		return nil, fmt.Errorf("expected value but got %s", p.nextToken.Kind)
	}

	constant.Value = ast.ParseValue(p.nextToken)

	p.scanToken()

	return constant, nil
}
