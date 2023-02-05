package lexer_test

import (
	"testing"

	"github.com/alinz/rpc.go/schema/lexer"
	"github.com/alinz/rpc.go/schema/token"
)

func TestLexStmt(t *testing.T) {
	runTestCase(t, -1, lexer.Stmt(nil), TestCases{
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
				{Kind: token.Identifier, Start: 4, End: 9, Val: "hello"},
				{Kind: token.Assign, Start: 10, End: 11, Val: "="},
				{Kind: token.Value, Start: 12, End: 13, Val: "1"},
				{Kind: token.Identifier, Start: 17, End: 20, Val: "bye"},
				{Kind: token.Assign, Start: 21, End: 22, Val: "="},
				{Kind: token.Value, Start: 23, End: 24, Val: "2"},
				{Kind: token.Enum, Start: 29, End: 33, Val: "enum"},
				{Kind: token.Identifier, Start: 34, End: 38, Val: "Cool"},
				{Kind: token.Type, Start: 39, End: 44, Val: "int64"},
				{Kind: token.OpenCurl, Start: 45, End: 46, Val: "{"},
				{Kind: token.Identifier, Start: 51, End: 52, Val: "A"},
				{Kind: token.Identifier, Start: 57, End: 58, Val: "B"},
				{Kind: token.Identifier, Start: 63, End: 64, Val: "C"},
				{Kind: token.Assign, Start: 65, End: 66, Val: "="},
				{Kind: token.Value, Start: 67, End: 68, Val: "2"},
				{Kind: token.CloseCurl, Start: 72, End: 73, Val: "}"},
				{Kind: token.EOF, Start: 77, End: 77, Val: ""},
			},
		},
		{
			input: `hello = 1
			bye = 2
			`,
			output: Tokens{
				{Kind: token.Identifier, Start: 0, End: 5, Val: "hello"},
				{Kind: token.Assign, Start: 6, End: 7, Val: "="},
				{Kind: token.Value, Start: 8, End: 9, Val: "1"},
				{Kind: token.Identifier, Start: 13, End: 16, Val: "bye"},
				{Kind: token.Assign, Start: 17, End: 18, Val: "="},
				{Kind: token.Value, Start: 19, End: 20, Val: "2"},
				{Kind: token.EOF, Start: 24, End: 24, Val: ""},
			},
		},
		{
			input: ``,
			output: Tokens{
				{Kind: token.EOF, Start: 0, End: 0, Val: ""},
			},
		},
	})

}
