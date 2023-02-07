package parser_test

import "testing"

func TestEnumParser(t *testing.T) {
	runTests(t, TestCases{
		{
			Input:  `enum A Int64 {}`,
			Output: `enum A Int64 {}`,
		},
		{
			Input: `enum A Int64 { A = 1 }`,
			Output: `
enum A Int64 {
	A = 1
}`,
		},
		{
			Input: `enum A Int64 { A = 1 B C }`,
			Output: `
enum A Int64 {
	A = 1
	B
	C
}`,
		},
		{
			Input: `enum A Uint8 { A = 1 B 
				C = 4 }`,
			Output: `
enum A Uint8 {
	A = 1
	B
	C = 4
}`,
		},
	})
}
