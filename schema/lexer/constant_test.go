package lexer_test

import (
	"testing"

	"github.com/alinz/rpc.go/schema/lexer"
	"github.com/alinz/rpc.go/schema/token"
)

func TestLexConstant(t *testing.T) {
	runTestCase(t, -1, lexer.Constant(nil), TestCases{
		{
			input: `c="ali's macbook"`,
			output: Tokens{
				{Kind: token.Identifier, Start: 0, End: 1, Val: "c"},
				{Kind: token.Assign, Start: 1, End: 2, Val: "="},
				{Kind: token.Value, Start: 3, End: 16, Val: "ali's macbook"},
			},
		},
		{
			input: `b = '2'`,
			output: Tokens{
				{Kind: token.Identifier, Start: 0, End: 1, Val: "b"},
				{Kind: token.Assign, Start: 2, End: 3, Val: "="},
				{Kind: token.Value, Start: 5, End: 6, Val: "2"},
			},
		},
		{
			input: `a = "1"`,
			output: Tokens{
				{Kind: token.Identifier, Start: 0, End: 1, Val: "a"},
				{Kind: token.Assign, Start: 2, End: 3, Val: "="},
				{Kind: token.Value, Start: 5, End: 6, Val: "1"},
			},
		},
		{
			input: `a = 1`,
			output: Tokens{
				{Kind: token.Identifier, Start: 0, End: 1, Val: "a"},
				{Kind: token.Assign, Start: 2, End: 3, Val: "="},
				{Kind: token.Value, Start: 4, End: 5, Val: "1"},
			},
		},
		{
			input: `a =1`,
			output: Tokens{
				{Kind: token.Identifier, Start: 0, End: 1, Val: "a"},
				{Kind: token.Assign, Start: 2, End: 3, Val: "="},
				{Kind: token.Value, Start: 3, End: 4, Val: "1"},
			},
		},
		{
			input: `a=1`,
			output: Tokens{
				{Kind: token.Identifier, Start: 0, End: 1, Val: "a"},
				{Kind: token.Assign, Start: 1, End: 2, Val: "="},
				{Kind: token.Value, Start: 2, End: 3, Val: "1"},
			},
		},
		{
			input: `a=       1           `,
			output: Tokens{
				{Kind: token.Identifier, Start: 0, End: 1, Val: "a"},
				{Kind: token.Assign, Start: 1, End: 2, Val: "="},
				{Kind: token.Value, Start: 9, End: 10, Val: "1"},
			},
		},
	})
}