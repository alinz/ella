package lex_test

import (
	"testing"

	"github.com/alinz/rpc.go/schema/lex"
	"github.com/alinz/rpc.go/schema/lex/token"
)

func TestLexComment(t *testing.T) {
	runTestCase(t, -1, lex.Stmt(nil), TestCases{
		{
			input: `
			# this is first comment
			message ComplexType {
				# this is second comment
				meta: map<string,any> # this is third comment
				metaNestedExample: map<string,map<string,uint32>>
				namesList: []string
				numsList: []int64
				doubleArray: [][]string # this is fourth comment
				listOfMaps: []map<string,uint32>
				listOfUsers: []User
				mapOfUsers: map<string,User>
				user: User
				# this is fifth comment
			} # this is sixth comment
			# this is seventh comment
			`,
			output: Tokens{
				{Type: token.Comment, Start: 4, End: 5, Val: "#"},
				{Type: token.Value, Start: 5, End: 27, Val: " this is first comment"},
				{Type: token.Message, Start: 31, End: 38, Val: "message"},
				{Type: token.Identifier, Start: 39, End: 50, Val: "ComplexType"},
				{Type: token.OpenCurl, Start: 51, End: 52, Val: "{"},
				{Type: token.Comment, Start: 57, End: 58, Val: "#"},
				{Type: token.Value, Start: 58, End: 81, Val: " this is second comment"},
				{Type: token.Identifier, Start: 86, End: 90, Val: "meta"},
				{Type: token.Colon, Start: 90, End: 91, Val: ":"},
				{Type: token.Type, Start: 92, End: 95, Val: "map"},
				{Type: token.OpenAngle, Start: 95, End: 96, Val: "<"},
				{Type: token.Type, Start: 96, End: 102, Val: "string"},
				{Type: token.Comma, Start: 102, End: 103, Val: ","},
				{Type: token.Type, Start: 103, End: 106, Val: "any"},
				{Type: token.CloseAngle, Start: 106, End: 107, Val: ">"},
				{Type: token.Type, Start: 108, End: 109, Val: "#"},
				{Type: token.Type, Start: 110, End: 114, Val: "this"},
				{Type: token.Type, Start: 115, End: 117, Val: "is"},
				{Type: token.Type, Start: 118, End: 123, Val: "third"},
				{Type: token.Type, Start: 124, End: 131, Val: "comment"},
				{Type: token.Identifier, Start: 136, End: 153, Val: "metaNestedExample"},
				{Type: token.Colon, Start: 153, End: 154, Val: ":"},
				{Type: token.Type, Start: 155, End: 158, Val: "map"},
				{Type: token.OpenAngle, Start: 158, End: 159, Val: "<"},
				{Type: token.Type, Start: 159, End: 165, Val: "string"},
				{Type: token.Comma, Start: 165, End: 166, Val: ","},
				{Type: token.Type, Start: 166, End: 169, Val: "map"},
				{Type: token.OpenAngle, Start: 169, End: 170, Val: "<"},
				{Type: token.Type, Start: 170, End: 176, Val: "string"},
				{Type: token.Comma, Start: 176, End: 177, Val: ","},
				{Type: token.Type, Start: 177, End: 183, Val: "uint32"},
				{Type: token.CloseAngle, Start: 183, End: 184, Val: ">"},
				{Type: token.CloseAngle, Start: 184, End: 185, Val: ">"},
				{Type: token.Identifier, Start: 190, End: 199, Val: "namesList"},
				{Type: token.Colon, Start: 199, End: 200, Val: ":"},
				{Type: token.OpenBracket, Start: 201, End: 202, Val: "["},
				{Type: token.CloseBracket, Start: 202, End: 203, Val: "]"},
				{Type: token.Type, Start: 203, End: 209, Val: "string"},
				{Type: token.Identifier, Start: 214, End: 222, Val: "numsList"},
				{Type: token.Colon, Start: 222, End: 223, Val: ":"},
				{Type: token.OpenBracket, Start: 224, End: 225, Val: "["},
				{Type: token.CloseBracket, Start: 225, End: 226, Val: "]"},
				{Type: token.Type, Start: 226, End: 231, Val: "int64"},
				{Type: token.Identifier, Start: 236, End: 247, Val: "doubleArray"},
				{Type: token.Colon, Start: 247, End: 248, Val: ":"},
				{Type: token.OpenBracket, Start: 249, End: 250, Val: "["},
				{Type: token.CloseBracket, Start: 250, End: 251, Val: "]"},
				{Type: token.OpenBracket, Start: 251, End: 252, Val: "["},
				{Type: token.CloseBracket, Start: 252, End: 253, Val: "]"},
				{Type: token.Type, Start: 253, End: 259, Val: "string"},
				{Type: token.Type, Start: 260, End: 261, Val: "#"},
				{Type: token.Type, Start: 262, End: 266, Val: "this"},
				{Type: token.Type, Start: 267, End: 269, Val: "is"},
				{Type: token.Type, Start: 270, End: 276, Val: "fourth"},
				{Type: token.Type, Start: 277, End: 284, Val: "comment"},
				{Type: token.Identifier, Start: 289, End: 299, Val: "listOfMaps"},
				{Type: token.Colon, Start: 299, End: 300, Val: ":"},
				{Type: token.OpenBracket, Start: 301, End: 302, Val: "["},
				{Type: token.CloseBracket, Start: 302, End: 303, Val: "]"},
				{Type: token.Type, Start: 303, End: 306, Val: "map"},
				{Type: token.OpenAngle, Start: 306, End: 307, Val: "<"},
				{Type: token.Type, Start: 307, End: 313, Val: "string"},
				{Type: token.Comma, Start: 313, End: 314, Val: ","},
				{Type: token.Type, Start: 314, End: 320, Val: "uint32"},
				{Type: token.CloseAngle, Start: 320, End: 321, Val: ">"},
				{Type: token.Identifier, Start: 326, End: 337, Val: "listOfUsers"},
				{Type: token.Colon, Start: 337, End: 338, Val: ":"},
				{Type: token.OpenBracket, Start: 339, End: 340, Val: "["},
				{Type: token.CloseBracket, Start: 340, End: 341, Val: "]"},
				{Type: token.Type, Start: 341, End: 345, Val: "User"},
				{Type: token.Identifier, Start: 350, End: 360, Val: "mapOfUsers"},
				{Type: token.Colon, Start: 360, End: 361, Val: ":"},
				{Type: token.Type, Start: 362, End: 365, Val: "map"},
				{Type: token.OpenAngle, Start: 365, End: 366, Val: "<"},
				{Type: token.Type, Start: 366, End: 372, Val: "string"},
				{Type: token.Comma, Start: 372, End: 373, Val: ","},
				{Type: token.Type, Start: 373, End: 377, Val: "User"},
				{Type: token.CloseAngle, Start: 377, End: 378, Val: ">"},
				{Type: token.Identifier, Start: 383, End: 387, Val: "user"},
				{Type: token.Colon, Start: 387, End: 388, Val: ":"},
				{Type: token.Type, Start: 389, End: 393, Val: "User"},
				{Type: token.Comment, Start: 398, End: 399, Val: "#"},
				{Type: token.Value, Start: 399, End: 421, Val: " this is fifth comment"},
				{Type: token.CloseCurl, Start: 425, End: 426, Val: "}"},
				{Type: token.Comment, Start: 427, End: 428, Val: "#"},
				{Type: token.Value, Start: 428, End: 450, Val: " this is sixth comment"},
				{Type: token.Comment, Start: 454, End: 455, Val: "#"},
				{Type: token.Value, Start: 455, End: 479, Val: " this is seventh comment"},
				{Type: token.EOF, Start: 483, End: 483, Val: ""},
			},
		},
		{
			input: `
			# this is first comment
			service TestService {
				# this is second comment
				Ping() => (status: stream int) # this is third comment
				# this is fourth comment
			} # this is fifth comment
			# this is sixth comment
			`,
			output: Tokens{
				{Type: token.Comment, Start: 4, End: 5, Val: "#"},
				{Type: token.Value, Start: 5, End: 27, Val: " this is first comment"},
				{Type: token.Service, Start: 31, End: 38, Val: "service"},
				{Type: token.Identifier, Start: 39, End: 50, Val: "TestService"},
				{Type: token.OpenCurl, Start: 51, End: 52, Val: "{"},
				{Type: token.Comment, Start: 57, End: 58, Val: "#"},
				{Type: token.Value, Start: 58, End: 81, Val: " this is second comment"},
				{Type: token.Identifier, Start: 86, End: 90, Val: "Ping"},
				{Type: token.OpenParen, Start: 90, End: 91, Val: "("},
				{Type: token.CloseParen, Start: 91, End: 92, Val: ")"},
				{Type: token.Return, Start: 93, End: 95, Val: "=>"},
				{Type: token.OpenParen, Start: 96, End: 97, Val: "("},
				{Type: token.Identifier, Start: 97, End: 103, Val: "status"},
				{Type: token.Colon, Start: 103, End: 104, Val: ":"},
				{Type: token.Stream, Start: 105, End: 111, Val: "stream"},
				{Type: token.Type, Start: 112, End: 115, Val: "int"},
				{Type: token.CloseParen, Start: 115, End: 116, Val: ")"},
				{Type: token.Comment, Start: 117, End: 118, Val: "#"},
				{Type: token.Value, Start: 118, End: 140, Val: " this is third comment"},
				{Type: token.Comment, Start: 145, End: 146, Val: "#"},
				{Type: token.Value, Start: 146, End: 169, Val: " this is fourth comment"},
				{Type: token.CloseCurl, Start: 173, End: 174, Val: "}"},
				{Type: token.Comment, Start: 175, End: 176, Val: "#"},
				{Type: token.Value, Start: 176, End: 198, Val: " this is fifth comment"},
				{Type: token.Comment, Start: 202, End: 203, Val: "#"},
				{Type: token.Value, Start: 203, End: 225, Val: " this is sixth comment"},
				{Type: token.EOF, Start: 229, End: 229, Val: ""},
			},
		},
		{
			input: `
			# this is first comment
			enum Foo int32 {
				# this is second comment
				a = 1
			} # this is third comment
			# this is fourth comment
			`,
			output: Tokens{
				{Type: token.Comment, Start: 4, End: 5, Val: "#"},
				{Type: token.Value, Start: 5, End: 27, Val: " this is first comment"},
				{Type: token.Enum, Start: 31, End: 35, Val: "enum"},
				{Type: token.Identifier, Start: 36, End: 39, Val: "Foo"},
				{Type: token.Type, Start: 40, End: 45, Val: "int32"},
				{Type: token.OpenCurl, Start: 46, End: 47, Val: "{"},
				{Type: token.Comment, Start: 52, End: 53, Val: "#"},
				{Type: token.Value, Start: 53, End: 76, Val: " this is second comment"},
				{Type: token.Identifier, Start: 81, End: 82, Val: "a"},
				{Type: token.Assign, Start: 83, End: 84, Val: "="},
				{Type: token.Value, Start: 85, End: 86, Val: "1"},
				{Type: token.CloseCurl, Start: 90, End: 91, Val: "}"},
				{Type: token.Comment, Start: 92, End: 93, Val: "#"},
				{Type: token.Value, Start: 93, End: 115, Val: " this is third comment"},
				{Type: token.Comment, Start: 119, End: 120, Val: "#"},
				{Type: token.Value, Start: 120, End: 143, Val: " this is fourth comment"},
				{Type: token.EOF, Start: 147, End: 147, Val: ""},
			},
		},
		{
			input: `
			# this is first comment
			enum Foo int32 {
				# this is second comment
			} # this is third comment
			# this is fourth comment
			`,
			output: Tokens{
				{Type: token.Comment, Start: 4, End: 5, Val: "#"},
				{Type: token.Value, Start: 5, End: 27, Val: " this is first comment"},
				{Type: token.Enum, Start: 31, End: 35, Val: "enum"},
				{Type: token.Identifier, Start: 36, End: 39, Val: "Foo"},
				{Type: token.Type, Start: 40, End: 45, Val: "int32"},
				{Type: token.OpenCurl, Start: 46, End: 47, Val: "{"},
				{Type: token.Comment, Start: 52, End: 53, Val: "#"},
				{Type: token.Value, Start: 53, End: 76, Val: " this is second comment"},
				{Type: token.CloseCurl, Start: 80, End: 81, Val: "}"},
				{Type: token.Comment, Start: 82, End: 83, Val: "#"},
				{Type: token.Value, Start: 83, End: 105, Val: " this is third comment"},
				{Type: token.Comment, Start: 109, End: 110, Val: "#"},
				{Type: token.Value, Start: 110, End: 133, Val: " this is fourth comment"},
				{Type: token.EOF, Start: 137, End: 137, Val: ""},
			},
		},
		{
			input: `
			# this is first comment
			enum Foo int32 {

			}
			`,
			output: Tokens{
				{Type: token.Comment, Start: 4, End: 5, Val: "#"},
				{Type: token.Value, Start: 5, End: 27, Val: " this is first comment"},
				{Type: token.Enum, Start: 31, End: 35, Val: "enum"},
				{Type: token.Identifier, Start: 36, End: 39, Val: "Foo"},
				{Type: token.Type, Start: 40, End: 45, Val: "int32"},
				{Type: token.OpenCurl, Start: 46, End: 47, Val: "{"},
				{Type: token.CloseCurl, Start: 52, End: 53, Val: "}"},
				{Type: token.EOF, Start: 57, End: 57, Val: ""},
			},
		},
		{
			input: `
				# this is a comment
				foo = 1 # this is second comment
				# this is third comment
				# this is fourth comment
			`,
			output: Tokens{
				{Type: token.Comment, Start: 5, End: 6, Val: "#"},
				{Type: token.Value, Start: 6, End: 24, Val: " this is a comment"},
				{Type: token.Identifier, Start: 29, End: 32, Val: "foo"},
				{Type: token.Assign, Start: 33, End: 34, Val: "="},
				{Type: token.Value, Start: 35, End: 36, Val: "1"},
				{Type: token.Comment, Start: 37, End: 38, Val: "#"},
				{Type: token.Value, Start: 38, End: 61, Val: " this is second comment"},
				{Type: token.Comment, Start: 66, End: 67, Val: "#"},
				{Type: token.Value, Start: 67, End: 89, Val: " this is third comment"},
				{Type: token.Comment, Start: 94, End: 95, Val: "#"},
				{Type: token.Value, Start: 95, End: 118, Val: " this is fourth comment"},
				{Type: token.EOF, Start: 122, End: 122, Val: ""},
			},
		},
		{
			input: `
				# this is a comment
				foo = 1 # this is second comment
			`,
			output: Tokens{
				{Type: token.Comment, Start: 5, End: 6, Val: "#"},
				{Type: token.Value, Start: 6, End: 24, Val: " this is a comment"},
				{Type: token.Identifier, Start: 29, End: 32, Val: "foo"},
				{Type: token.Assign, Start: 33, End: 34, Val: "="},
				{Type: token.Value, Start: 35, End: 36, Val: "1"},
				{Type: token.Comment, Start: 37, End: 38, Val: "#"},
				{Type: token.Value, Start: 38, End: 61, Val: " this is second comment"},
				{Type: token.EOF, Start: 65, End: 65, Val: ""},
			},
		},
		{
			input: `
				# this is a comment
				foo = 1
			`,
			output: Tokens{
				{Type: token.Comment, Start: 5, End: 6, Val: "#"},
				{Type: token.Value, Start: 6, End: 24, Val: " this is a comment"},
				{Type: token.Identifier, Start: 29, End: 32, Val: "foo"},
				{Type: token.Assign, Start: 33, End: 34, Val: "="},
				{Type: token.Value, Start: 35, End: 36, Val: "1"},
				{Type: token.EOF, Start: 40, End: 40, Val: ""},
			},
		},
		{
			input: `
				# this is a comment
			`,
			output: Tokens{
				{Type: token.Comment, Start: 5, End: 6, Val: "#"},
				{Type: token.Value, Start: 6, End: 24, Val: " this is a comment"},
				{Type: token.EOF, Start: 28, End: 28, Val: ""},
			},
		},
	})
}
