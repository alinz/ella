package token

import "github.com/alinz/rpc.go/pkg/lexer"

const (
	Error lexer.Type = -1   // Error token type which indicates error
	EOF   lexer.Type = iota // EOF token type which indicates end of input

	Identifier

	Assign // =
	Value  // anything after assign char
	Type

	Enum
	Message
	Service
	Stream

	// Byte
	// Bool
	// Any
	// Null
	// Uint8
	// Uint16
	// Uint32
	// Uint64
	// Int8
	// Int16
	// Int32
	// Int64
	// Float32
	// Float64
	// String
	// Timestamp

	Colon     // :
	Comma     // ,
	Underline // _
	Optional  // ?
	Ellipsis  // ...
	Return    // =>

	OpenCurl     // {
	CloseCurl    // }
	OpenParen    // (
	CloseParen   // )
	OpenAngle    // <
	CloseAngle   // >
	OpenBracket  // [
	CloseBracket // ]
)

var names = map[lexer.Type]string{
	EOF:        "EOF",
	Identifier: "Identifier",
	Assign:     "Assign",
	Value:      "Value",
	Type:       "Type",
	Enum:       "Enum",
	Message:    "Message",
	Service:    "Service",
	Stream:     "Stream",
	// Byte:       "Byte",
	// Bool:       "Bool",
	// Any:        "Any",
	// Null:       "Null",
	// Uint8:      "Uint8",
	// Uint16:     "Uint16",
	// Uint32:     "Uint32",
	// Uint64:     "Uint64",
	// Int8:       "Int8",
	// Int16:      "Int16",
	// Int32:      "Int32",
	// Int64:      "Int64",
	// Float32:    "Float32",
	// Float64:    "Float64",
	// String:     "String",
	// Timestamp:  "Timestamp",
	Colon:        "Colon",
	Comma:        "Comma",
	Underline:    "Underline",
	Optional:     "Optional",
	Ellipsis:     "Ellipsis",
	Return:       "Return",
	OpenCurl:     "OpenCurl",
	CloseCurl:    "CloseCurl",
	OpenParen:    "OpenParen",
	CloseParen:   "CloseParen",
	OpenAngle:    "OpenAngle",
	CloseAngle:   "CloseAngle",
	OpenBracket:  "OpenBracket",
	CloseBracket: "CloseBracket",
}

func Name(t lexer.Type) string {
	if name, ok := names[t]; ok {
		return name
	}
	return "Unknown"
}
