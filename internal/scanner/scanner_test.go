package scanner_test

import (
	"testing"

	"ella.to/internal/scanner"
	"ella.to/internal/token"
)

func TestLex(t *testing.T) {
	runTestCase(t, -1, scanner.Lex, TestCases{
		{
			input: `
			
			# this is a comment 1
			# this is another comment 2
			a = 1 # this is a comment 3
			# this is another comment 4

			message A {
				# this is a comment 5
				# this is another comment 6
				firstname: string
			}
			
			`,
			output: Tokens{
				{Type: token.TopComment, Start: 9, End: 29, Val: " this is a comment 1"},
				{Type: token.TopComment, Start: 34, End: 60, Val: " this is another comment 2"},
				{Type: token.Identifier, Start: 64, End: 65, Val: "a"},
				{Type: token.Assign, Start: 66, End: 67, Val: "="},
				{Type: token.ConstInt, Start: 68, End: 69, Val: "1"},
				{Type: token.RightComment, Start: 71, End: 91, Val: " this is a comment 3"},
				{Type: token.TopComment, Start: 96, End: 122, Val: " this is another comment 4"},
				{Type: token.Message, Start: 127, End: 134, Val: "message"},
				{Type: token.Identifier, Start: 135, End: 136, Val: "A"},
				{Type: token.OpenCurly, Start: 137, End: 138, Val: "{"},
				{Type: token.TopComment, Start: 144, End: 164, Val: " this is a comment 5"},
				{Type: token.TopComment, Start: 170, End: 196, Val: " this is another comment 6"},
				{Type: token.Identifier, Start: 201, End: 210, Val: "firstname"},
				{Type: token.Colon, Start: 210, End: 211, Val: ":"},
				{Type: token.String, Start: 212, End: 218, Val: "string"},
				{Type: token.CloseCurly, Start: 222, End: 223, Val: "}"},
				{Type: token.EOF, Start: 231, End: 231, Val: ""},
			},
		},
		{
			input: `

			# This is a first comment
			a = 1 # this is the second comment
			# this is the third comment


			`,
			output: Tokens{
				{Type: token.TopComment, Start: 6, End: 30, Val: " This is a first comment"},
				{Type: token.Identifier, Start: 34, End: 35, Val: "a"},
				{Type: token.Assign, Start: 36, End: 37, Val: "="},
				{Type: token.ConstInt, Start: 38, End: 39, Val: "1"},
				{Type: token.RightComment, Start: 41, End: 68, Val: " this is the second comment"},
				{Type: token.TopComment, Start: 73, End: 99, Val: " this is the third comment"},
				{Type: token.EOF, Start: 105, End: 105, Val: ""},
			},
		},
		{
			input: `ella = "1.0.0-b01"`,
			output: Tokens{
				{Type: token.Identifier, Start: 0, End: 4, Val: "ella"},
				{Type: token.Assign, Start: 5, End: 6, Val: "="},
				{Type: token.ConstStringDoubleQuote, Start: 8, End: 17, Val: "1.0.0-b01"},
				{Type: token.EOF, Start: 18, End: 18, Val: ""},
			},
		},
		{
			input: `message A {
				...B
				...C

				first: int64
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 9, Val: "A"},
				{Type: token.OpenCurly, Start: 10, End: 11, Val: "{"},
				{Type: token.Extend, Start: 16, End: 19, Val: "..."},
				{Type: token.Identifier, Start: 19, End: 20, Val: "B"},
				{Type: token.Extend, Start: 25, End: 28, Val: "..."},
				{Type: token.Identifier, Start: 28, End: 29, Val: "C"},
				{Type: token.Identifier, Start: 35, End: 40, Val: "first"},
				{Type: token.Colon, Start: 40, End: 41, Val: ":"},
				{Type: token.Int64, Start: 42, End: 47, Val: "int64"},
				{Type: token.CloseCurly, Start: 51, End: 52, Val: "}"},
				{Type: token.EOF, Start: 52, End: 52, Val: ""},
			},
		},
		{
			input: `enum a int64 {
				one = 1 # comment
				two = 2# comment2
				three
			}`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Type: token.Identifier, Start: 5, End: 6, Val: "a"},
				{Type: token.Int64, Start: 7, End: 12, Val: "int64"},
				{Type: token.OpenCurly, Start: 13, End: 14, Val: "{"},
				{Type: token.Identifier, Start: 19, End: 22, Val: "one"},
				{Type: token.Assign, Start: 23, End: 24, Val: "="},
				{Type: token.ConstInt, Start: 25, End: 26, Val: "1"},
				{Type: token.RightComment, Start: 28, End: 36, Val: " comment"},
				{Type: token.Identifier, Start: 41, End: 44, Val: "two"},
				{Type: token.Assign, Start: 45, End: 46, Val: "="},
				{Type: token.ConstInt, Start: 47, End: 48, Val: "2"},
				{Type: token.RightComment, Start: 49, End: 58, Val: " comment2"},
				{Type: token.Identifier, Start: 63, End: 68, Val: "three"},
				{Type: token.CloseCurly, Start: 72, End: 73, Val: "}"},
				{Type: token.EOF, Start: 73, End: 73, Val: ""},
			},
		},
		{
			input: `enum a int64 {
				one = 1
				two = 2
				three
			}`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Type: token.Identifier, Start: 5, End: 6, Val: "a"},
				{Type: token.Int64, Start: 7, End: 12, Val: "int64"},
				{Type: token.OpenCurly, Start: 13, End: 14, Val: "{"},
				{Type: token.Identifier, Start: 19, End: 22, Val: "one"},
				{Type: token.Assign, Start: 23, End: 24, Val: "="},
				{Type: token.ConstInt, Start: 25, End: 26, Val: "1"},
				{Type: token.Identifier, Start: 31, End: 34, Val: "two"},
				{Type: token.Assign, Start: 35, End: 36, Val: "="},
				{Type: token.ConstInt, Start: 37, End: 38, Val: "2"},
				{Type: token.Identifier, Start: 43, End: 48, Val: "three"},
				{Type: token.CloseCurly, Start: 52, End: 53, Val: "}"},
				{Type: token.EOF, Start: 53, End: 53, Val: ""},
			},
		},
		{
			input: `enum a int64 {}`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Val: "enum"},
				{Type: token.Identifier, Start: 5, End: 6, Val: "a"},
				{Type: token.Int64, Start: 7, End: 12, Val: "int64"},
				{Type: token.OpenCurly, Start: 13, End: 14, Val: "{"},
				{Type: token.CloseCurly, Start: 14, End: 15, Val: "}"},
				{Type: token.EOF, Start: 15, End: 15, Val: ""},
			},
		},
		{
			input: `a=1`,
			output: Tokens{
				{Type: token.Identifier, Start: 0, End: 1, Val: "a"},
				{Type: token.Assign, Start: 1, End: 2, Val: "="},
				{Type: token.ConstInt, Start: 2, End: 3, Val: "1"},
				{Type: token.EOF, Start: 3, End: 3, Val: ""},
			},
		},
		{
			input: `
			
			a = 1.0

			message A {
				firstname: string {
					required
					pattern = "^[a-zA-Z]+$"
				}
			}

			http MyService {
				GetUserById (id: int64) => (user: User) {
					method = "GET"
				}
			}
			
			`,
			output: Tokens{
				{Type: token.Identifier, Start: 8, End: 9, Val: "a"},
				{Type: token.Assign, Start: 10, End: 11, Val: "="},
				{Type: token.ConstFloat, Start: 12, End: 15, Val: "1.0"},
				{Type: token.Message, Start: 20, End: 27, Val: "message"},
				{Type: token.Identifier, Start: 28, End: 29, Val: "A"},
				{Type: token.OpenCurly, Start: 30, End: 31, Val: "{"},
				{Type: token.Identifier, Start: 36, End: 45, Val: "firstname"},
				{Type: token.Colon, Start: 45, End: 46, Val: ":"},
				{Type: token.String, Start: 47, End: 53, Val: "string"},
				{Type: token.OpenCurly, Start: 54, End: 55, Val: "{"},
				{Type: token.Identifier, Start: 61, End: 69, Val: "required"},
				{Type: token.Identifier, Start: 75, End: 82, Val: "pattern"},
				{Type: token.Assign, Start: 83, End: 84, Val: "="},
				{Type: token.ConstStringDoubleQuote, Start: 86, End: 97, Val: "^[a-zA-Z]+$"},
				{Type: token.CloseCurly, Start: 103, End: 104, Val: "}"},
				{Type: token.CloseCurly, Start: 108, End: 109, Val: "}"},
				{Type: token.Http, Start: 114, End: 118, Val: "http"},
				{Type: token.Identifier, Start: 119, End: 128, Val: "MyService"},
				{Type: token.OpenCurly, Start: 129, End: 130, Val: "{"},
				{Type: token.Identifier, Start: 135, End: 146, Val: "GetUserById"},
				{Type: token.OpenParen, Start: 147, End: 148, Val: "("},
				{Type: token.Identifier, Start: 148, End: 150, Val: "id"},
				{Type: token.Colon, Start: 150, End: 151, Val: ":"},
				{Type: token.Int64, Start: 152, End: 157, Val: "int64"},
				{Type: token.CloseParen, Start: 157, End: 158, Val: ")"},
				{Type: token.Return, Start: 159, End: 161, Val: "=>"},
				{Type: token.OpenParen, Start: 162, End: 163, Val: "("},
				{Type: token.Identifier, Start: 163, End: 167, Val: "user"},
				{Type: token.Colon, Start: 167, End: 168, Val: ":"},
				{Type: token.Identifier, Start: 169, End: 173, Val: "User"},
				{Type: token.CloseParen, Start: 173, End: 174, Val: ")"},
				{Type: token.OpenCurly, Start: 175, End: 176, Val: "{"},
				{Type: token.Identifier, Start: 182, End: 188, Val: "method"},
				{Type: token.Assign, Start: 189, End: 190, Val: "="},
				{Type: token.ConstStringDoubleQuote, Start: 192, End: 195, Val: "GET"},
				{Type: token.CloseCurly, Start: 201, End: 202, Val: "}"},
				{Type: token.CloseCurly, Start: 206, End: 207, Val: "}"},
				{Type: token.EOF, Start: 215, End: 215, Val: ""},
			},
		},
	})
}

func TestNumber(t *testing.T) {

	runTestCase(t, -1, scanner.Number,
		TestCases{
			{
				input: `1`,
				output: Tokens{
					{Type: token.ConstInt, Start: 0, End: 1, Val: "1"},
				},
			},
			{
				input: `1.0`,
				output: Tokens{
					{Type: token.ConstFloat, Start: 0, End: 3, Val: "1.0"},
				},
			},
			{
				input: `1.`,
				output: Tokens{
					{Type: token.Error, Start: 0, End: 2, Val: "expected digit after decimal point"},
				},
			},
			{
				input: `1.0.0`,
				output: Tokens{
					{Type: token.Error, Start: 0, End: 3, Val: "unexpected character after number: ."},
				},
			},
			{
				input: `1_0_0`,
				output: Tokens{
					{Type: token.ConstInt, Start: 0, End: 5, Val: "1_0_0"},
				},
			},
			{
				input:  `_1_0_0`,
				output: Tokens{},
			},
			{
				input: `1_0_0_`,
				output: Tokens{
					{Type: token.Error, Start: 0, End: 6, Val: "expected digit after each underscore"},
				},
			},
			{
				input: `0.1_0_0`,
				output: Tokens{
					{Type: token.ConstFloat, Start: 0, End: 7, Val: "0.1_0_0"},
				},
			},
			{
				input: `0.1__0_0`,
				output: Tokens{
					{Type: token.Error, Start: 0, End: 8, Val: "expected digit after each underscore"},
				},
			},
			{
				input:  `hello`,
				output: Tokens{},
			},
		},
	)
}
