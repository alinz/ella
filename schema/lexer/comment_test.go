package lexer_test

import (
	"testing"

	"github.com/alinz/rpc.go/schema/lexer"
	"github.com/alinz/rpc.go/schema/token"
)

func TestLexComment(t *testing.T) {
	runTestCase(t, -1, lexer.Stmt(nil), TestCases{
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
				{Kind: token.Comment, Start: 4, End: 5, Val: "#"},
				{Kind: token.Value, Start: 5, End: 27, Val: " this is first comment"},
				{Kind: token.Message, Start: 31, End: 38, Val: "message"},
				{Kind: token.Identifier, Start: 39, End: 50, Val: "ComplexType"},
				{Kind: token.OpenCurl, Start: 51, End: 52, Val: "{"},
				{Kind: token.Comment, Start: 57, End: 58, Val: "#"},
				{Kind: token.Value, Start: 58, End: 81, Val: " this is second comment"},
				{Kind: token.Identifier, Start: 86, End: 90, Val: "meta"},
				{Kind: token.Colon, Start: 90, End: 91, Val: ":"},
				{Kind: token.Type, Start: 92, End: 95, Val: "map"},
				{Kind: token.OpenAngle, Start: 95, End: 96, Val: "<"},
				{Kind: token.Type, Start: 96, End: 102, Val: "string"},
				{Kind: token.Comma, Start: 102, End: 103, Val: ","},
				{Kind: token.Type, Start: 103, End: 106, Val: "any"},
				{Kind: token.CloseAngle, Start: 106, End: 107, Val: ">"},
				{Kind: token.Type, Start: 108, End: 109, Val: "#"},
				{Kind: token.Type, Start: 110, End: 114, Val: "this"},
				{Kind: token.Type, Start: 115, End: 117, Val: "is"},
				{Kind: token.Type, Start: 118, End: 123, Val: "third"},
				{Kind: token.Type, Start: 124, End: 131, Val: "comment"},
				{Kind: token.Identifier, Start: 136, End: 153, Val: "metaNestedExample"},
				{Kind: token.Colon, Start: 153, End: 154, Val: ":"},
				{Kind: token.Type, Start: 155, End: 158, Val: "map"},
				{Kind: token.OpenAngle, Start: 158, End: 159, Val: "<"},
				{Kind: token.Type, Start: 159, End: 165, Val: "string"},
				{Kind: token.Comma, Start: 165, End: 166, Val: ","},
				{Kind: token.Type, Start: 166, End: 169, Val: "map"},
				{Kind: token.OpenAngle, Start: 169, End: 170, Val: "<"},
				{Kind: token.Type, Start: 170, End: 176, Val: "string"},
				{Kind: token.Comma, Start: 176, End: 177, Val: ","},
				{Kind: token.Type, Start: 177, End: 183, Val: "uint32"},
				{Kind: token.CloseAngle, Start: 183, End: 184, Val: ">"},
				{Kind: token.CloseAngle, Start: 184, End: 185, Val: ">"},
				{Kind: token.Identifier, Start: 190, End: 199, Val: "namesList"},
				{Kind: token.Colon, Start: 199, End: 200, Val: ":"},
				{Kind: token.OpenBracket, Start: 201, End: 202, Val: "["},
				{Kind: token.CloseBracket, Start: 202, End: 203, Val: "]"},
				{Kind: token.Type, Start: 203, End: 209, Val: "string"},
				{Kind: token.Identifier, Start: 214, End: 222, Val: "numsList"},
				{Kind: token.Colon, Start: 222, End: 223, Val: ":"},
				{Kind: token.OpenBracket, Start: 224, End: 225, Val: "["},
				{Kind: token.CloseBracket, Start: 225, End: 226, Val: "]"},
				{Kind: token.Type, Start: 226, End: 231, Val: "int64"},
				{Kind: token.Identifier, Start: 236, End: 247, Val: "doubleArray"},
				{Kind: token.Colon, Start: 247, End: 248, Val: ":"},
				{Kind: token.OpenBracket, Start: 249, End: 250, Val: "["},
				{Kind: token.CloseBracket, Start: 250, End: 251, Val: "]"},
				{Kind: token.OpenBracket, Start: 251, End: 252, Val: "["},
				{Kind: token.CloseBracket, Start: 252, End: 253, Val: "]"},
				{Kind: token.Type, Start: 253, End: 259, Val: "string"},
				{Kind: token.Type, Start: 260, End: 261, Val: "#"},
				{Kind: token.Type, Start: 262, End: 266, Val: "this"},
				{Kind: token.Type, Start: 267, End: 269, Val: "is"},
				{Kind: token.Type, Start: 270, End: 276, Val: "fourth"},
				{Kind: token.Type, Start: 277, End: 284, Val: "comment"},
				{Kind: token.Identifier, Start: 289, End: 299, Val: "listOfMaps"},
				{Kind: token.Colon, Start: 299, End: 300, Val: ":"},
				{Kind: token.OpenBracket, Start: 301, End: 302, Val: "["},
				{Kind: token.CloseBracket, Start: 302, End: 303, Val: "]"},
				{Kind: token.Type, Start: 303, End: 306, Val: "map"},
				{Kind: token.OpenAngle, Start: 306, End: 307, Val: "<"},
				{Kind: token.Type, Start: 307, End: 313, Val: "string"},
				{Kind: token.Comma, Start: 313, End: 314, Val: ","},
				{Kind: token.Type, Start: 314, End: 320, Val: "uint32"},
				{Kind: token.CloseAngle, Start: 320, End: 321, Val: ">"},
				{Kind: token.Identifier, Start: 326, End: 337, Val: "listOfUsers"},
				{Kind: token.Colon, Start: 337, End: 338, Val: ":"},
				{Kind: token.OpenBracket, Start: 339, End: 340, Val: "["},
				{Kind: token.CloseBracket, Start: 340, End: 341, Val: "]"},
				{Kind: token.Type, Start: 341, End: 345, Val: "User"},
				{Kind: token.Identifier, Start: 350, End: 360, Val: "mapOfUsers"},
				{Kind: token.Colon, Start: 360, End: 361, Val: ":"},
				{Kind: token.Type, Start: 362, End: 365, Val: "map"},
				{Kind: token.OpenAngle, Start: 365, End: 366, Val: "<"},
				{Kind: token.Type, Start: 366, End: 372, Val: "string"},
				{Kind: token.Comma, Start: 372, End: 373, Val: ","},
				{Kind: token.Type, Start: 373, End: 377, Val: "User"},
				{Kind: token.CloseAngle, Start: 377, End: 378, Val: ">"},
				{Kind: token.Identifier, Start: 383, End: 387, Val: "user"},
				{Kind: token.Colon, Start: 387, End: 388, Val: ":"},
				{Kind: token.Type, Start: 389, End: 393, Val: "User"},
				{Kind: token.Comment, Start: 398, End: 399, Val: "#"},
				{Kind: token.Value, Start: 399, End: 421, Val: " this is fifth comment"},
				{Kind: token.CloseCurl, Start: 425, End: 426, Val: "}"},
				{Kind: token.Comment, Start: 427, End: 428, Val: "#"},
				{Kind: token.Value, Start: 428, End: 450, Val: " this is sixth comment"},
				{Kind: token.Comment, Start: 454, End: 455, Val: "#"},
				{Kind: token.Value, Start: 455, End: 479, Val: " this is seventh comment"},
				{Kind: token.EOF, Start: 483, End: 483, Val: ""},
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
				{Kind: token.Comment, Start: 4, End: 5, Val: "#"},
				{Kind: token.Value, Start: 5, End: 27, Val: " this is first comment"},
				{Kind: token.Service, Start: 31, End: 38, Val: "service"},
				{Kind: token.Identifier, Start: 39, End: 50, Val: "TestService"},
				{Kind: token.OpenCurl, Start: 51, End: 52, Val: "{"},
				{Kind: token.Comment, Start: 57, End: 58, Val: "#"},
				{Kind: token.Value, Start: 58, End: 81, Val: " this is second comment"},
				{Kind: token.Identifier, Start: 86, End: 90, Val: "Ping"},
				{Kind: token.OpenParen, Start: 90, End: 91, Val: "("},
				{Kind: token.CloseParen, Start: 91, End: 92, Val: ")"},
				{Kind: token.Return, Start: 93, End: 95, Val: "=>"},
				{Kind: token.OpenParen, Start: 96, End: 97, Val: "("},
				{Kind: token.Identifier, Start: 97, End: 103, Val: "status"},
				{Kind: token.Colon, Start: 103, End: 104, Val: ":"},
				{Kind: token.Stream, Start: 105, End: 111, Val: "stream"},
				{Kind: token.Type, Start: 112, End: 115, Val: "int"},
				{Kind: token.CloseParen, Start: 115, End: 116, Val: ")"},
				{Kind: token.Comment, Start: 117, End: 118, Val: "#"},
				{Kind: token.Value, Start: 118, End: 140, Val: " this is third comment"},
				{Kind: token.Comment, Start: 145, End: 146, Val: "#"},
				{Kind: token.Value, Start: 146, End: 169, Val: " this is fourth comment"},
				{Kind: token.CloseCurl, Start: 173, End: 174, Val: "}"},
				{Kind: token.Comment, Start: 175, End: 176, Val: "#"},
				{Kind: token.Value, Start: 176, End: 198, Val: " this is fifth comment"},
				{Kind: token.Comment, Start: 202, End: 203, Val: "#"},
				{Kind: token.Value, Start: 203, End: 225, Val: " this is sixth comment"},
				{Kind: token.EOF, Start: 229, End: 229, Val: ""},
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
				{Kind: token.Comment, Start: 4, End: 5, Val: "#"},
				{Kind: token.Value, Start: 5, End: 27, Val: " this is first comment"},
				{Kind: token.Enum, Start: 31, End: 35, Val: "enum"},
				{Kind: token.Identifier, Start: 36, End: 39, Val: "Foo"},
				{Kind: token.Type, Start: 40, End: 45, Val: "int32"},
				{Kind: token.OpenCurl, Start: 46, End: 47, Val: "{"},
				{Kind: token.Comment, Start: 52, End: 53, Val: "#"},
				{Kind: token.Value, Start: 53, End: 76, Val: " this is second comment"},
				{Kind: token.Identifier, Start: 81, End: 82, Val: "a"},
				{Kind: token.Assign, Start: 83, End: 84, Val: "="},
				{Kind: token.Value, Start: 85, End: 86, Val: "1"},
				{Kind: token.CloseCurl, Start: 90, End: 91, Val: "}"},
				{Kind: token.Comment, Start: 92, End: 93, Val: "#"},
				{Kind: token.Value, Start: 93, End: 115, Val: " this is third comment"},
				{Kind: token.Comment, Start: 119, End: 120, Val: "#"},
				{Kind: token.Value, Start: 120, End: 143, Val: " this is fourth comment"},
				{Kind: token.EOF, Start: 147, End: 147, Val: ""},
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
				{Kind: token.Comment, Start: 4, End: 5, Val: "#"},
				{Kind: token.Value, Start: 5, End: 27, Val: " this is first comment"},
				{Kind: token.Enum, Start: 31, End: 35, Val: "enum"},
				{Kind: token.Identifier, Start: 36, End: 39, Val: "Foo"},
				{Kind: token.Type, Start: 40, End: 45, Val: "int32"},
				{Kind: token.OpenCurl, Start: 46, End: 47, Val: "{"},
				{Kind: token.Comment, Start: 52, End: 53, Val: "#"},
				{Kind: token.Value, Start: 53, End: 76, Val: " this is second comment"},
				{Kind: token.CloseCurl, Start: 80, End: 81, Val: "}"},
				{Kind: token.Comment, Start: 82, End: 83, Val: "#"},
				{Kind: token.Value, Start: 83, End: 105, Val: " this is third comment"},
				{Kind: token.Comment, Start: 109, End: 110, Val: "#"},
				{Kind: token.Value, Start: 110, End: 133, Val: " this is fourth comment"},
				{Kind: token.EOF, Start: 137, End: 137, Val: ""},
			},
		},
		{
			input: `
			# this is first comment
			enum Foo int32 {

			}
			`,
			output: Tokens{
				{Kind: token.Comment, Start: 4, End: 5, Val: "#"},
				{Kind: token.Value, Start: 5, End: 27, Val: " this is first comment"},
				{Kind: token.Enum, Start: 31, End: 35, Val: "enum"},
				{Kind: token.Identifier, Start: 36, End: 39, Val: "Foo"},
				{Kind: token.Type, Start: 40, End: 45, Val: "int32"},
				{Kind: token.OpenCurl, Start: 46, End: 47, Val: "{"},
				{Kind: token.CloseCurl, Start: 52, End: 53, Val: "}"},
				{Kind: token.EOF, Start: 57, End: 57, Val: ""},
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
				{Kind: token.Comment, Start: 5, End: 6, Val: "#"},
				{Kind: token.Value, Start: 6, End: 24, Val: " this is a comment"},
				{Kind: token.Identifier, Start: 29, End: 32, Val: "foo"},
				{Kind: token.Assign, Start: 33, End: 34, Val: "="},
				{Kind: token.Value, Start: 35, End: 36, Val: "1"},
				{Kind: token.Comment, Start: 37, End: 38, Val: "#"},
				{Kind: token.Value, Start: 38, End: 61, Val: " this is second comment"},
				{Kind: token.Comment, Start: 66, End: 67, Val: "#"},
				{Kind: token.Value, Start: 67, End: 89, Val: " this is third comment"},
				{Kind: token.Comment, Start: 94, End: 95, Val: "#"},
				{Kind: token.Value, Start: 95, End: 118, Val: " this is fourth comment"},
				{Kind: token.EOF, Start: 122, End: 122, Val: ""},
			},
		},
		{
			input: `
				# this is a comment
				foo = 1 # this is second comment
			`,
			output: Tokens{
				{Kind: token.Comment, Start: 5, End: 6, Val: "#"},
				{Kind: token.Value, Start: 6, End: 24, Val: " this is a comment"},
				{Kind: token.Identifier, Start: 29, End: 32, Val: "foo"},
				{Kind: token.Assign, Start: 33, End: 34, Val: "="},
				{Kind: token.Value, Start: 35, End: 36, Val: "1"},
				{Kind: token.Comment, Start: 37, End: 38, Val: "#"},
				{Kind: token.Value, Start: 38, End: 61, Val: " this is second comment"},
				{Kind: token.EOF, Start: 65, End: 65, Val: ""},
			},
		},
		{
			input: `
				# this is a comment
				foo = 1
			`,
			output: Tokens{
				{Kind: token.Comment, Start: 5, End: 6, Val: "#"},
				{Kind: token.Value, Start: 6, End: 24, Val: " this is a comment"},
				{Kind: token.Identifier, Start: 29, End: 32, Val: "foo"},
				{Kind: token.Assign, Start: 33, End: 34, Val: "="},
				{Kind: token.Value, Start: 35, End: 36, Val: "1"},
				{Kind: token.EOF, Start: 40, End: 40, Val: ""},
			},
		},
		{
			input: `
				# this is a comment
			`,
			output: Tokens{
				{Kind: token.Comment, Start: 5, End: 6, Val: "#"},
				{Kind: token.Value, Start: 6, End: 24, Val: " this is a comment"},
				{Kind: token.EOF, Start: 28, End: 28, Val: ""},
			},
		},
	})
}
