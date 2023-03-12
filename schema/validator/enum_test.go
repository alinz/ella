package validator_test

import "testing"

func TestValidatorEnum(t *testing.T) {
	runTests(t, TestCases{
		{
			Input: `
hello = 1
Ella = "0.0.1"

enum MyEnum int32 {
	A = 1
	B
}
			`,
			Output: `
Ella = "0.0.1"
hello = 1

enum MyEnum int32 {
	A = 1
	B = 2
}
`,
		},
		{
			Input: `
		hello = 1
		Ella = "0.0.1"

		enum MyEnum int32 {
			A
			B
		}
					`,
			Output: `
Ella = "0.0.1"
hello = 1

enum MyEnum int32 {
	A = 0
	B = 1
}
		`,
		},
	})
}
