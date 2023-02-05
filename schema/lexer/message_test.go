package lexer_test

import (
	"testing"

	"github.com/alinz/rpc.go/schema/lexer"
	"github.com/alinz/rpc.go/schema/token"
)

func TestMessage(t *testing.T) {
	runTestCase(t, -1, lexer.Message(nil), TestCases{
		{
			input: `message ComplexType {
				meta: map<string,any>
				metaNestedExample: map<string,map<string,uint32>>
				namesList: []string
				numsList: []int64
				doubleArray: [][]string
				listOfMaps: []map<string,uint32>
				listOfUsers: []User
				mapOfUsers: map<string,User>
				user: User
			}`,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Identifier, Start: 8, End: 19, Val: "ComplexType"},
				{Kind: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Kind: token.Identifier, Start: 26, End: 30, Val: "meta"},
				{Kind: token.Colon, Start: 30, End: 31, Val: ":"},
				{Kind: token.Type, Start: 32, End: 35, Val: "map"},
				{Kind: token.OpenAngle, Start: 35, End: 36, Val: "<"},
				{Kind: token.Type, Start: 36, End: 42, Val: "string"},
				{Kind: token.Comma, Start: 42, End: 43, Val: ","},
				{Kind: token.Type, Start: 43, End: 46, Val: "any"},
				{Kind: token.CloseAngle, Start: 46, End: 47, Val: ">"},
				{Kind: token.Identifier, Start: 52, End: 69, Val: "metaNestedExample"},
				{Kind: token.Colon, Start: 69, End: 70, Val: ":"},
				{Kind: token.Type, Start: 71, End: 74, Val: "map"},
				{Kind: token.OpenAngle, Start: 74, End: 75, Val: "<"},
				{Kind: token.Type, Start: 75, End: 81, Val: "string"},
				{Kind: token.Comma, Start: 81, End: 82, Val: ","},
				{Kind: token.Type, Start: 82, End: 85, Val: "map"},
				{Kind: token.OpenAngle, Start: 85, End: 86, Val: "<"},
				{Kind: token.Type, Start: 86, End: 92, Val: "string"},
				{Kind: token.Comma, Start: 92, End: 93, Val: ","},
				{Kind: token.Type, Start: 93, End: 99, Val: "uint32"},
				{Kind: token.CloseAngle, Start: 99, End: 100, Val: ">"},
				{Kind: token.CloseAngle, Start: 100, End: 101, Val: ">"},
				{Kind: token.Identifier, Start: 106, End: 115, Val: "namesList"},
				{Kind: token.Colon, Start: 115, End: 116, Val: ":"},
				{Kind: token.OpenBracket, Start: 117, End: 118, Val: "["},
				{Kind: token.CloseBracket, Start: 118, End: 119, Val: "]"},
				{Kind: token.Type, Start: 119, End: 125, Val: "string"},
				{Kind: token.Identifier, Start: 130, End: 138, Val: "numsList"},
				{Kind: token.Colon, Start: 138, End: 139, Val: ":"},
				{Kind: token.OpenBracket, Start: 140, End: 141, Val: "["},
				{Kind: token.CloseBracket, Start: 141, End: 142, Val: "]"},
				{Kind: token.Type, Start: 142, End: 147, Val: "int64"},
				{Kind: token.Identifier, Start: 152, End: 163, Val: "doubleArray"},
				{Kind: token.Colon, Start: 163, End: 164, Val: ":"},
				{Kind: token.OpenBracket, Start: 165, End: 166, Val: "["},
				{Kind: token.CloseBracket, Start: 166, End: 167, Val: "]"},
				{Kind: token.OpenBracket, Start: 167, End: 168, Val: "["},
				{Kind: token.CloseBracket, Start: 168, End: 169, Val: "]"},
				{Kind: token.Type, Start: 169, End: 175, Val: "string"},
				{Kind: token.Identifier, Start: 180, End: 190, Val: "listOfMaps"},
				{Kind: token.Colon, Start: 190, End: 191, Val: ":"},
				{Kind: token.OpenBracket, Start: 192, End: 193, Val: "["},
				{Kind: token.CloseBracket, Start: 193, End: 194, Val: "]"},
				{Kind: token.Type, Start: 194, End: 197, Val: "map"},
				{Kind: token.OpenAngle, Start: 197, End: 198, Val: "<"},
				{Kind: token.Type, Start: 198, End: 204, Val: "string"},
				{Kind: token.Comma, Start: 204, End: 205, Val: ","},
				{Kind: token.Type, Start: 205, End: 211, Val: "uint32"},
				{Kind: token.CloseAngle, Start: 211, End: 212, Val: ">"},
				{Kind: token.Identifier, Start: 217, End: 228, Val: "listOfUsers"},
				{Kind: token.Colon, Start: 228, End: 229, Val: ":"},
				{Kind: token.OpenBracket, Start: 230, End: 231, Val: "["},
				{Kind: token.CloseBracket, Start: 231, End: 232, Val: "]"},
				{Kind: token.Type, Start: 232, End: 236, Val: "User"},
				{Kind: token.Identifier, Start: 241, End: 251, Val: "mapOfUsers"},
				{Kind: token.Colon, Start: 251, End: 252, Val: ":"},
				{Kind: token.Type, Start: 253, End: 256, Val: "map"},
				{Kind: token.OpenAngle, Start: 256, End: 257, Val: "<"},
				{Kind: token.Type, Start: 257, End: 263, Val: "string"},
				{Kind: token.Comma, Start: 263, End: 264, Val: ","},
				{Kind: token.Type, Start: 264, End: 268, Val: "User"},
				{Kind: token.CloseAngle, Start: 268, End: 269, Val: ">"},
				{Kind: token.Identifier, Start: 274, End: 278, Val: "user"},
				{Kind: token.Colon, Start: 278, End: 279, Val: ":"},
				{Kind: token.Type, Start: 280, End: 284, Val: "User"},
				{Kind: token.CloseCurl, Start: 288, End: 289, Val: "}"},
			},
		},
		{
			input: `message User {
				email: string {
					validation_regex = ^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$
				}
			}`,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.Identifier, Start: 19, End: 24, Val: "email"},
				{Kind: token.Colon, Start: 24, End: 25, Val: ":"},
				{Kind: token.Type, Start: 26, End: 32, Val: "string"},
				{Kind: token.OpenCurl, Start: 33, End: 34, Val: "{"},
				{Kind: token.Identifier, Start: 40, End: 56, Val: "validation_regex"},
				{Kind: token.Assign, Start: 57, End: 58, Val: "="},
				{Kind: token.Value, Start: 59, End: 101, Val: `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`},
				{Kind: token.CloseCurl, Start: 106, End: 107, Val: "}"},
				{Kind: token.CloseCurl, Start: 111, End: 112, Val: "}"},
			},
		},
		{
			input: `message User {
				...Base
				firstname: string
				lastname: string
				ops: map<string, string>
			}`,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.Ellipsis, Start: 19, End: 22, Val: "..."},
				{Kind: token.Type, Start: 22, End: 26, Val: "Base"},
				{Kind: token.Identifier, Start: 31, End: 40, Val: "firstname"},
				{Kind: token.Colon, Start: 40, End: 41, Val: ":"},
				{Kind: token.Type, Start: 42, End: 48, Val: "string"},
				{Kind: token.Identifier, Start: 53, End: 61, Val: "lastname"},
				{Kind: token.Colon, Start: 61, End: 62, Val: ":"},
				{Kind: token.Type, Start: 63, End: 69, Val: "string"},
				{Kind: token.Identifier, Start: 74, End: 77, Val: "ops"},
				{Kind: token.Colon, Start: 77, End: 78, Val: ":"},
				{Kind: token.Type, Start: 79, End: 82, Val: "map"},
				{Kind: token.OpenAngle, Start: 82, End: 83, Val: "<"},
				{Kind: token.Type, Start: 83, End: 89, Val: "string"},
				{Kind: token.Comma, Start: 89, End: 90, Val: ","},
				{Kind: token.Type, Start: 91, End: 97, Val: "string"},
				{Kind: token.CloseAngle, Start: 97, End: 98, Val: ">"},
				{Kind: token.CloseCurl, Start: 102, End: 103, Val: "}"},
			},
		},
		{
			input: `message User {
				...Base
			}`,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.Ellipsis, Start: 19, End: 22, Val: "..."},
				{Kind: token.Type, Start: 22, End: 26, Val: "Base"},
				{Kind: token.CloseCurl, Start: 30, End: 31, Val: "}"},
			},
		},
		{
			input: `message User { 
				users: map<int,string>
				ids: []string
			}`,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.Identifier, Start: 20, End: 25, Val: "users"},
				{Kind: token.Colon, Start: 25, End: 26, Val: ":"},
				{Kind: token.Type, Start: 27, End: 30, Val: "map"},
				{Kind: token.OpenAngle, Start: 30, End: 31, Val: "<"},
				{Kind: token.Type, Start: 31, End: 34, Val: "int"},
				{Kind: token.Comma, Start: 34, End: 35, Val: ","},
				{Kind: token.Type, Start: 35, End: 41, Val: "string"},
				{Kind: token.CloseAngle, Start: 41, End: 42, Val: ">"},
				{Kind: token.Identifier, Start: 47, End: 50, Val: "ids"},
				{Kind: token.Colon, Start: 50, End: 51, Val: ":"},
				{Kind: token.OpenBracket, Start: 52, End: 53, Val: "["},
				{Kind: token.CloseBracket, Start: 53, End: 54, Val: "]"},
				{Kind: token.Type, Start: 54, End: 60, Val: "string"},
				{Kind: token.CloseCurl, Start: 64, End: 65, Val: "}"},
			},
		},
		{
			input: `message User { 
				users: map<int,string>
			}`,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.Identifier, Start: 20, End: 25, Val: "users"},
				{Kind: token.Colon, Start: 25, End: 26, Val: ":"},
				{Kind: token.Type, Start: 27, End: 30, Val: "map"},
				{Kind: token.OpenAngle, Start: 30, End: 31, Val: "<"},
				{Kind: token.Type, Start: 31, End: 34, Val: "int"},
				{Kind: token.Comma, Start: 34, End: 35, Val: ","},
				{Kind: token.Type, Start: 35, End: 41, Val: "string"},
				{Kind: token.CloseAngle, Start: 41, End: 42, Val: ">"},
				{Kind: token.CloseCurl, Start: 46, End: 47, Val: "}"},
			},
		},
		{
			input: `message User { 
				id: []string
			}`,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.Identifier, Start: 20, End: 22, Val: "id"},
				{Kind: token.Colon, Start: 22, End: 23, Val: ":"},
				{Kind: token.OpenBracket, Start: 24, End: 25, Val: "["},
				{Kind: token.CloseBracket, Start: 25, End: 26, Val: "]"},
				{Kind: token.Type, Start: 26, End: 32, Val: "string"},
				{Kind: token.CloseCurl, Start: 36, End: 37, Val: "}"},
			},
		},
		{
			input: `message User { 
				id: string {
					json = id
				}
				name: string {

				}
			}`,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.Identifier, Start: 20, End: 22, Val: "id"},
				{Kind: token.Colon, Start: 22, End: 23, Val: ":"},
				{Kind: token.Type, Start: 24, End: 30, Val: "string"},
				{Kind: token.OpenCurl, Start: 31, End: 32, Val: "{"},
				{Kind: token.Identifier, Start: 38, End: 42, Val: "json"},
				{Kind: token.Assign, Start: 43, End: 44, Val: "="},
				{Kind: token.Value, Start: 45, End: 47, Val: "id"},
				{Kind: token.CloseCurl, Start: 52, End: 53, Val: "}"},
				{Kind: token.Identifier, Start: 58, End: 62, Val: "name"},
				{Kind: token.Colon, Start: 62, End: 63, Val: ":"},
				{Kind: token.Type, Start: 64, End: 70, Val: "string"},
				{Kind: token.OpenCurl, Start: 71, End: 72, Val: "{"},
				{Kind: token.CloseCurl, Start: 78, End: 79, Val: "}"},
				{Kind: token.CloseCurl, Start: 83, End: 84, Val: "}"},
			},
		},
		{
			input: `message User { 
				id: string {
					json = id
				}
				name: string
			}`,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.Identifier, Start: 20, End: 22, Val: "id"},
				{Kind: token.Colon, Start: 22, End: 23, Val: ":"},
				{Kind: token.Type, Start: 24, End: 30, Val: "string"},
				{Kind: token.OpenCurl, Start: 31, End: 32, Val: "{"},
				{Kind: token.Identifier, Start: 38, End: 42, Val: "json"},
				{Kind: token.Assign, Start: 43, End: 44, Val: "="},
				{Kind: token.Value, Start: 45, End: 47, Val: "id"},
				{Kind: token.CloseCurl, Start: 52, End: 53, Val: "}"},
				{Kind: token.Identifier, Start: 58, End: 62, Val: "name"},
				{Kind: token.Colon, Start: 62, End: 63, Val: ":"},
				{Kind: token.Type, Start: 64, End: 70, Val: "string"},
				{Kind: token.CloseCurl, Start: 74, End: 75, Val: "}"},
			},
		},
		{
			input: `message User { 
				id: string {
					json = id
				}
			}`,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.Identifier, Start: 20, End: 22, Val: "id"},
				{Kind: token.Colon, Start: 22, End: 23, Val: ":"},
				{Kind: token.Type, Start: 24, End: 30, Val: "string"},
				{Kind: token.OpenCurl, Start: 31, End: 32, Val: "{"},
				{Kind: token.Identifier, Start: 38, End: 42, Val: "json"},
				{Kind: token.Assign, Start: 43, End: 44, Val: "="},
				{Kind: token.Value, Start: 45, End: 47, Val: "id"},
				{Kind: token.CloseCurl, Start: 52, End: 53, Val: "}"},
				{Kind: token.CloseCurl, Start: 57, End: 58, Val: "}"},
			},
		},
		{
			input: `message User { 
				id: string {

				}
			}`,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.Identifier, Start: 20, End: 22, Val: "id"},
				{Kind: token.Colon, Start: 22, End: 23, Val: ":"},
				{Kind: token.Type, Start: 24, End: 30, Val: "string"},
				{Kind: token.OpenCurl, Start: 31, End: 32, Val: "{"},
				{Kind: token.CloseCurl, Start: 38, End: 39, Val: "}"},
				{Kind: token.CloseCurl, Start: 43, End: 44, Val: "}"},
			},
		},
		{
			input: `message User { 
				id: string
			}`,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.Identifier, Start: 20, End: 22, Val: "id"},
				{Kind: token.Colon, Start: 22, End: 23, Val: ":"},
				{Kind: token.Type, Start: 24, End: 30, Val: "string"},
				{Kind: token.CloseCurl, Start: 34, End: 35, Val: "}"},
			},
		},
		{
			input: `message User { }       `,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.CloseCurl, Start: 15, End: 16, Val: "}"},
			},
		},
		{
			input: `message User {}`,
			output: Tokens{
				{Kind: token.Message, Start: 0, End: 7, Val: "message"},
				{Kind: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Kind: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Kind: token.CloseCurl, Start: 14, End: 15, Val: "}"},
			},
		},
	})
}
