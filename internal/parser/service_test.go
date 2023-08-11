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
		{
			Input: `
service Foo {
	rpc GetFoo() => (value: int64) {
		Required
		A = 1
	}
}
`,
			Output: `
service Foo {
	rpc GetFoo() => (value: int64) {
		Required
		A = 1
	}
}
`,
		},
		{
			Input: `
service Foo {
	rpc GetFoo() => (value: int64) {
		Required
		A = 1mb
		B = 1ms
	}
}
`,
			Output: `
service Foo {
	rpc GetFoo() => (value: int64) {
		Required
		A = 1mb
		B = 1ms
	}
}
`,
		},
		{
			Input: `
service Foo {
	rpc GetFoo() => (value: stream int64)
}
`,
			Output: `
service Foo {
	rpc GetFoo() => (value: stream int64)
}
`,
		},
	}

	runTests(t, func(p *parser.Parser) (ast.Node, error) {
		return parser.ParseService(p)
	}, testCases)
}
