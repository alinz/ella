package lex_test

import (
	"testing"

	"github.com/alinz/rpc.go/schema/lex"
	"github.com/alinz/rpc.go/schema/lex/token"
)

func TestLexConstant(t *testing.T) {
	runTestCase(t, -1, lex.Constant(nil), TestCases{
		{
			input: `a = 1`,
			output: Tokens{
				{Type: token.Identifier, Start: 0, End: 1, Val: "a"},
				{Type: token.Assign, Start: 2, End: 3, Val: "="},
				{Type: token.Value, Start: 4, End: 5, Val: "1"},
			},
		},
		{
			input: `a =1`,
			output: Tokens{
				{Type: token.Identifier, Start: 0, End: 1, Val: "a"},
				{Type: token.Assign, Start: 2, End: 3, Val: "="},
				{Type: token.Value, Start: 3, End: 4, Val: "1"},
			},
		},
		{
			input: `a=1`,
			output: Tokens{
				{Type: token.Identifier, Start: 0, End: 1, Val: "a"},
				{Type: token.Assign, Start: 1, End: 2, Val: "="},
				{Type: token.Value, Start: 2, End: 3, Val: "1"},
			},
		},
		{
			input: `a=       1           `,
			output: Tokens{
				{Type: token.Identifier, Start: 0, End: 1, Val: "a"},
				{Type: token.Assign, Start: 1, End: 2, Val: "="},
				{Type: token.Value, Start: 9, End: 10, Val: "1"},
			},
		},
	})
}
