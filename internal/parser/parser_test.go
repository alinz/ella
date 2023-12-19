package parser_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"compiler.ella.to/internal/ast"
	"compiler.ella.to/internal/parser"
)

type TestCase struct {
	Input  string
	Output string
	Error  string
}

type TestCases []TestCase

func runTests(t *testing.T, fn func(*parser.Parser) (ast.Node, error), testCases TestCases) {
	for _, testCase := range testCases {
		p := parser.New(testCase.Input)

		node, err := fn(p)
		if err != nil {
			if testCase.Error == "" {
				t.Fatal(err)
			} else {
				assert.Equal(t, strings.TrimSpace(testCase.Error), strings.TrimSpace(err.Error()))
				continue
			}
		}

		assert.Equal(t, strings.TrimSpace(testCase.Output), strings.TrimSpace(node.String()))
	}
}
