package scanner_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"ella.to/schema/scanner"
	"ella.to/schema/token"
)

type Tokens []token.Token

type TestCase struct {
	input  string
	output Tokens
}

type TestCases []TestCase

func (t Tokens) String() string {
	var sb strings.Builder
	sb.WriteString("\n")
	for i := range t {
		sb.WriteString(fmt.Sprintf("{Kind: token.%s, Start: %d, End: %d, Val: \"%s\"},\n", t[i].Kind, t[i].Start, t[i].End, t[i].Val))
	}
	return sb.String()
}

func runTestCase(t *testing.T, target int, initState scanner.State, testCases TestCases) {
	if target > -1 && target < len(testCases) {
		testCases = TestCases{testCases[target]}
	}

	for i, tc := range testCases {

		output := make(Tokens, 0)
		emitter := token.EmitterFunc(func(token *token.Token) {
			output = append(output, *token)
		})

		scanner.Start(tc.input, emitter, initState)

		assert.Equal(t, tc.output, output, "Failed scanner at %d: %s", i, output)
	}
}

func TestLex(t *testing.T) {
	runTestCase(t, -1, scanner.Lex, TestCases{
		{
			input: `ella = "1.0.0-b01"`,
			output: Tokens{
				{Kind: token.Word, Start: 0, End: 4, Val: "ella"},
				{Kind: token.Assign, Start: 5, End: 6, Val: "="},
				{Kind: token.ConstantString, Start: 8, End: 17, Val: "1.0.0-b01"},
				{Kind: token.EOF, Start: 18, End: 18, Val: ""},
			},
		},
		{
			input: `message A {
				...B
				...C

				first: int64
			}`,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Word, Start: 8, End: 9, Val: "A"},
				{Kind: token.OpenCurly, Start: 10, End: 11, Val: "{"},
				{Kind: token.Dot, Start: 16, End: 17, Val: "."},
				{Kind: token.Dot, Start: 17, End: 18, Val: "."},
				{Kind: token.Dot, Start: 18, End: 19, Val: "."},
				{Kind: token.Word, Start: 19, End: 20, Val: "B"},
				{Kind: token.Dot, Start: 25, End: 26, Val: "."},
				{Kind: token.Dot, Start: 26, End: 27, Val: "."},
				{Kind: token.Dot, Start: 27, End: 28, Val: "."},
				{Kind: token.Word, Start: 28, End: 29, Val: "C"},
				{Kind: token.Word, Start: 35, End: 40, Val: "first"},
				{Kind: token.Colon, Start: 40, End: 41, Val: ":"},
				{Kind: token.Int64, Start: 42, End: 47, Val: "int64"},
				{Kind: token.CloseCurly, Start: 51, End: 52, Val: "}"},
				{Kind: token.EOF, Start: 52, End: 52, Val: ""},
			},
		},
		{
			input: `enum a int64 {
				one = 1 # comment
				two = 2# comment2
				three
			}`,
			output: Tokens{
				{Kind: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Kind: token.Word, Start: 5, End: 6, Val: "a"},
				{Kind: token.Int64, Start: 7, End: 12, Val: "int64"},
				{Kind: token.OpenCurly, Start: 13, End: 14, Val: "{"},
				{Kind: token.Word, Start: 19, End: 22, Val: "one"},
				{Kind: token.Assign, Start: 23, End: 24, Val: "="},
				{Kind: token.ConstantNumber, Start: 25, End: 26, Val: "1"},
				{Kind: token.Comment, Start: 27, End: 36, Val: "# comment"},
				{Kind: token.Word, Start: 41, End: 44, Val: "two"},
				{Kind: token.Assign, Start: 45, End: 46, Val: "="},
				{Kind: token.ConstantNumber, Start: 47, End: 48, Val: "2"},
				{Kind: token.Comment, Start: 48, End: 58, Val: "# comment2"},
				{Kind: token.Word, Start: 63, End: 68, Val: "three"},
				{Kind: token.CloseCurly, Start: 72, End: 73, Val: "}"},
				{Kind: token.EOF, Start: 73, End: 73, Val: ""},
			},
		},
		{
			input: `enum a int64 {
				one = 1
				two = 2
				three
			}`,
			output: Tokens{
				{Kind: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Kind: token.Word, Start: 5, End: 6, Val: "a"},
				{Kind: token.Int64, Start: 7, End: 12, Val: "int64"},
				{Kind: token.OpenCurly, Start: 13, End: 14, Val: "{"},
				{Kind: token.Word, Start: 19, End: 22, Val: "one"},
				{Kind: token.Assign, Start: 23, End: 24, Val: "="},
				{Kind: token.ConstantNumber, Start: 25, End: 26, Val: "1"},
				{Kind: token.Word, Start: 31, End: 34, Val: "two"},
				{Kind: token.Assign, Start: 35, End: 36, Val: "="},
				{Kind: token.ConstantNumber, Start: 37, End: 38, Val: "2"},
				{Kind: token.Word, Start: 43, End: 48, Val: "three"},
				{Kind: token.CloseCurly, Start: 52, End: 53, Val: "}"},
				{Kind: token.EOF, Start: 53, End: 53, Val: ""},
			},
		},
		{
			input: `enum a int64 {}`,
			output: Tokens{
				{Kind: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Kind: token.Word, Start: 5, End: 6, Val: "a"},
				{Kind: token.Int64, Start: 7, End: 12, Val: "int64"},
				{Kind: token.OpenCurly, Start: 13, End: 14, Val: "{"},
				{Kind: token.CloseCurly, Start: 14, End: 15, Val: "}"},
				{Kind: token.EOF, Start: 15, End: 15, Val: ""},
			},
		},
		{
			input: `a=1`,
			output: Tokens{
				{Kind: token.Word, Start: 0, End: 1, Val: "a"},
				{Kind: token.Assign, Start: 1, End: 2, Val: "="},
				{Kind: token.ConstantNumber, Start: 2, End: 3, Val: "1"},
				{Kind: token.EOF, Start: 3, End: 3, Val: ""},
			},
		},
	})
}
