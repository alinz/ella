package typescript_test

import (
	"testing"

	"ella.to/schema/ast"
	"ella.to/transform"
	"ella.to/transform/typescript"
)

func TestConstants(t *testing.T) {
	runTests(t, TestCases{
		{
			Input: `
Ella = "0.0.1"
A = 1
B = true
			`,
			Fn: func(prog *ast.Program) []transform.Func {
				constants := getValues[*ast.Constant](prog)

				return []transform.Func{
					typescript.Constants(constants),
				}
			},
			Output: `
// Constants

export const Ella = "0.0.1";
export const A = 1;
export const B = true;
			`,
		},
	})
}
