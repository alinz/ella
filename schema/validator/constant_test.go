package validator_test

import "testing"

func TestValidatorConstant(t *testing.T) {
	runTests(t, TestCases{
		{
			Input: `hello = 1
			ella = "0.0.1"`,
			Output: `
ella = "0.0.1"
hello = 1
`,
		},
	})
}
