package parser_test

import "testing"

func TestParseConst(t *testing.T) {
	runTests(t, TestCases{
		{
			Input:  `ella = "1.0.0-b01"`,
			Output: `ella = "1.0.0-b01"`,
		},
		{
			Input:  `ella = 'cool this is a string'`,
			Output: `ella = 'cool this is a string'`,
		},
		{
			Input:  `ella = 1.33333`,
			Output: `ella = 1.33333`,
		},
		{
			Input: `
# a
b = 1 # b




# hello
a = 2 # a
# world
`,
			Output: `
# a
b = 1 # b

# hello
a = 2 # a

# world
`,
		},
	})
}
