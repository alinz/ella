package lex_test

import (
	"testing"

	"github.com/alinz/rpc.go/schema/lex"
	"github.com/alinz/rpc.go/schema/lex/token"
)

func TestLexStmt(t *testing.T) {
	runTestCase(t, -1, lex.Stmt(nil), TestCases{
		{
			input: `
			hello = 1
			bye = 2

			enum Cool int64 {
				A
				B
				C = 2
			}
			`,
			output: Tokens{
				{Type: token.Identifier, Start: 4, End: 9, Val: "hello"},
				{Type: token.Assign, Start: 10, End: 11, Val: "="},
				{Type: token.Value, Start: 12, End: 13, Val: "1"},
				{Type: token.Identifier, Start: 17, End: 20, Val: "bye"},
				{Type: token.Assign, Start: 21, End: 22, Val: "="},
				{Type: token.Value, Start: 23, End: 24, Val: "2"},
				{Type: token.Enum, Start: 29, End: 33, Val: "enum"},
				{Type: token.Identifier, Start: 34, End: 38, Val: "Cool"},
				{Type: token.Type, Start: 39, End: 44, Val: "int64"},
				{Type: token.OpenCurl, Start: 45, End: 46, Val: "{"},
				{Type: token.Identifier, Start: 51, End: 52, Val: "A"},
				{Type: token.Identifier, Start: 57, End: 58, Val: "B"},
				{Type: token.Identifier, Start: 63, End: 64, Val: "C"},
				{Type: token.Assign, Start: 65, End: 66, Val: "="},
				{Type: token.Value, Start: 67, End: 68, Val: "2"},
				{Type: token.CloseCurl, Start: 72, End: 73, Val: "}"},
				{Type: token.EOF, Start: 77, End: 77, Val: ""},
			},
		},
		{
			input: `hello = 1
			bye = 2
			`,
			output: Tokens{
				{Type: token.Identifier, Start: 0, End: 5, Val: "hello"},
				{Type: token.Assign, Start: 6, End: 7, Val: "="},
				{Type: token.Value, Start: 8, End: 9, Val: "1"},
				{Type: token.Identifier, Start: 13, End: 16, Val: "bye"},
				{Type: token.Assign, Start: 17, End: 18, Val: "="},
				{Type: token.Value, Start: 19, End: 20, Val: "2"},
				{Type: token.EOF, Start: 24, End: 24, Val: ""},
			},
		},
		{
			input: ``,
			output: Tokens{
				{Type: token.EOF, Start: 0, End: 0, Val: ""},
			},
		},
	})

}
