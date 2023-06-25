package parser_test

import (
	"testing"

	"ella.to/internal/ast"
	"ella.to/internal/parser"
)

func TestParseConst(t *testing.T) {
	runTests(t, func(p *parser.Parser) (ast.Node, error) {
		return parser.ParseConst(p)
	}, TestCases{
		{
			Input:  `ella = "1.0.0-b01"`,
			Output: `ella = "1.0.0-b01"`,
		},
		{
			Input:  `ella = 'cool this is a string'`,
			Output: `ella = 'cool this is a string'`,
		},
		{
			Input:  `ella = 1.33333`,
			Output: `ella = 1.33333`,
		},
		{
			Input:  `ella`,
			Output: ``,
			Error: `
expected '=' after an identifier for defining a constant
ella
`,
		},
	})
}
