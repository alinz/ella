package parser_test

import (
	"testing"

	"ella.to/internal/ast"
	"ella.to/internal/parser"
)

func TestParseService(t *testing.T) {
	testCases := TestCases{
		{
			Input:  `service Foo {}`,
			Output: `service Foo {}`,
		},
		{
			Input: `
	service Foo {
		rpc GetFoo() => (value: int64) {}
	}
		`,
			Output: `
service Foo {
	rpc GetFoo() => (value: int64)
}
		`,
		},
		{
			Input: `
service Foo {
	rpc GetFoo() => (value: int64) {
		Required
	}
}
`,
			Output: `
service Foo {
	rpc GetFoo() => (value: int64) {
		Required
	}
}
`,
		},
	}

	runTests(t, func(p *parser.Parser) (ast.Node, error) {
		return parser.ParseService(p)
	}, testCases)
}
