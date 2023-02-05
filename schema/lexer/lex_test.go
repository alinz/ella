package lexer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alinz/rpc.go/schema/lexer"
	"github.com/alinz/rpc.go/schema/token"
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

func runTestCase(t *testing.T, target int, initState lexer.StateFn, testCases TestCases) {
	if target > -1 && target < len(testCases) {
		testCases = TestCases{testCases[target]}
	}

	for i, tc := range testCases {

		output := make(Tokens, 0)
		emitter := token.EmitterFunc(func(token *token.Token) {
			output = append(output, *token)
		})

		lexer.Start(tc.input, emitter, initState)

		assert.Equal(t, tc.output, output, "Failed lexer at %d: %s", i, output)
	}
}
