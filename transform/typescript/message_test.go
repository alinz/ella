package typescript_test

import (
	"testing"

	"ella.to/schema/ast"
	"ella.to/transform"
	"ella.to/transform/typescript"
)

func TestMessages(t *testing.T) {
	runTests(t, TestCases{
		{
			Input: `
Ella = "0.0.1"

message Complex {
	Firstname: string
	Lastname: string

	Age: int8
	Jobs: []string
	AddressMap: map<string,string>
}

message A {
	Firstname: string {
		json = first_name
	}
	Lastname: string
}
			`,
			Fn: func(prog *ast.Program) []transform.Func {
				messages := getValues[*ast.Message](prog)

				return []transform.Func{
					typescript.Messages(messages),
				}
			},
			Output: `
// Messages

export interface Complex {
	firstname: string
	lastname: string
	age: number
	jobs: string[]
	address_map: { [key: string]: string }
}

export interface A {
	first_name: string
	lastname: string
}			
			`,
		},
		{
			Input: `
Ella = "0.0.1"

message A {
	Firstname: string {
		json = first_name
	}
	Lastname: string
}
			`,
			Fn: func(prog *ast.Program) []transform.Func {
				messages := getValues[*ast.Message](prog)

				return []transform.Func{
					typescript.Messages(messages),
				}
			},
			Output: `
// Messages

export interface A {
	first_name: string
	lastname: string
}			
			`,
		},
	})
}
