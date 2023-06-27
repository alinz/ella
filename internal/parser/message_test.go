package parser_test

import (
	"testing"

	"ella.to/internal/ast"
	"ella.to/internal/parser"
)

func TestParseMessage(t *testing.T) {
	testCases := TestCases{
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
	}

	runTests(t, func(p *parser.Parser) (ast.Node, error) {
		return parser.ParseMessage(p)
	}, testCases)
}
