package scanner_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"ella.to/internal/scanner"
	"ella.to/internal/token"
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
		sb.WriteString(fmt.Sprintf("{Type: token.%s, Start: %d, End: %d, Literal: \"%s\"},\n", t[i].Type, t[i].Start, t[i].End, t[i].Literal))
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

		scanner.Start(emitter, initState, tc.input)
		assert.Equal(t, tc.output, output, "Failed scanner at %d: %s", i, output)
	}
}
