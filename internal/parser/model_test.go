package parser_test

import (
	"testing"

	"ella.to/internal/ast"
	"ella.to/internal/parser"
)

func TestParseMessage(t *testing.T) {
	testCases := TestCases{
		{
			Input: `model Foo {
				...Hello
			}`,
			Output: `
model Foo {
	...Hello
}`,
		},
		{
			Input:  `model Foo {}`,
			Output: `model Foo {}`,
		},
		{
			Input: `model Foo {
				FirstName: string {
					Required
				}
			}`,
			Output: `
model Foo {
	FirstName: string {
		Required
	}
}
`,
		},
		{
			Input: `model Foo {
				FirstName: string {
					Required = true
				}
			}`,
			Output: `
model Foo {
	FirstName: string {
		Required = true
	}
}
`,
		},
	}

	runTests(t, func(p *parser.Parser) (ast.Node, error) {
		return parser.ParseModel(p)
	}, testCases)
}
