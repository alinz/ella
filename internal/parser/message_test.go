package parser_test

import (
	"testing"

	"ella.to/internal/ast"
	"ella.to/internal/parser"
)

func TestParseMessage(t *testing.T) {
	testCases := TestCases{
		{
			Input: `message Foo {
				...Hello
			}`,
			Output: `
message Foo {
	...Hello
}`,
		},
		{
			Input:  `message Foo {}`,
			Output: `message Foo {}`,
		},
		{
			Input: `message Foo {
				FirstName: string {
					Required
				}
			}`,
			Output: `
message Foo {
	FirstName: string {
		Required
	}
}
`,
		},
		{
			Input: `message Foo {
				FirstName: string {
					Required = true
				}
			}`,
			Output: `
message Foo {
	FirstName: string {
		Required = true
	}
}
`,
		},
	}

	runTests(t, func(p *parser.Parser) (ast.Node, error) {
		return parser.ParseMessage(p)
	}, testCases)
}
