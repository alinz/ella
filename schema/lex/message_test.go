package lex_test

import (
	"testing"

	"github.com/alinz/rpc.go/schema/lex"
	"github.com/alinz/rpc.go/schema/lex/token"
)

func TestMessage(t *testing.T) {
	runTestCase(t, -1, lex.Message(nil), TestCases{
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
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 19, Val: "ComplexType"},
				{Type: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Type: token.Identifier, Start: 26, End: 30, Val: "meta"},
				{Type: token.Colon, Start: 30, End: 31, Val: ":"},
				{Type: token.Type, Start: 32, End: 35, Val: "map"},
				{Type: token.OpenAngle, Start: 35, End: 36, Val: "<"},
				{Type: token.Type, Start: 36, End: 42, Val: "string"},
				{Type: token.Comma, Start: 42, End: 43, Val: ","},
				{Type: token.Type, Start: 43, End: 46, Val: "any"},
				{Type: token.CloseAngle, Start: 46, End: 47, Val: ">"},
				{Type: token.Identifier, Start: 52, End: 69, Val: "metaNestedExample"},
				{Type: token.Colon, Start: 69, End: 70, Val: ":"},
				{Type: token.Type, Start: 71, End: 74, Val: "map"},
				{Type: token.OpenAngle, Start: 74, End: 75, Val: "<"},
				{Type: token.Type, Start: 75, End: 81, Val: "string"},
				{Type: token.Comma, Start: 81, End: 82, Val: ","},
				{Type: token.Type, Start: 82, End: 85, Val: "map"},
				{Type: token.OpenAngle, Start: 85, End: 86, Val: "<"},
				{Type: token.Type, Start: 86, End: 92, Val: "string"},
				{Type: token.Comma, Start: 92, End: 93, Val: ","},
				{Type: token.Type, Start: 93, End: 99, Val: "uint32"},
				{Type: token.CloseAngle, Start: 99, End: 100, Val: ">"},
				{Type: token.CloseAngle, Start: 100, End: 101, Val: ">"},
				{Type: token.Identifier, Start: 106, End: 115, Val: "namesList"},
				{Type: token.Colon, Start: 115, End: 116, Val: ":"},
				{Type: token.OpenBracket, Start: 117, End: 118, Val: "["},
				{Type: token.CloseBracket, Start: 118, End: 119, Val: "]"},
				{Type: token.Type, Start: 119, End: 125, Val: "string"},
				{Type: token.Identifier, Start: 130, End: 138, Val: "numsList"},
				{Type: token.Colon, Start: 138, End: 139, Val: ":"},
				{Type: token.OpenBracket, Start: 140, End: 141, Val: "["},
				{Type: token.CloseBracket, Start: 141, End: 142, Val: "]"},
				{Type: token.Type, Start: 142, End: 147, Val: "int64"},
				{Type: token.Identifier, Start: 152, End: 163, Val: "doubleArray"},
				{Type: token.Colon, Start: 163, End: 164, Val: ":"},
				{Type: token.OpenBracket, Start: 165, End: 166, Val: "["},
				{Type: token.CloseBracket, Start: 166, End: 167, Val: "]"},
				{Type: token.OpenBracket, Start: 167, End: 168, Val: "["},
				{Type: token.CloseBracket, Start: 168, End: 169, Val: "]"},
				{Type: token.Type, Start: 169, End: 175, Val: "string"},
				{Type: token.Identifier, Start: 180, End: 190, Val: "listOfMaps"},
				{Type: token.Colon, Start: 190, End: 191, Val: ":"},
				{Type: token.OpenBracket, Start: 192, End: 193, Val: "["},
				{Type: token.CloseBracket, Start: 193, End: 194, Val: "]"},
				{Type: token.Type, Start: 194, End: 197, Val: "map"},
				{Type: token.OpenAngle, Start: 197, End: 198, Val: "<"},
				{Type: token.Type, Start: 198, End: 204, Val: "string"},
				{Type: token.Comma, Start: 204, End: 205, Val: ","},
				{Type: token.Type, Start: 205, End: 211, Val: "uint32"},
				{Type: token.CloseAngle, Start: 211, End: 212, Val: ">"},
				{Type: token.Identifier, Start: 217, End: 228, Val: "listOfUsers"},
				{Type: token.Colon, Start: 228, End: 229, Val: ":"},
				{Type: token.OpenBracket, Start: 230, End: 231, Val: "["},
				{Type: token.CloseBracket, Start: 231, End: 232, Val: "]"},
				{Type: token.Type, Start: 232, End: 236, Val: "User"},
				{Type: token.Identifier, Start: 241, End: 251, Val: "mapOfUsers"},
				{Type: token.Colon, Start: 251, End: 252, Val: ":"},
				{Type: token.Type, Start: 253, End: 256, Val: "map"},
				{Type: token.OpenAngle, Start: 256, End: 257, Val: "<"},
				{Type: token.Type, Start: 257, End: 263, Val: "string"},
				{Type: token.Comma, Start: 263, End: 264, Val: ","},
				{Type: token.Type, Start: 264, End: 268, Val: "User"},
				{Type: token.CloseAngle, Start: 268, End: 269, Val: ">"},
				{Type: token.Identifier, Start: 274, End: 278, Val: "user"},
				{Type: token.Colon, Start: 278, End: 279, Val: ":"},
				{Type: token.Type, Start: 280, End: 284, Val: "User"},
				{Type: token.CloseCurl, Start: 288, End: 289, Val: "}"},
			},
		},
		{
			input: `message User {
				email: string {
					validation_regex = ^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$
				}
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Type: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Type: token.Identifier, Start: 19, End: 24, Val: "email"},
				{Type: token.Colon, Start: 24, End: 25, Val: ":"},
				{Type: token.Type, Start: 26, End: 32, Val: "string"},
				{Type: token.OpenCurl, Start: 33, End: 34, Val: "{"},
				{Type: token.Identifier, Start: 40, End: 56, Val: "validation_regex"},
				{Type: token.Assign, Start: 57, End: 58, Val: "="},
				{Type: token.Value, Start: 59, End: 101, Val: `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`},
				{Type: token.CloseCurl, Start: 106, End: 107, Val: "}"},
				{Type: token.CloseCurl, Start: 111, End: 112, Val: "}"},
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
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Type: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Type: token.Ellipsis, Start: 19, End: 22, Val: "..."},
				{Type: token.Type, Start: 22, End: 26, Val: "Base"},
				{Type: token.Identifier, Start: 31, End: 40, Val: "firstname"},
				{Type: token.Colon, Start: 40, End: 41, Val: ":"},
				{Type: token.Type, Start: 42, End: 48, Val: "string"},
				{Type: token.Identifier, Start: 53, End: 61, Val: "lastname"},
				{Type: token.Colon, Start: 61, End: 62, Val: ":"},
				{Type: token.Type, Start: 63, End: 69, Val: "string"},
				{Type: token.Identifier, Start: 74, End: 77, Val: "ops"},
				{Type: token.Colon, Start: 77, End: 78, Val: ":"},
				{Type: token.Type, Start: 79, End: 82, Val: "map"},
				{Type: token.OpenAngle, Start: 82, End: 83, Val: "<"},
				{Type: token.Type, Start: 83, End: 89, Val: "string"},
				{Type: token.Comma, Start: 89, End: 90, Val: ","},
				{Type: token.Type, Start: 91, End: 97, Val: "string"},
				{Type: token.CloseAngle, Start: 97, End: 98, Val: ">"},
				{Type: token.CloseCurl, Start: 102, End: 103, Val: "}"},
			},
		},
		{
			input: `message User {
				...Base
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Type: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Type: token.Ellipsis, Start: 19, End: 22, Val: "..."},
				{Type: token.Type, Start: 22, End: 26, Val: "Base"},
				{Type: token.CloseCurl, Start: 30, End: 31, Val: "}"},
			},
		},
		{
			input: `message User { 
				users: map<int,string>
				ids: []string
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Type: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Type: token.Identifier, Start: 20, End: 25, Val: "users"},
				{Type: token.Colon, Start: 25, End: 26, Val: ":"},
				{Type: token.Type, Start: 27, End: 30, Val: "map"},
				{Type: token.OpenAngle, Start: 30, End: 31, Val: "<"},
				{Type: token.Type, Start: 31, End: 34, Val: "int"},
				{Type: token.Comma, Start: 34, End: 35, Val: ","},
				{Type: token.Type, Start: 35, End: 41, Val: "string"},
				{Type: token.CloseAngle, Start: 41, End: 42, Val: ">"},
				{Type: token.Identifier, Start: 47, End: 50, Val: "ids"},
				{Type: token.Colon, Start: 50, End: 51, Val: ":"},
				{Type: token.OpenBracket, Start: 52, End: 53, Val: "["},
				{Type: token.CloseBracket, Start: 53, End: 54, Val: "]"},
				{Type: token.Type, Start: 54, End: 60, Val: "string"},
				{Type: token.CloseCurl, Start: 64, End: 65, Val: "}"},
			},
		},
		{
			input: `message User { 
				users: map<int,string>
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Type: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Type: token.Identifier, Start: 20, End: 25, Val: "users"},
				{Type: token.Colon, Start: 25, End: 26, Val: ":"},
				{Type: token.Type, Start: 27, End: 30, Val: "map"},
				{Type: token.OpenAngle, Start: 30, End: 31, Val: "<"},
				{Type: token.Type, Start: 31, End: 34, Val: "int"},
				{Type: token.Comma, Start: 34, End: 35, Val: ","},
				{Type: token.Type, Start: 35, End: 41, Val: "string"},
				{Type: token.CloseAngle, Start: 41, End: 42, Val: ">"},
				{Type: token.CloseCurl, Start: 46, End: 47, Val: "}"},
			},
		},
		{
			input: `message User { 
				id: []string
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Type: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Type: token.Identifier, Start: 20, End: 22, Val: "id"},
				{Type: token.Colon, Start: 22, End: 23, Val: ":"},
				{Type: token.OpenBracket, Start: 24, End: 25, Val: "["},
				{Type: token.CloseBracket, Start: 25, End: 26, Val: "]"},
				{Type: token.Type, Start: 26, End: 32, Val: "string"},
				{Type: token.CloseCurl, Start: 36, End: 37, Val: "}"},
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
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Type: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Type: token.Identifier, Start: 20, End: 22, Val: "id"},
				{Type: token.Colon, Start: 22, End: 23, Val: ":"},
				{Type: token.Type, Start: 24, End: 30, Val: "string"},
				{Type: token.OpenCurl, Start: 31, End: 32, Val: "{"},
				{Type: token.Identifier, Start: 38, End: 42, Val: "json"},
				{Type: token.Assign, Start: 43, End: 44, Val: "="},
				{Type: token.Value, Start: 45, End: 47, Val: "id"},
				{Type: token.CloseCurl, Start: 52, End: 53, Val: "}"},
				{Type: token.Identifier, Start: 58, End: 62, Val: "name"},
				{Type: token.Colon, Start: 62, End: 63, Val: ":"},
				{Type: token.Type, Start: 64, End: 70, Val: "string"},
				{Type: token.OpenCurl, Start: 71, End: 72, Val: "{"},
				{Type: token.CloseCurl, Start: 78, End: 79, Val: "}"},
				{Type: token.CloseCurl, Start: 83, End: 84, Val: "}"},
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
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Type: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Type: token.Identifier, Start: 20, End: 22, Val: "id"},
				{Type: token.Colon, Start: 22, End: 23, Val: ":"},
				{Type: token.Type, Start: 24, End: 30, Val: "string"},
				{Type: token.OpenCurl, Start: 31, End: 32, Val: "{"},
				{Type: token.Identifier, Start: 38, End: 42, Val: "json"},
				{Type: token.Assign, Start: 43, End: 44, Val: "="},
				{Type: token.Value, Start: 45, End: 47, Val: "id"},
				{Type: token.CloseCurl, Start: 52, End: 53, Val: "}"},
				{Type: token.Identifier, Start: 58, End: 62, Val: "name"},
				{Type: token.Colon, Start: 62, End: 63, Val: ":"},
				{Type: token.Type, Start: 64, End: 70, Val: "string"},
				{Type: token.CloseCurl, Start: 74, End: 75, Val: "}"},
			},
		},
		{
			input: `message User { 
				id: string {
					json = id
				}
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Type: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Type: token.Identifier, Start: 20, End: 22, Val: "id"},
				{Type: token.Colon, Start: 22, End: 23, Val: ":"},
				{Type: token.Type, Start: 24, End: 30, Val: "string"},
				{Type: token.OpenCurl, Start: 31, End: 32, Val: "{"},
				{Type: token.Identifier, Start: 38, End: 42, Val: "json"},
				{Type: token.Assign, Start: 43, End: 44, Val: "="},
				{Type: token.Value, Start: 45, End: 47, Val: "id"},
				{Type: token.CloseCurl, Start: 52, End: 53, Val: "}"},
				{Type: token.CloseCurl, Start: 57, End: 58, Val: "}"},
			},
		},
		{
			input: `message User { 
				id: string {

				}
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Type: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Type: token.Identifier, Start: 20, End: 22, Val: "id"},
				{Type: token.Colon, Start: 22, End: 23, Val: ":"},
				{Type: token.Type, Start: 24, End: 30, Val: "string"},
				{Type: token.OpenCurl, Start: 31, End: 32, Val: "{"},
				{Type: token.CloseCurl, Start: 38, End: 39, Val: "}"},
				{Type: token.CloseCurl, Start: 43, End: 44, Val: "}"},
			},
		},
		{
			input: `message User { 
				id: string
			}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Type: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Type: token.Identifier, Start: 20, End: 22, Val: "id"},
				{Type: token.Colon, Start: 22, End: 23, Val: ":"},
				{Type: token.Type, Start: 24, End: 30, Val: "string"},
				{Type: token.CloseCurl, Start: 34, End: 35, Val: "}"},
			},
		},
		{
			input: `message User { }       `,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Type: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Type: token.CloseCurl, Start: 15, End: 16, Val: "}"},
			},
		},
		{
			input: `message User {}`,
			output: Tokens{
				{Type: token.Message, Start: 0, End: 7, Val: "message"},
				{Type: token.Identifier, Start: 8, End: 12, Val: "User"},
				{Type: token.OpenCurl, Start: 13, End: 14, Val: "{"},
				{Type: token.CloseCurl, Start: 14, End: 15, Val: "}"},
			},
		},
	})
}
