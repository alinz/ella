package parser_test

import (
	"testing"

	"ella.to/internal/ast"
	"ella.to/internal/parser"
)

func TestParseEnum(t *testing.T) {
	testCases := TestCases{
		{
			Input: `enum Foo {}`,
			Error: `
expected enum type
enum Foo {			
			`,
		},
		{
			Input: `
enum Foo int8 {
	A = 1
	B
	C
}
			`,
			Output: `
enum Foo int8 {
	A = 1
	B
	C
}			
			`,
		},
		{
			Input: `

					enum Foo int8 {
						A = 1
					}

					`,
			Output: `
enum Foo int8 {
	A = 1
}
		`,
		},
		{
			Input: `enum Foo int8 {

					}`,
			Output: `enum Foo int8 {}`,
		},
	}

	runTests(t, func(p *parser.Parser) (ast.Node, error) {
		return parser.ParseEnum(p)
	}, testCases)
}
