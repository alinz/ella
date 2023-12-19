package parser_test

import (
	"testing"

	"compiler.ella.to/internal/ast"
	"compiler.ella.to/internal/parser"
)

func TestParseEnum(t *testing.T) {
	testCases := TestCases{
		{
			Input: `enum Foo {}`,
			Output: `
enum Foo {}
			`,
		},
		{
			Input: `
		enum Foo {
			A = 1
			B
			C
		}
					`,
			Output: `
enum Foo {
	A = 1
	B
	C
}
					`,
		},
		{
			Input: `

							enum Foo {
								A = 1
							}

							`,
			Output: `
enum Foo {
	A = 1
}
				`,
		},
		{
			Input: `enum Foo {

					}`,
			Output: `enum Foo {}`,
		},
	}

	runTests(t, func(p *parser.Parser) (ast.Node, error) {
		return parser.ParseEnum(p)
	}, testCases)
}
