package lex_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alinz/rpc.go/pkg/lexer"
	"github.com/alinz/rpc.go/schema/lex"
	"github.com/alinz/rpc.go/schema/lex/token"
)

type Tokens []lexer.Token

type TestCase struct {
	input  string
	output Tokens
}

type TestCases []TestCase

func (t Tokens) String() string {
	var sb strings.Builder
	sb.WriteString("\n")
	for i, _ := range t {
		sb.WriteString(fmt.Sprintf("{Type: token.%s, Start: %d, End: %d, Val: \"%s\"},\n", token.Name(t[i].Type), t[i].Start, t[i].End, t[i].Val))
	}
	return sb.String()
}

func runTestCase(t *testing.T, initState lexer.State, testCases TestCases, indicies ...int) {
	if len(indicies) > 0 {
		result := make(TestCases, 0)
		for _, index := range indicies {
			result = append(result, testCases[index])
		}
		testCases = result
	}

	for i, tc := range testCases {

		output := make(Tokens, 0)
		emitter := lexer.EmitterFunc(func(token *lexer.Token) {
			output = append(output, *token)
		})

		lexer.Start(tc.input, emitter, initState)

		assert.Equal(t, tc.output, output, "Failed lexer at %d: %s", i, output)
	}
}

func TestLexConstant(t *testing.T) {
	runTestCase(t, lex.Stmt, TestCases{
		{
			input: "",
			output: Tokens{
				{Type: token.EOF},
			},
		},
		{
			input: " ",
			output: Tokens{
				{Type: token.EOF, Start: 1, End: 1},
			},
		},
		{
			input: "rpc_version = 1",
			output: Tokens{
				{Type: token.Identifier, Start: 0, End: 11, Val: "rpc_version"},
				{Type: token.Assign, Val: "=", Start: 12, End: 13},
				{Type: token.Value, Start: 14, End: 15, Val: "1"},
				{Type: token.EOF, Start: 15, End: 15},
			},
		},
		{
			input: "rpc_version=1",
			output: Tokens{
				{Type: token.Identifier, Start: 0, End: 11, Val: "rpc_version"},
				{Type: token.Assign, Start: 11, End: 12, Val: "="},
				{Type: token.Value, Start: 12, End: 13, Val: "1"},
				{Type: token.EOF, Start: 13, End: 13, Val: ""},
			},
		},
		{
			input: "rpc_version= 1",
			output: Tokens{
				{Type: token.Identifier, Start: 0, End: 11, Val: "rpc_version"},
				{Type: token.Assign, Start: 11, End: 12, Val: "="},
				{Type: token.Value, Start: 13, End: 14, Val: "1"},
				{Type: token.EOF, Start: 14, End: 14, Val: ""},
			},
		},
	})
}

func TestLexEnum(t *testing.T) {
	runTestCase(t, lex.Stmt, TestCases{
		{
			input: `enum Kind int { }`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Type: token.Identifier, Start: 5, End: 9, Val: "Kind"},
				{Type: token.Type, Start: 10, End: 13, Val: "int"},
				{Type: token.OpenCurl, Start: 14, End: 15, Val: "{"},
				{Type: token.CloseCurl, Start: 16, End: 17, Val: "}"},
				{Type: token.EOF, Start: 17, End: 17, Val: ""},
			},
		},
		{
			input: `enum Kind int { A = 1 }`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Type: token.Identifier, Start: 5, End: 9, Val: "Kind"},
				{Type: token.Type, Start: 10, End: 13, Val: "int"},
				{Type: token.OpenCurl, Start: 14, End: 15, Val: "{"},
				{Type: token.Identifier, Start: 16, End: 17, Val: "A"},
				{Type: token.Assign, Start: 18, End: 19, Val: "="},
				{Type: token.Value, Start: 20, End: 21, Val: "1"},
				{Type: token.CloseCurl, Start: 22, End: 23, Val: "}"},
				{Type: token.EOF, Start: 23, End: 23, Val: ""},
			},
		},
		{
			input: `enum Kind int {
				A = 1
				B = 2
			}`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Type: token.Identifier, Start: 5, End: 9, Val: "Kind"},
				{Type: token.Type, Start: 10, End: 13, Val: "int"},
				{Type: token.OpenCurl, Start: 14, End: 15, Val: "{"},
				{Type: token.Identifier, Start: 20, End: 21, Val: "A"},
				{Type: token.Assign, Start: 22, End: 23, Val: "="},
				{Type: token.Value, Start: 24, End: 25, Val: "1"},
				{Type: token.Identifier, Start: 30, End: 31, Val: "B"},
				{Type: token.Assign, Start: 32, End: 33, Val: "="},
				{Type: token.Value, Start: 34, End: 35, Val: "2"},
				{Type: token.CloseCurl, Start: 39, End: 40, Val: "}"},
				{Type: token.EOF, Start: 40, End: 40, Val: ""},
			},
		},
		{
			input: `enum Kind int {
				A = 1
				B
				C
			}`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Type: token.Identifier, Start: 5, End: 9, Val: "Kind"},
				{Type: token.Type, Start: 10, End: 13, Val: "int"},
				{Type: token.OpenCurl, Start: 14, End: 15, Val: "{"},
				{Type: token.Identifier, Start: 20, End: 21, Val: "A"},
				{Type: token.Assign, Start: 22, End: 23, Val: "="},
				{Type: token.Value, Start: 24, End: 25, Val: "1"},
				{Type: token.Identifier, Start: 30, End: 31, Val: "B"},
				{Type: token.Identifier, Start: 36, End: 37, Val: "C"},
				{Type: token.CloseCurl, Start: 41, End: 42, Val: "}"},
				{Type: token.EOF, Start: 42, End: 42, Val: ""},
			},
		},
	})
}

func TestLexMessage(t *testing.T) {
	runTestCase(t, lex.Stmt, TestCases{
		{
			input: `message Foo {}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 11, Val: "Foo"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.CloseCurl, Start: 13, End: 14, Val: "}"},
				{Type: token.EOF, Start: 14, End: 14, Val: ""},
			},
		},
		{
			input: `message Foo {
				a: string
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 11, Val: "Foo"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Type: token.Colon, Start: 19, End: 20, Val: ":"},
				{Type: token.Type, Start: 21, End: 27, Val: "string"},
				{Type: token.CloseCurl, Start: 31, End: 32, Val: "}"},
				{Type: token.EOF, Start: 32, End: 32, Val: ""},
			},
		},
		{
			input: `message Foo {
				a: string
				b: int

				c: bool
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 11, Val: "Foo"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Type: token.Colon, Start: 19, End: 20, Val: ":"},
				{Type: token.Type, Start: 21, End: 27, Val: "string"},
				{Type: token.Identifier, Start: 32, End: 33, Val: "b"},
				{Type: token.Colon, Start: 33, End: 34, Val: ":"},
				{Type: token.Type, Start: 35, End: 38, Val: "int"},
				{Type: token.Identifier, Start: 44, End: 45, Val: "c"},
				{Type: token.Colon, Start: 45, End: 46, Val: ":"},
				{Type: token.Type, Start: 47, End: 51, Val: "bool"},
				{Type: token.CloseCurl, Start: 55, End: 56, Val: "}"},
				{Type: token.EOF, Start: 56, End: 56, Val: ""},
			},
		},
		{
			input: `message Foo {
				a?: string
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 11, Val: "Foo"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Type: token.Optional, Start: 19, End: 20, Val: "?"},
				{Type: token.Colon, Start: 20, End: 21, Val: ":"},
				{Type: token.Type, Start: 22, End: 28, Val: "string"},
				{Type: token.CloseCurl, Start: 32, End: 33, Val: "}"},
				{Type: token.EOF, Start: 33, End: 33, Val: ""},
			},
		},
		{
			input: `message Foo {
				a: []string
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 11, Val: "Foo"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Type: token.Colon, Start: 19, End: 20, Val: ":"},
				{Type: token.OpenBracket, Start: 21, End: 22, Val: "["},
				{Type: token.CloseBracket, Start: 22, End: 23, Val: "]"},
				{Type: token.Type, Start: 23, End: 29, Val: "string"},
				{Type: token.CloseCurl, Start: 33, End: 34, Val: "}"},
				{Type: token.EOF, Start: 34, End: 34, Val: ""},
			},
		},
		{
			input: `message Foo {
				a: []string
				b: int
				c?:uint32
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 11, Val: "Foo"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Type: token.Colon, Start: 19, End: 20, Val: ":"},
				{Type: token.OpenBracket, Start: 21, End: 22, Val: "["},
				{Type: token.CloseBracket, Start: 22, End: 23, Val: "]"},
				{Type: token.Type, Start: 23, End: 29, Val: "string"},
				{Type: token.Identifier, Start: 34, End: 35, Val: "b"},
				{Type: token.Colon, Start: 35, End: 36, Val: ":"},
				{Type: token.Type, Start: 37, End: 40, Val: "int"},
				{Type: token.Identifier, Start: 45, End: 46, Val: "c"},
				{Type: token.Optional, Start: 46, End: 47, Val: "?"},
				{Type: token.Colon, Start: 47, End: 48, Val: ":"},
				{Type: token.Type, Start: 48, End: 54, Val: "uint32"},
				{Type: token.CloseCurl, Start: 58, End: 59, Val: "}"},
				{Type: token.EOF, Start: 59, End: 59, Val: ""},
			},
		},
		{
			input: `message Foo {
				a: map<string,int>
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 11, Val: "Foo"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Type: token.Colon, Start: 19, End: 20, Val: ":"},
				{Type: token.Type, Start: 21, End: 24, Val: "map"},
				{Type: token.OpenAngle, Start: 24, End: 25, Val: "<"},
				{Type: token.Type, Start: 25, End: 31, Val: "string"},
				{Type: token.Comma, Start: 31, End: 32, Val: ","},
				{Type: token.Type, Start: 32, End: 35, Val: "int"},
				{Type: token.CloseAngle, Start: 35, End: 36, Val: ">"},
				{Type: token.CloseCurl, Start: 40, End: 41, Val: "}"},
				{Type: token.EOF, Start: 41, End: 41, Val: ""},
			},
		},
		{
			input: `
			message Foo {
				a?: string
			}
			
			message Bar {
				b: []Foo
			}
			`,
			output: Tokens{
				{Type: token.Message, Start: 4, End: 11, Val: "message"},
				{Type: token.Identifier, Start: 12, End: 15, Val: "Foo"},
				{Type: token.OpenCurl, Start: 16, End: 17, Val: "{"},
				{Type: token.Identifier, Start: 22, End: 23, Val: "a"},
				{Type: token.Optional, Start: 23, End: 24, Val: "?"},
				{Type: token.Colon, Start: 24, End: 25, Val: ":"},
				{Type: token.Type, Start: 26, End: 32, Val: "string"},
				{Type: token.CloseCurl, Start: 36, End: 37, Val: "}"},
				{Type: token.Message, Start: 45, End: 52, Val: "message"},
				{Type: token.Identifier, Start: 53, End: 56, Val: "Bar"},
				{Type: token.OpenCurl, Start: 57, End: 58, Val: "{"},
				{Type: token.Identifier, Start: 63, End: 64, Val: "b"},
				{Type: token.Colon, Start: 64, End: 65, Val: ":"},
				{Type: token.OpenBracket, Start: 66, End: 67, Val: "["},
				{Type: token.CloseBracket, Start: 67, End: 68, Val: "]"},
				{Type: token.Type, Start: 68, End: 71, Val: "Foo"},
				{Type: token.CloseCurl, Start: 75, End: 76, Val: "}"},
				{Type: token.EOF, Start: 80, End: 80, Val: ""},
			},
		},
		{
			input: `
			message Foo {
				a?: string
			}
			
			message Bar {
				...Foo
				b: int
			}
			`,
			output: Tokens{
				{Type: token.Message, Start: 4, End: 11, Val: "message"},
				{Type: token.Identifier, Start: 12, End: 15, Val: "Foo"},
				{Type: token.OpenCurl, Start: 16, End: 17, Val: "{"},
				{Type: token.Identifier, Start: 22, End: 23, Val: "a"},
				{Type: token.Optional, Start: 23, End: 24, Val: "?"},
				{Type: token.Colon, Start: 24, End: 25, Val: ":"},
				{Type: token.Type, Start: 26, End: 32, Val: "string"},
				{Type: token.CloseCurl, Start: 36, End: 37, Val: "}"},
				{Type: token.Message, Start: 45, End: 52, Val: "message"},
				{Type: token.Identifier, Start: 53, End: 56, Val: "Bar"},
				{Type: token.OpenCurl, Start: 57, End: 58, Val: "{"},
				{Type: token.Ellipsis, Start: 63, End: 66, Val: "..."},
				{Type: token.Type, Start: 66, End: 69, Val: "Foo"},
				{Type: token.Identifier, Start: 74, End: 75, Val: "b"},
				{Type: token.Colon, Start: 75, End: 76, Val: ":"},
				{Type: token.Type, Start: 77, End: 80, Val: "int"},
				{Type: token.CloseCurl, Start: 84, End: 85, Val: "}"},
				{Type: token.EOF, Start: 89, End: 89, Val: ""},
			},
		},
		{
			input: `message Foo {
				a: int {

				}
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 11, Val: "Foo"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Type: token.Colon, Start: 19, End: 20, Val: ":"},
				{Type: token.Type, Start: 21, End: 24, Val: "int"},
				{Type: token.OpenCurl, Start: 25, End: 26, Val: "{"},
				{Type: token.CloseCurl, Start: 32, End: 33, Val: "}"},
				{Type: token.CloseCurl, Start: 37, End: 38, Val: "}"},
				{Type: token.EOF, Start: 38, End: 38, Val: ""},
			},
		},
		{
			input: `message Foo {
				a: int {
					json = 
				}
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 11, Val: "Foo"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Type: token.Colon, Start: 19, End: 20, Val: ":"},
				{Type: token.Type, Start: 21, End: 24, Val: "int"},
				{Type: token.OpenCurl, Start: 25, End: 26, Val: "{"},
				{Type: token.Identifier, Start: 32, End: 36, Val: "json"},
				{Type: token.Assign, Start: 37, End: 38, Val: "="},
				{Type: token.Value, Start: 39, End: 39, Val: ""},
				{Type: token.CloseCurl, Start: 44, End: 45, Val: "}"},
				{Type: token.CloseCurl, Start: 49, End: 50, Val: "}"},
				{Type: token.EOF, Start: 50, End: 50, Val: ""},
			},
		},
		{
			input: `message Foo {
				a: int {
					json = 
					go.field.name = A
				}
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 11, Val: "Foo"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Type: token.Colon, Start: 19, End: 20, Val: ":"},
				{Type: token.Type, Start: 21, End: 24, Val: "int"},
				{Type: token.OpenCurl, Start: 25, End: 26, Val: "{"},
				{Type: token.Identifier, Start: 32, End: 36, Val: "json"},
				{Type: token.Assign, Start: 37, End: 38, Val: "="},
				{Type: token.Value, Start: 39, End: 39, Val: ""},
				{Type: token.Identifier, Start: 45, End: 58, Val: "go.field.name"},
				{Type: token.Assign, Start: 59, End: 60, Val: "="},
				{Type: token.Value, Start: 61, End: 62, Val: "A"},
				{Type: token.CloseCurl, Start: 67, End: 68, Val: "}"},
				{Type: token.CloseCurl, Start: 72, End: 73, Val: "}"},
				{Type: token.EOF, Start: 73, End: 73, Val: ""},
			},
		},
		{
			input: `message Foo {
				a: int {
					json = 
					go.field.name = A
				}
				b: int 
				c: float {
					json =
				}
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 11, Val: "Foo"},
				{Type: token.OpenCurl, Start: 12, End: 13, Val: "{"},
				{Type: token.Identifier, Start: 18, End: 19, Val: "a"},
				{Type: token.Colon, Start: 19, End: 20, Val: ":"},
				{Type: token.Type, Start: 21, End: 24, Val: "int"},
				{Type: token.OpenCurl, Start: 25, End: 26, Val: "{"},
				{Type: token.Identifier, Start: 32, End: 36, Val: "json"},
				{Type: token.Assign, Start: 37, End: 38, Val: "="},
				{Type: token.Value, Start: 39, End: 39, Val: ""},
				{Type: token.Identifier, Start: 45, End: 58, Val: "go.field.name"},
				{Type: token.Assign, Start: 59, End: 60, Val: "="},
				{Type: token.Value, Start: 61, End: 62, Val: "A"},
				{Type: token.CloseCurl, Start: 67, End: 68, Val: "}"},
				{Type: token.Identifier, Start: 73, End: 74, Val: "b"},
				{Type: token.Colon, Start: 74, End: 75, Val: ":"},
				{Type: token.Type, Start: 76, End: 79, Val: "int"},
				{Type: token.Identifier, Start: 85, End: 86, Val: "c"},
				{Type: token.Colon, Start: 86, End: 87, Val: ":"},
				{Type: token.Type, Start: 88, End: 93, Val: "float"},
				{Type: token.OpenCurl, Start: 94, End: 95, Val: "{"},
				{Type: token.Identifier, Start: 101, End: 105, Val: "json"},
				{Type: token.Assign, Start: 106, End: 107, Val: "="},
				{Type: token.Value, Start: 107, End: 107, Val: ""},
				{Type: token.CloseCurl, Start: 112, End: 113, Val: "}"},
				{Type: token.CloseCurl, Start: 117, End: 118, Val: "}"},
				{Type: token.EOF, Start: 118, End: 118, Val: ""},
			},
		},
	})
}
