package parser_test

import (
	"testing"
)

func TestConstantParser(t *testing.T) {
	runTests(t, TestCases{
		{
			Input:  `rpc = 1.0.0-b01`,
			Output: `rpc = 1.0.0-b01`,
		},
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
