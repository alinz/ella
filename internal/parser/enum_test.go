package parser_test

import (
	"testing"

	"ella.to/internal/ast"
	"ella.to/internal/parser"
)

func TestParseEnum(t *testing.T) {
	testCases := TestCases{
		{
			Input: `enum foo {}`,
			Error: `
expected enum type
enum foo {			
			`,
		},
		{
			Input: `
enum foo int8 {
	a = 1
	b
	c
}
			`,
			Output: `
enum foo int8 {
	a = 1
	b
	c
}			
			`,
		},
		{
			Input: `

					enum foo int8 {
						a = 1
					}

					`,
			Output: `
enum foo int8 {
	a = 1
}
		`,
		},
		{
			Input: `enum foo int8 {

					}`,
			Output: `enum foo int8 {}`,
		},
	}

	runTests(t, func(p *parser.Parser) (ast.Node, error) {
		return parser.ParseEnum(p)
	}, testCases)
}
