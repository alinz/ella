package lex_test

import (
	"testing"

	"github.com/alinz/rpc.go/schema/lex"
	"github.com/alinz/rpc.go/schema/lex/token"
)

func TestLexEnum(t *testing.T) {
	runTestCase(t, -1, lex.Enum(nil), TestCases{
		{
			input: `enum a int32 { }`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Type: token.Identifier, Start: 5, End: 6, Val: "a"},
				{Type: token.Type, Start: 7, End: 12, Val: "int32"},
				{Type: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Type: token.CloseCurl, Start: 15, End: 16, Val: "}"},
			},
		},
		{
			input: `enum a int32{}`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Type: token.Identifier, Start: 5, End: 6, Val: "a"},
				{Type: token.Type, Start: 7, End: 12, Val: "int32"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.CloseCurl, Start: 13, End: 14, Val: "}"},
			},
		},
		{
			input: `enum a int32{
				a
			}`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Type: token.Identifier, Start: 5, End: 6, Val: "a"},
				{Type: token.Type, Start: 7, End: 12, Val: "int32"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Type: token.CloseCurl, Start: 23, End: 24, Val: "}"},
			},
		},
		{
			input: `enum a int32{
				a = 1
			}`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Type: token.Identifier, Start: 5, End: 6, Val: "a"},
				{Type: token.Type, Start: 7, End: 12, Val: "int32"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Type: token.Assign, Start: 20, End: 21, Val: "="},
				{Type: token.Value, Start: 22, End: 23, Val: "1"},
				{Type: token.CloseCurl, Start: 27, End: 28, Val: "}"},
			},
		},
		{
			input: `enum a int32{
				a = 1
				b
				c = 2
			}`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Type: token.Identifier, Start: 5, End: 6, Val: "a"},
				{Type: token.Type, Start: 7, End: 12, Val: "int32"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Type: token.Assign, Start: 20, End: 21, Val: "="},
				{Type: token.Value, Start: 22, End: 23, Val: "1"},
				{Type: token.Identifier, Start: 28, End: 29, Val: "b"},
				{Type: token.Identifier, Start: 34, End: 35, Val: "c"},
				{Type: token.Assign, Start: 36, End: 37, Val: "="},
				{Type: token.Value, Start: 38, End: 39, Val: "2"},
				{Type: token.CloseCurl, Start: 43, End: 44, Val: "}"},
			},
		},
	})
}
