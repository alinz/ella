package parser_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alinz/rpc.go/schema/parser"
)

type TestCase struct {
	Input  string
	Output string
}

type TestCases []TestCase

func runTests(t *testing.T, testCases TestCases) {
	for _, testCase := range testCases {
		p := parser.New(testCase.Input)

		program, err := p.Parse()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, strings.TrimSpace(testCase.Output), strings.TrimSpace(program.TokenLiteral()))
	}
}
