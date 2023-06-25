package parser_test

import (
	"testing"

	"ella.to/internal/ast"
	"ella.to/internal/parser"
)

func TestParseMessage(t *testing.T) {
	testCases := TestCases{
		{
			Input:  `message foo {}`,
			Output: `message foo {}`,
		},
		{
			Input: `message foo {
				first_name: string {
					required
				}
			}`,
			Output: `
message foo {
	first_name: string {
		required
	}
}
`,
		},
	}

	runTests(t, func(p *parser.Parser) (ast.Node, error) {
		return parser.ParseMessage(p)
	}, testCases)
}
