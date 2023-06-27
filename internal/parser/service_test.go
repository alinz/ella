package parser_test

import (
	"testing"

	"ella.to/internal/ast"
	"ella.to/internal/parser"
)

func TestParseService(t *testing.T) {
	testCases := TestCases{
		{
			Input:  `service foo {}`,
			Output: `service foo {}`,
		},
		{
			Input: `
	service foo {
		rpc GetFoo() => (value: int64) {}
	}
		`,
			Output: `
service foo {
	rpc GetFoo() => (value: int64)
}
		`,
		},
		{
			Input: `
service foo {
	rpc GetFoo() => (value: int64) {
		required
	}
}
`,
			Output: `
service foo {
	rpc GetFoo() => (value: int64) {
		required
	}
}
`,
		},
	}

	runTests(t, func(p *parser.Parser) (ast.Node, error) {
		return parser.ParseService(p)
	}, testCases)
}
