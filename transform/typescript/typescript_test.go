package typescript_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"ella.to/schema/ast"
	"ella.to/schema/parser"
	"ella.to/transform"
)

type TestCase struct {
	Input  string
	Fn     func(prog *ast.Program) []transform.Func
	Output string
}

type TestCases []TestCase

func runTests(t *testing.T, testCases TestCases) {
	for _, testCase := range testCases {
		p := parser.New(testCase.Input)

		program, err := p.Parse()
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		var buffer bytes.Buffer
		err = transform.Run(&buffer, testCase.Fn(program)...)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		if !assert.Equal(
			t,
			strings.TrimSpace(testCase.Output),
			strings.TrimSpace(buffer.String()),
		) {
			t.FailNow()
		}
	}
}

func getValues[T any](prog *ast.Program) []T {
	var values []T
	for _, node := range prog.Nodes {
		if value, ok := node.(T); ok {
			values = append(values, value)
		}
	}
	return values
}
