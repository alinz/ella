package scanner_test

import (
	"testing"

	"compiler.ella.to/internal/scanner"
	"compiler.ella.to/internal/token"
)

func TestLex(t *testing.T) {
	runTestCase(t, -1, scanner.Lex, TestCases{
		{
			input: `service Foo {
				rpc GetFoo() => (value: int64) {
					Required
					A = 1mb
					B = 100h
				}
			}`,
			output: Tokens{
				{Type: token.Service, Start: 0, End: 7, Literal: "service"},
				{Type: token.Identifier, Start: 8, End: 11, Literal: "Foo"},
				{Type: token.OpenCurly, Start: 12, End: 13, Literal: "{"},
				{Type: token.Rpc, Start: 18, End: 21, Literal: "rpc"},
				{Type: token.Identifier, Start: 22, End: 28, Literal: "GetFoo"},
				{Type: token.OpenParen, Start: 28, End: 29, Literal: "("},
				{Type: token.CloseParen, Start: 29, End: 30, Literal: ")"},
				{Type: token.Return, Start: 31, End: 33, Literal: "=>"},
				{Type: token.OpenParen, Start: 34, End: 35, Literal: "("},
				{Type: token.Identifier, Start: 35, End: 40, Literal: "value"},
				{Type: token.Colon, Start: 40, End: 41, Literal: ":"},
				{Type: token.Int64, Start: 42, End: 47, Literal: "int64"},
				{Type: token.CloseParen, Start: 47, End: 48, Literal: ")"},
				{Type: token.OpenCurly, Start: 49, End: 50, Literal: "{"},
				{Type: token.Identifier, Start: 56, End: 64, Literal: "Required"},
				{Type: token.Identifier, Start: 70, End: 71, Literal: "A"},
				{Type: token.Assign, Start: 72, End: 73, Literal: "="},
				{Type: token.ConstBytes, Start: 74, End: 77, Literal: "1mb"},
				{Type: token.Identifier, Start: 83, End: 84, Literal: "B"},
				{Type: token.Assign, Start: 85, End: 86, Literal: "="},
				{Type: token.ConstDuration, Start: 87, End: 91, Literal: "100h"},
				{Type: token.CloseCurly, Start: 96, End: 97, Literal: "}"},
				{Type: token.CloseCurly, Start: 101, End: 102, Literal: "}"},
				{Type: token.EOF, Start: 102, End: 102, Literal: ""},
			},
		},
		{
			input: `A = 1mb`,
			output: Tokens{
				{Type: token.Identifier, Start: 0, End: 1, Literal: "A"},
				{Type: token.Assign, Start: 2, End: 3, Literal: "="},
				{Type: token.ConstBytes, Start: 4, End: 7, Literal: "1mb"},
				{Type: token.EOF, Start: 7, End: 7, Literal: ""},
			},
		},
		{
			skip: true,
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
				{Type: token.TopComment, Start: 9, End: 29, Literal: " this is a comment 1"},
				{Type: token.TopComment, Start: 34, End: 60, Literal: " this is another comment 2"},
				{Type: token.Identifier, Start: 64, End: 65, Literal: "a"},
				{Type: token.Assign, Start: 66, End: 67, Literal: "="},
				{Type: token.ConstInt, Start: 68, End: 69, Literal: "1"},
				{Type: token.RightComment, Start: 71, End: 91, Literal: " this is a comment 3"},
				{Type: token.TopComment, Start: 96, End: 122, Literal: " this is another comment 4"},
				{Type: token.Identifier, Start: 127, End: 134, Literal: "message"},
				{Type: token.Identifier, Start: 135, End: 136, Literal: "A"},
				{Type: token.OpenCurly, Start: 137, End: 138, Literal: "{"},
				{Type: token.TopComment, Start: 144, End: 164, Literal: " this is a comment 5"},
				{Type: token.TopComment, Start: 170, End: 196, Literal: " this is another comment 6"},
				{Type: token.Identifier, Start: 201, End: 210, Literal: "firstname"},
				{Type: token.Colon, Start: 210, End: 211, Literal: ":"},
				{Type: token.String, Start: 212, End: 218, Literal: "string"},
				{Type: token.CloseCurly, Start: 222, End: 223, Literal: "}"},
				{Type: token.EOF, Start: 231, End: 231, Literal: ""},
			},
		},
		{
			skip: true,
			input: `

			# This is a first comment
			a = 1 # this is the second comment
			# this is the third comment


			`,
			output: Tokens{
				{Type: token.TopComment, Start: 6, End: 30, Literal: " This is a first comment"},
				{Type: token.Identifier, Start: 34, End: 35, Literal: "a"},
				{Type: token.Assign, Start: 36, End: 37, Literal: "="},
				{Type: token.ConstInt, Start: 38, End: 39, Literal: "1"},
				{Type: token.RightComment, Start: 41, End: 68, Literal: " this is the second comment"},
				{Type: token.TopComment, Start: 73, End: 99, Literal: " this is the third comment"},
				{Type: token.EOF, Start: 105, End: 105, Literal: ""},
			},
		},
		{
			input: `ella = "1.0.0-b01"`,
			output: Tokens{
				{Type: token.Identifier, Start: 0, End: 4, Literal: "ella"},
				{Type: token.Assign, Start: 5, End: 6, Literal: "="},
				{Type: token.ConstStringDoubleQuote, Start: 8, End: 17, Literal: "1.0.0-b01"},
				{Type: token.EOF, Start: 18, End: 18, Literal: ""},
			},
		},
		{
			input: `message A {
				...B
				...C

				first: int64
			}`,
			output: Tokens{
				{Type: token.Identifier, Start: 0, End: 7, Literal: "message"},
				{Type: token.Identifier, Start: 8, End: 9, Literal: "A"},
				{Type: token.OpenCurly, Start: 10, End: 11, Literal: "{"},
				{Type: token.Extend, Start: 16, End: 19, Literal: "..."},
				{Type: token.Identifier, Start: 19, End: 20, Literal: "B"},
				{Type: token.Extend, Start: 25, End: 28, Literal: "..."},
				{Type: token.Identifier, Start: 28, End: 29, Literal: "C"},
				{Type: token.Identifier, Start: 35, End: 40, Literal: "first"},
				{Type: token.Colon, Start: 40, End: 41, Literal: ":"},
				{Type: token.Int64, Start: 42, End: 47, Literal: "int64"},
				{Type: token.CloseCurly, Start: 51, End: 52, Literal: "}"},
				{Type: token.EOF, Start: 52, End: 52, Literal: ""},
			},
		},
		{
			skip: true,
			input: `enum a int64 {
				one = 1 # comment
				two = 2# comment2
				three
			}`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Literal: "enum"},
				{Type: token.Identifier, Start: 5, End: 6, Literal: "a"},
				{Type: token.Int64, Start: 7, End: 12, Literal: "int64"},
				{Type: token.OpenCurly, Start: 13, End: 14, Literal: "{"},
				{Type: token.Identifier, Start: 19, End: 22, Literal: "one"},
				{Type: token.Assign, Start: 23, End: 24, Literal: "="},
				{Type: token.ConstInt, Start: 25, End: 26, Literal: "1"},
				{Type: token.RightComment, Start: 28, End: 36, Literal: " comment"},
				{Type: token.Identifier, Start: 41, End: 44, Literal: "two"},
				{Type: token.Assign, Start: 45, End: 46, Literal: "="},
				{Type: token.ConstInt, Start: 47, End: 48, Literal: "2"},
				{Type: token.RightComment, Start: 49, End: 58, Literal: " comment2"},
				{Type: token.Identifier, Start: 63, End: 68, Literal: "three"},
				{Type: token.CloseCurly, Start: 72, End: 73, Literal: "}"},
				{Type: token.EOF, Start: 73, End: 73, Literal: ""},
			},
		},
		{
			input: `enum a int64 {
				one = 1
				two = 2
				three
			}`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Literal: "enum"},
				{Type: token.Identifier, Start: 5, End: 6, Literal: "a"},
				{Type: token.Int64, Start: 7, End: 12, Literal: "int64"},
				{Type: token.OpenCurly, Start: 13, End: 14, Literal: "{"},
				{Type: token.Identifier, Start: 19, End: 22, Literal: "one"},
				{Type: token.Assign, Start: 23, End: 24, Literal: "="},
				{Type: token.ConstInt, Start: 25, End: 26, Literal: "1"},
				{Type: token.Identifier, Start: 31, End: 34, Literal: "two"},
				{Type: token.Assign, Start: 35, End: 36, Literal: "="},
				{Type: token.ConstInt, Start: 37, End: 38, Literal: "2"},
				{Type: token.Identifier, Start: 43, End: 48, Literal: "three"},
				{Type: token.CloseCurly, Start: 52, End: 53, Literal: "}"},
				{Type: token.EOF, Start: 53, End: 53, Literal: ""},
			},
		},
		{
			input: `enum a int64 {}`,
			output: Tokens{
				{Type: token.Enum, Start: 0, End: 4, Literal: "enum"},
				{Type: token.Identifier, Start: 5, End: 6, Literal: "a"},
				{Type: token.Int64, Start: 7, End: 12, Literal: "int64"},
				{Type: token.OpenCurly, Start: 13, End: 14, Literal: "{"},
				{Type: token.CloseCurly, Start: 14, End: 15, Literal: "}"},
				{Type: token.EOF, Start: 15, End: 15, Literal: ""},
			},
		},
		{
			input: `a=1`,
			output: Tokens{
				{Type: token.Identifier, Start: 0, End: 1, Literal: "a"},
				{Type: token.Assign, Start: 1, End: 2, Literal: "="},
				{Type: token.ConstInt, Start: 2, End: 3, Literal: "1"},
				{Type: token.EOF, Start: 3, End: 3, Literal: ""},
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

			service MyService {
				http GetUserById (id: int64) => (user: User) {
					method = "GET"
				}
			}
			
			`,
			output: Tokens{
				{Type: token.Identifier, Start: 8, End: 9, Literal: "a"},
				{Type: token.Assign, Start: 10, End: 11, Literal: "="},
				{Type: token.ConstFloat, Start: 12, End: 15, Literal: "1.0"},
				{Type: token.Identifier, Start: 20, End: 27, Literal: "message"},
				{Type: token.Identifier, Start: 28, End: 29, Literal: "A"},
				{Type: token.OpenCurly, Start: 30, End: 31, Literal: "{"},
				{Type: token.Identifier, Start: 36, End: 45, Literal: "firstname"},
				{Type: token.Colon, Start: 45, End: 46, Literal: ":"},
				{Type: token.String, Start: 47, End: 53, Literal: "string"},
				{Type: token.OpenCurly, Start: 54, End: 55, Literal: "{"},
				{Type: token.Identifier, Start: 61, End: 69, Literal: "required"},
				{Type: token.Identifier, Start: 75, End: 82, Literal: "pattern"},
				{Type: token.Assign, Start: 83, End: 84, Literal: "="},
				{Type: token.ConstStringDoubleQuote, Start: 86, End: 97, Literal: "^[a-zA-Z]+$"},
				{Type: token.CloseCurly, Start: 103, End: 104, Literal: "}"},
				{Type: token.CloseCurly, Start: 108, End: 109, Literal: "}"},
				{Type: token.Service, Start: 114, End: 121, Literal: "service"},
				{Type: token.Identifier, Start: 122, End: 131, Literal: "MyService"},
				{Type: token.OpenCurly, Start: 132, End: 133, Literal: "{"},
				{Type: token.Http, Start: 138, End: 142, Literal: "http"},
				{Type: token.Identifier, Start: 143, End: 154, Literal: "GetUserById"},
				{Type: token.OpenParen, Start: 155, End: 156, Literal: "("},
				{Type: token.Identifier, Start: 156, End: 158, Literal: "id"},
				{Type: token.Colon, Start: 158, End: 159, Literal: ":"},
				{Type: token.Int64, Start: 160, End: 165, Literal: "int64"},
				{Type: token.CloseParen, Start: 165, End: 166, Literal: ")"},
				{Type: token.Return, Start: 167, End: 169, Literal: "=>"},
				{Type: token.OpenParen, Start: 170, End: 171, Literal: "("},
				{Type: token.Identifier, Start: 171, End: 175, Literal: "user"},
				{Type: token.Colon, Start: 175, End: 176, Literal: ":"},
				{Type: token.Identifier, Start: 177, End: 181, Literal: "User"},
				{Type: token.CloseParen, Start: 181, End: 182, Literal: ")"},
				{Type: token.OpenCurly, Start: 183, End: 184, Literal: "{"},
				{Type: token.Identifier, Start: 190, End: 196, Literal: "method"},
				{Type: token.Assign, Start: 197, End: 198, Literal: "="},
				{Type: token.ConstStringDoubleQuote, Start: 200, End: 203, Literal: "GET"},
				{Type: token.CloseCurly, Start: 209, End: 210, Literal: "}"},
				{Type: token.CloseCurly, Start: 214, End: 215, Literal: "}"},
				{Type: token.EOF, Start: 223, End: 223, Literal: ""},
			},
		},
		{
			input: `error ErrUserNotFound { Code = 1000 HttpStatus = NotFound Msg = "user not found" }`,
			output: Tokens{
				{Type: token.CustomError, Start: 0, End: 5, Literal: "error"},
				{Type: token.Identifier, Start: 6, End: 21, Literal: "ErrUserNotFound"},
				{Type: token.OpenCurly, Start: 22, End: 23, Literal: "{"},
				{Type: token.Identifier, Start: 24, End: 28, Literal: "Code"},
				{Type: token.Assign, Start: 29, End: 30, Literal: "="},
				{Type: token.ConstInt, Start: 31, End: 35, Literal: "1000"},
				{Type: token.Identifier, Start: 36, End: 46, Literal: "HttpStatus"},
				{Type: token.Assign, Start: 47, End: 48, Literal: "="},
				{Type: token.Identifier, Start: 49, End: 57, Literal: "NotFound"},
				{Type: token.Identifier, Start: 58, End: 61, Literal: "Msg"},
				{Type: token.Assign, Start: 62, End: 63, Literal: "="},
				{Type: token.ConstStringDoubleQuote, Start: 65, End: 79, Literal: "user not found"},
				{Type: token.CloseCurly, Start: 81, End: 82, Literal: "}"},
				{Type: token.EOF, Start: 82, End: 82, Literal: ""},
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
					{Type: token.ConstInt, Start: 0, End: 1, Literal: "1"},
				},
			},
			{
				input: `1.0`,
				output: Tokens{
					{Type: token.ConstFloat, Start: 0, End: 3, Literal: "1.0"},
				},
			},
			{
				input: `1.`,
				output: Tokens{
					{Type: token.Error, Start: 0, End: 2, Literal: "expected digit after decimal point"},
				},
			},
			{
				input: `1.0.0`,
				output: Tokens{
					{Type: token.Error, Start: 0, End: 3, Literal: "unexpected character after number: ."},
				},
			},
			{
				input: `1_0_0`,
				output: Tokens{
					{Type: token.ConstInt, Start: 0, End: 5, Literal: "1_0_0"},
				},
			},
			{
				input:  `_1_0_0`,
				output: Tokens{},
			},
			{
				input: `1_0_0_`,
				output: Tokens{
					{Type: token.Error, Start: 0, End: 6, Literal: "expected digit after each underscore"},
				},
			},
			{
				input: `0.1_0_0`,
				output: Tokens{
					{Type: token.ConstFloat, Start: 0, End: 7, Literal: "0.1_0_0"},
				},
			},
			{
				input: `0.1__0_0`,
				output: Tokens{
					{Type: token.Error, Start: 0, End: 8, Literal: "expected digit after each underscore"},
				},
			},
			{
				input:  `hello`,
				output: Tokens{},
			},
			{
				input: `1_200kb`,
				output: Tokens{
					{Type: token.ConstBytes, Start: 0, End: 7, Literal: "1_200kb"},
				},
			},
		},
	)
}
