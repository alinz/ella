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
			Input:  `Ella = "1.0.0-b01"`,
			Output: `Ella = "1.0.0-b01"`,
		},
		{
			Input:  `Ella = 'cool this is a string'`,
			Output: `Ella = 'cool this is a string'`,
		},
		{
			Input:  `Ella = 1.33333`,
			Output: `Ella = 1.33333`,
		},
		{
			Input:  `Ella`,
			Output: ``,
			Error: `
expected '=' after an identifier for defining a constant
Ella
`,
		},
		{
			Input: `
			Ella = version`,
			Output: `
Ella = version
			`,
		},
	})
}
