package parser_test

import "testing"

func TestEnumParser(t *testing.T) {
	runTests(t, TestCases{
		{
			Input:  `enum A int64 {}`,
			Output: `enum A int64 {}`,
		},
		{
			Input: `enum A int64 { A = 1 }`,
			Output: `
enum A int64 {
	A = 1
}`,
		},
		{
			Input: `enum A int64 { A = 1 B C }`,
			Output: `
enum A int64 {
	A = 1
	B
	C
}`,
		},
		{
			Input: `enum A int64 { A = 1 B 
				C = 4 }`,
			Output: `
enum A int64 {
	A = 1
	B
	C = 4
}`,
		},
	})
}
