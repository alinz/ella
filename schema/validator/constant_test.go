package validator_test

import "testing"

func TestValidatorConstant(t *testing.T) {
	runTests(t, TestCases{
		{
			Input: `hello = 1
			Ella = "0.0.1"`,
			Output: `
Ella = "0.0.1"
hello = 1
`,
		},
	})
}
