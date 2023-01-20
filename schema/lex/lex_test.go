package lex_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alinz/rpc.go/pkg/lexer"
	"github.com/alinz/rpc.go/schema/lex/token"
)

type Tokens []lexer.Token

type TestCase struct {
	input  string
	output Tokens
}

type TestCases []TestCase

func (t Tokens) String() string {
	var sb strings.Builder
	sb.WriteString("\n")
	for i, _ := range t {
		sb.WriteString(fmt.Sprintf("{Type: token.%s, Start: %d, End: %d, Val: \"%s\"},\n", token.Name(t[i].Type), t[i].Start, t[i].End, t[i].Val))
	}
	return sb.String()
}

func runTestCase(t *testing.T, target int, initState lexer.State, testCases TestCases) {
	if target > -1 && target < len(testCases) {
		testCases = TestCases{testCases[target]}
	}

	for i, tc := range testCases {

		output := make(Tokens, 0)
		emitter := lexer.EmitterFunc(func(token *lexer.Token) {
			output = append(output, *token)
		})

		lexer.Start(tc.input, emitter, initState)

		assert.Equal(t, tc.output, output, "Failed lexer at %d: %s", i, output)
	}
}
