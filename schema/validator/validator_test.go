package validator_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alinz/rpc.go/schema/parser"
	"github.com/alinz/rpc.go/schema/validator"
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
		assert.NoError(t, err)

		err = validator.Validate(program)
		assert.NoError(t, err)

		assert.Equal(t, strings.TrimSpace(testCase.Output), strings.TrimSpace(program.TokenLiteral()))
	}
}
