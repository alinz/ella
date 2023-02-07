package lexer_test

import (
	"testing"

	"github.com/alinz/rpc.go/schema/lexer"
	"github.com/alinz/rpc.go/schema/token"
)

func TestLexEnum(t *testing.T) {
	runTestCase(t, -1, lexer.Enum(nil), TestCases{
		{
			input: `enum a int32 { a = 1 b = 2 }`,
			output: Tokens{
				{Kind: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Kind: token.Identifier, Start: 5, End: 6, Val: "a"},
				{Kind: token.Type, Start: 7, End: 12, Val: "int32"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.Identifier, Start: 15, End: 16, Val: "a"},
				{Kind: token.Assign, Start: 17, End: 18, Val: "="},
				{Kind: token.Value, Start: 19, End: 20, Val: "1"},
				{Kind: token.Identifier, Start: 21, End: 22, Val: "b"},
				{Kind: token.Assign, Start: 23, End: 24, Val: "="},
				{Kind: token.Value, Start: 25, End: 26, Val: "2"},
				{Kind: token.CloseCurl, Start: 27, End: 28, Val: "}"},
			},
		},
		{
			input: `enum a int32 { }`,
			output: Tokens{
				{Kind: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Kind: token.Identifier, Start: 5, End: 6, Val: "a"},
				{Kind: token.Type, Start: 7, End: 12, Val: "int32"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.CloseCurl, Start: 15, End: 16, Val: "}"},
			},
		},
		{
			input: `enum a int32{}`,
			output: Tokens{
				{Kind: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Kind: token.Identifier, Start: 5, End: 6, Val: "a"},
				{Kind: token.Type, Start: 7, End: 12, Val: "int32"},
				{Kind: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Kind: token.CloseCurl, Start: 13, End: 14, Val: "}"},
			},
		},
		{
			input: `enum a int32{
				a
			}`,
			output: Tokens{
				{Kind: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Kind: token.Identifier, Start: 5, End: 6, Val: "a"},
				{Kind: token.Type, Start: 7, End: 12, Val: "int32"},
				{Kind: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Kind: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Kind: token.CloseCurl, Start: 23, End: 24, Val: "}"},
			},
		},
		{
			input: `enum a int32{
				a = 1
			}`,
			output: Tokens{
				{Kind: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Kind: token.Identifier, Start: 5, End: 6, Val: "a"},
				{Kind: token.Type, Start: 7, End: 12, Val: "int32"},
				{Kind: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Kind: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Kind: token.Assign, Start: 20, End: 21, Val: "="},
				{Kind: token.Value, Start: 22, End: 23, Val: "1"},
				{Kind: token.CloseCurl, Start: 27, End: 28, Val: "}"},
			},
		},
		{
			input: `enum a int32{
				a = 1
				b
				c = 2
			}`,
			output: Tokens{
				{Kind: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Kind: token.Identifier, Start: 5, End: 6, Val: "a"},
				{Kind: token.Type, Start: 7, End: 12, Val: "int32"},
				{Kind: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Kind: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Kind: token.Assign, Start: 20, End: 21, Val: "="},
				{Kind: token.Value, Start: 22, End: 23, Val: "1"},
				{Kind: token.Identifier, Start: 28, End: 29, Val: "b"},
				{Kind: token.Identifier, Start: 34, End: 35, Val: "c"},
				{Kind: token.Assign, Start: 36, End: 37, Val: "="},
				{Kind: token.Value, Start: 38, End: 39, Val: "2"},
				{Kind: token.CloseCurl, Start: 43, End: 44, Val: "}"},
			},
		},
	})
}
