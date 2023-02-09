package validator_test

import "testing"

func TestValidatorEnum(t *testing.T) {
	runTests(t, TestCases{
		{
			Input: `
hello = 1
rpc = 0.0.1

enum MyEnum Int32 {
	A = 1
	B
}
			`,
			Output: `
rpc = 0.0.1
hello = 1

enum MyEnum Int32 {
	A = 1
	B = 2
}
`,
		},
	})
}
