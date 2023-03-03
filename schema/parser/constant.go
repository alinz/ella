package parser

import (
	"fmt"

	"github.com/alinz/ella.to/schema/ast"
	"github.com/alinz/ella.to/schema/token"
)

// T -> ID = V
// OR
// T -> ID = V
// T -> ID
func (p *Parser) parseConstant(permitEmpty bool) (*ast.Constant, error) {
	constant := &ast.Constant{}

	if p.nextToken.Kind != token.Word {
		return nil, fmt.Errorf("expected a name for constant but got %s", p.nextToken.Kind)
	}

	if !isConstantName(p.nextToken.Val) {
		return nil, fmt.Errorf("invalid constant name %s", p.nextToken.Val)
	}

	constant.Name = &ast.Identifier{
		Name:  p.nextToken.Val,
		Token: p.nextToken,
	}

	p.scanToken()

	if permitEmpty && p.nextToken.Kind != token.Assign {
		return constant, nil
	} else if p.nextToken.Kind != token.Assign {
		return nil, fmt.Errorf("expected '=' after constant name but got %s", p.nextToken.Kind)
	}

	p.scanToken() // skip '='

	if !p.nextToken.OneOf(token.ConstantNumber, token.ConstantString, token.Word) {
		return nil, fmt.Errorf("expected value for constant but got %s", p.nextToken.Kind)
	}

	constant.Value = ast.ParseValue(p.nextToken)

	p.scanToken()

	return constant, nil
}

func isConstantName(word string) bool {
	return isIdentifier(word, false) || isIdentifier(word, true)
}
