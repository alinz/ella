package parser_test

import (
	"testing"
)

func TestConstantParser(t *testing.T) {
	runTests(t, TestCases{
		{
			Input: `a = 1 b= 2`,
			Output: `
a = 1
b = 2
`,
		},
		{
			Input:  `a = 1`,
			Output: `a = 1`,
		},
	})
}
