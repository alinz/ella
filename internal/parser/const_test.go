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
			Input:  `const Integer = 1`,
			Output: `const Integer = 1`,
		},
		{
			Input:  `const Ella = "1.0.0-b01"`,
			Output: `const Ella = "1.0.0-b01"`,
		},
		{
			Input:  `const Ella = 'cool this is a string'`,
			Output: `const Ella = 'cool this is a string'`,
		},
		{
			Input:  `const Ella = 1.33333`,
			Output: `const Ella = 1.33333`,
		},
		{
			Input:  `const Ella`,
			Output: ``,
			Error: `
expected '=' after an identifier for defining a constant: -><-
const Ella
`,
		},
		{
			Input: `
			const Ella = version`,
			Output: `
const Ella = version
			`,
		},
	})
}
