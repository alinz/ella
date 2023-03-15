package typescript_test

import (
	"testing"

	"ella.to/schema/ast"
	"ella.to/transform"
	"ella.to/transform/typescript"
)

func TestEnums(t *testing.T) {
	runTests(t, TestCases{
		{
			Input: `
Ella = "0.0.1"
enum A int8 {}
			`,
			Fn: func(prog *ast.Program) []transform.Func {
				enums := getValues[*ast.Enum](prog)

				return []transform.Func{
					typescript.Enums(enums),
				}
			},
			Output: `
// Enums

export enum A {
}
			`,
		},
		{
			Input: `
Ella = "0.0.1"
enum A int8 {
	A = 1
	B
	C
}
			`,
			Fn: func(prog *ast.Program) []transform.Func {
				enums := getValues[*ast.Enum](prog)

				return []transform.Func{
					typescript.Enums(enums),
				}
			},
			Output: `
// Enums

export enum A {
	A = "A"
	B = "B"
	C = "C"
}
			`,
		},
		{
			Input: `
Ella = "0.0.1"

enum B int8 {
	
}

enum A int8 {
	A = 1
	B
	C
}
			`,
			Fn: func(prog *ast.Program) []transform.Func {
				enums := getValues[*ast.Enum](prog)

				return []transform.Func{
					typescript.Enums(enums),
				}
			},
			Output: `
// Enums

export enum B {
}

export enum A {
	A = "A"
	B = "B"
	C = "C"
}
			`,
		},
	})
}
