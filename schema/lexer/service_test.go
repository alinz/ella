package lexer_test

import (
	"testing"

	"github.com/alinz/rpc.go/schema/lexer"
	"github.com/alinz/rpc.go/schema/token"
)

func TestLexService(t *testing.T) {
	runTestCase(t, -1, lexer.Service(nil), TestCases{
		{
			input: `service TestService {
				Ping() {
					Method = GET
				}
			}`,
			output: Tokens{
				{Kind: token.Service, Start: 0, End: 7, Val: "service"},
				{Kind: token.Identifier, Start: 8, End: 19, Val: "TestService"},
				{Kind: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Kind: token.Identifier, Start: 26, End: 30, Val: "Ping"},
				{Kind: token.OpenParen, Start: 30, End: 31, Val: "("},
				{Kind: token.CloseParen, Start: 31, End: 32, Val: ")"},
				{Kind: token.OpenCurl, Start: 33, End: 34, Val: "{"},
				{Kind: token.Identifier, Start: 40, End: 46, Val: "Method"},
				{Kind: token.Assign, Start: 47, End: 48, Val: "="},
				{Kind: token.Value, Start: 49, End: 52, Val: "GET"},
				{Kind: token.CloseCurl, Start: 57, End: 58, Val: "}"},
				{Kind: token.CloseCurl, Start: 62, End: 63, Val: "}"},
			},
		},
		{
			input: `service TestService {
				Ping() {}
			}`,
			output: Tokens{
				{Kind: token.Service, Start: 0, End: 7, Val: "service"},
				{Kind: token.Identifier, Start: 8, End: 19, Val: "TestService"},
				{Kind: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Kind: token.Identifier, Start: 26, End: 30, Val: "Ping"},
				{Kind: token.OpenParen, Start: 30, End: 31, Val: "("},
				{Kind: token.CloseParen, Start: 31, End: 32, Val: ")"},
				{Kind: token.OpenCurl, Start: 33, End: 34, Val: "{"},
				{Kind: token.CloseCurl, Start: 34, End: 35, Val: "}"},
				{Kind: token.CloseCurl, Start: 39, End: 40, Val: "}"},
			},
		},
		{
			input: `service TestService {
				Ping() => (status: stream int) {
					Method = GET
				}
			}`,
			output: Tokens{
				{Kind: token.Service, Start: 0, End: 7, Val: "service"},
				{Kind: token.Identifier, Start: 8, End: 19, Val: "TestService"},
				{Kind: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Kind: token.Identifier, Start: 26, End: 30, Val: "Ping"},
				{Kind: token.OpenParen, Start: 30, End: 31, Val: "("},
				{Kind: token.CloseParen, Start: 31, End: 32, Val: ")"},
				{Kind: token.Return, Start: 33, End: 35, Val: "=>"},
				{Kind: token.OpenParen, Start: 36, End: 37, Val: "("},
				{Kind: token.Identifier, Start: 37, End: 43, Val: "status"},
				{Kind: token.Colon, Start: 43, End: 44, Val: ":"},
				{Kind: token.Stream, Start: 45, End: 51, Val: "stream"},
				{Kind: token.Type, Start: 52, End: 55, Val: "int"},
				{Kind: token.CloseParen, Start: 55, End: 56, Val: ")"},
				{Kind: token.OpenCurl, Start: 57, End: 58, Val: "{"},
				{Kind: token.Identifier, Start: 64, End: 70, Val: "Method"},
				{Kind: token.Assign, Start: 71, End: 72, Val: "="},
				{Kind: token.Value, Start: 73, End: 76, Val: "GET"},
				{Kind: token.CloseCurl, Start: 81, End: 82, Val: "}"},
				{Kind: token.CloseCurl, Start: 86, End: 87, Val: "}"},
			},
		},
		{
			input: `service TestService {
				Ping() => (status: stream int) {}
			}`,
			output: Tokens{
				{Kind: token.Service, Start: 0, End: 7, Val: "service"},
				{Kind: token.Identifier, Start: 8, End: 19, Val: "TestService"},
				{Kind: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Kind: token.Identifier, Start: 26, End: 30, Val: "Ping"},
				{Kind: token.OpenParen, Start: 30, End: 31, Val: "("},
				{Kind: token.CloseParen, Start: 31, End: 32, Val: ")"},
				{Kind: token.Return, Start: 33, End: 35, Val: "=>"},
				{Kind: token.OpenParen, Start: 36, End: 37, Val: "("},
				{Kind: token.Identifier, Start: 37, End: 43, Val: "status"},
				{Kind: token.Colon, Start: 43, End: 44, Val: ":"},
				{Kind: token.Stream, Start: 45, End: 51, Val: "stream"},
				{Kind: token.Type, Start: 52, End: 55, Val: "int"},
				{Kind: token.CloseParen, Start: 55, End: 56, Val: ")"},
				{Kind: token.OpenCurl, Start: 57, End: 58, Val: "{"},
				{Kind: token.CloseCurl, Start: 58, End: 59, Val: "}"},
				{Kind: token.CloseCurl, Start: 63, End: 64, Val: "}"},
			},
		},
		{
			input: `service TestService {
				Ping() => (status: stream int)
			}`,
			output: Tokens{
				{Kind: token.Service, Start: 0, End: 7, Val: "service"},
				{Kind: token.Identifier, Start: 8, End: 19, Val: "TestService"},
				{Kind: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Kind: token.Identifier, Start: 26, End: 30, Val: "Ping"},
				{Kind: token.OpenParen, Start: 30, End: 31, Val: "("},
				{Kind: token.CloseParen, Start: 31, End: 32, Val: ")"},
				{Kind: token.Return, Start: 33, End: 35, Val: "=>"},
				{Kind: token.OpenParen, Start: 36, End: 37, Val: "("},
				{Kind: token.Identifier, Start: 37, End: 43, Val: "status"},
				{Kind: token.Colon, Start: 43, End: 44, Val: ":"},
				{Kind: token.Stream, Start: 45, End: 51, Val: "stream"},
				{Kind: token.Type, Start: 52, End: 55, Val: "int"},
				{Kind: token.CloseParen, Start: 55, End: 56, Val: ")"},
				{Kind: token.CloseCurl, Start: 60, End: 61, Val: "}"},
			},
		},
		{
			input: `service TestService {
				Ping() => (a: int, b: map<string, []int>)
			}`,
			output: Tokens{
				{Kind: token.Service, Start: 0, End: 7, Val: "service"},
				{Kind: token.Identifier, Start: 8, End: 19, Val: "TestService"},
				{Kind: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Kind: token.Identifier, Start: 26, End: 30, Val: "Ping"},
				{Kind: token.OpenParen, Start: 30, End: 31, Val: "("},
				{Kind: token.CloseParen, Start: 31, End: 32, Val: ")"},
				{Kind: token.Return, Start: 33, End: 35, Val: "=>"},
				{Kind: token.OpenParen, Start: 36, End: 37, Val: "("},
				{Kind: token.Identifier, Start: 37, End: 38, Val: "a"},
				{Kind: token.Colon, Start: 38, End: 39, Val: ":"},
				{Kind: token.Type, Start: 40, End: 43, Val: "int"},
				{Kind: token.Comma, Start: 43, End: 44, Val: ","},
				{Kind: token.Identifier, Start: 45, End: 46, Val: "b"},
				{Kind: token.Colon, Start: 46, End: 47, Val: ":"},
				{Kind: token.Type, Start: 48, End: 51, Val: "map"},
				{Kind: token.OpenAngle, Start: 51, End: 52, Val: "<"},
				{Kind: token.Type, Start: 52, End: 58, Val: "string"},
				{Kind: token.Comma, Start: 58, End: 59, Val: ","},
				{Kind: token.OpenBracket, Start: 60, End: 61, Val: "["},
				{Kind: token.CloseBracket, Start: 61, End: 62, Val: "]"},
				{Kind: token.Type, Start: 62, End: 65, Val: "int"},
				{Kind: token.CloseAngle, Start: 65, End: 66, Val: ">"},
				{Kind: token.CloseParen, Start: 66, End: 67, Val: ")"},
				{Kind: token.CloseCurl, Start: 71, End: 72, Val: "}"},
			},
		},
		{
			input: `service TestService {
				Ping() => (a: int, b: string)
			}`,
			output: Tokens{
				{Kind: token.Service, Start: 0, End: 7, Val: "service"},
				{Kind: token.Identifier, Start: 8, End: 19, Val: "TestService"},
				{Kind: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Kind: token.Identifier, Start: 26, End: 30, Val: "Ping"},
				{Kind: token.OpenParen, Start: 30, End: 31, Val: "("},
				{Kind: token.CloseParen, Start: 31, End: 32, Val: ")"},
				{Kind: token.Return, Start: 33, End: 35, Val: "=>"},
				{Kind: token.OpenParen, Start: 36, End: 37, Val: "("},
				{Kind: token.Identifier, Start: 37, End: 38, Val: "a"},
				{Kind: token.Colon, Start: 38, End: 39, Val: ":"},
				{Kind: token.Type, Start: 40, End: 43, Val: "int"},
				{Kind: token.Comma, Start: 43, End: 44, Val: ","},
				{Kind: token.Identifier, Start: 45, End: 46, Val: "b"},
				{Kind: token.Colon, Start: 46, End: 47, Val: ":"},
				{Kind: token.Type, Start: 48, End: 54, Val: "string"},
				{Kind: token.CloseParen, Start: 54, End: 55, Val: ")"},
				{Kind: token.CloseCurl, Start: 59, End: 60, Val: "}"},
			},
		},
		{
			input: `service TestService {
				Ping() => (a: int)
			}`,
			output: Tokens{
				{Kind: token.Service, Start: 0, End: 7, Val: "service"},
				{Kind: token.Identifier, Start: 8, End: 19, Val: "TestService"},
				{Kind: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Kind: token.Identifier, Start: 26, End: 30, Val: "Ping"},
				{Kind: token.OpenParen, Start: 30, End: 31, Val: "("},
				{Kind: token.CloseParen, Start: 31, End: 32, Val: ")"},
				{Kind: token.Return, Start: 33, End: 35, Val: "=>"},
				{Kind: token.OpenParen, Start: 36, End: 37, Val: "("},
				{Kind: token.Identifier, Start: 37, End: 38, Val: "a"},
				{Kind: token.Colon, Start: 38, End: 39, Val: ":"},
				{Kind: token.Type, Start: 40, End: 43, Val: "int"},
				{Kind: token.CloseParen, Start: 43, End: 44, Val: ")"},
				{Kind: token.CloseCurl, Start: 48, End: 49, Val: "}"},
			},
		},
		{
			input: `service TestService {
				Ping() => ()
			}`,
			output: Tokens{
				{Kind: token.Service, Start: 0, End: 7, Val: "service"},
				{Kind: token.Identifier, Start: 8, End: 19, Val: "TestService"},
				{Kind: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Kind: token.Identifier, Start: 26, End: 30, Val: "Ping"},
				{Kind: token.OpenParen, Start: 30, End: 31, Val: "("},
				{Kind: token.CloseParen, Start: 31, End: 32, Val: ")"},
				{Kind: token.Return, Start: 33, End: 35, Val: "=>"},
				{Kind: token.OpenParen, Start: 36, End: 37, Val: "("},
				{Kind: token.CloseParen, Start: 37, End: 38, Val: ")"},
				{Kind: token.CloseCurl, Start: 42, End: 43, Val: "}"},
			},
		},
		{
			input: `service TestService {
				Ping(a: int, b: map<string, int>)
			}`,
			output: Tokens{
				{Kind: token.Service, Start: 0, End: 7, Val: "service"},
				{Kind: token.Identifier, Start: 8, End: 19, Val: "TestService"},
				{Kind: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Kind: token.Identifier, Start: 26, End: 30, Val: "Ping"},
				{Kind: token.OpenParen, Start: 30, End: 31, Val: "("},
				{Kind: token.Identifier, Start: 31, End: 32, Val: "a"},
				{Kind: token.Colon, Start: 32, End: 33, Val: ":"},
				{Kind: token.Type, Start: 34, End: 37, Val: "int"},
				{Kind: token.Comma, Start: 37, End: 38, Val: ","},
				{Kind: token.Identifier, Start: 39, End: 40, Val: "b"},
				{Kind: token.Colon, Start: 40, End: 41, Val: ":"},
				{Kind: token.Type, Start: 42, End: 45, Val: "map"},
				{Kind: token.OpenAngle, Start: 45, End: 46, Val: "<"},
				{Kind: token.Type, Start: 46, End: 52, Val: "string"},
				{Kind: token.Comma, Start: 52, End: 53, Val: ","},
				{Kind: token.Type, Start: 54, End: 57, Val: "int"},
				{Kind: token.CloseAngle, Start: 57, End: 58, Val: ">"},
				{Kind: token.CloseParen, Start: 58, End: 59, Val: ")"},
				{Kind: token.CloseCurl, Start: 63, End: 64, Val: "}"},
			},
		},
		{
			input: `service TestService {
				Ping(a: int)
			}`,
			output: Tokens{
				{Kind: token.Service, Start: 0, End: 7, Val: "service"},
				{Kind: token.Identifier, Start: 8, End: 19, Val: "TestService"},
				{Kind: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Kind: token.Identifier, Start: 26, End: 30, Val: "Ping"},
				{Kind: token.OpenParen, Start: 30, End: 31, Val: "("},
				{Kind: token.Identifier, Start: 31, End: 32, Val: "a"},
				{Kind: token.Colon, Start: 32, End: 33, Val: ":"},
				{Kind: token.Type, Start: 34, End: 37, Val: "int"},
				{Kind: token.CloseParen, Start: 37, End: 38, Val: ")"},
				{Kind: token.CloseCurl, Start: 42, End: 43, Val: "}"},
			},
		},
		{
			input: `service TestService {
				Ping()
			}`,
			output: Tokens{
				{Kind: token.Service, Start: 0, End: 7, Val: "service"},
				{Kind: token.Identifier, Start: 8, End: 19, Val: "TestService"},
				{Kind: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Kind: token.Identifier, Start: 26, End: 30, Val: "Ping"},
				{Kind: token.OpenParen, Start: 30, End: 31, Val: "("},
				{Kind: token.CloseParen, Start: 31, End: 32, Val: ")"},
				{Kind: token.CloseCurl, Start: 36, End: 37, Val: "}"},
			},
		},
		{
			input: `service TestService {}`,
			output: Tokens{
				{Kind: token.Service, Start: 0, End: 7, Val: "service"},
				{Kind: token.Identifier, Start: 8, End: 19, Val: "TestService"},
				{Kind: token.OpenCurl, Start: 20, End: 21, Val: "{"},
				{Kind: token.CloseCurl, Start: 21, End: 22, Val: "}"},
			},
		},
	})
}
