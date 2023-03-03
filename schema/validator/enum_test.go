package validator_test

import "testing"

func TestValidatorEnum(t *testing.T) {
	runTests(t, TestCases{
		{
			Input: `
hello = 1
ella = "0.0.1"

enum MyEnum int32 {
	A = 1
	B
}
			`,
			Output: `
ella = "0.0.1"
hello = 1

enum MyEnum int32 {
	A = 1
	B = 2
}
`,
		},
	})
}
