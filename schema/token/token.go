package token

import (
	"fmt"
)

type Kind int

var _ fmt.Stringer = Kind(0)

type Token struct {
	Val   string
	Kind  Kind
	Start int
	End   int
}

func (t Token) String() string {
	return fmt.Sprintf("Kind: %s, Val: %s, Start: %d, End: %d", t.Kind, t.Val, t.Start, t.End)
}

func (t Token) OneOf(kinds ...Kind) bool {
	for _, kind := range kinds {
		if t.Kind == kind {
			return true
		}
	}
	return false
}

const (
	Error Kind = -1   // Error token type which indicates error
	EOF   Kind = iota // EOF token type which indicates end of input
	Word              // Word token type which indicates word,identifier,keyword,value
	ConstantNumber
	ConstantString
	Enum        // enum
	Message     // message
	Service     // service
	Stream      // stream
	Map         // map
	String      // string
	Byte        // byte
	Bool        // bool
	Int8        // int8
	Int16       // int16
	Int32       // int32
	Int64       // int64
	Uint8       // uint8
	Uint16      // uint16
	Uint32      // uint32
	Uint64      // uint64
	Float32     // float32
	Float64     // float64
	Timestamp   // timestamp
	Any         // any
	Assign      // =
	Colon       // :
	Comma       // ,
	Dot         // .
	OpenCurly   // {
	CloseCurly  // }
	OpenParen   // (
	CloseParen  // )
	OpenAngle   // <
	CloseAngle  // >
	OpenSquare  // [
	CloseSquare // ]
	Comment     // #
)

var names = map[Kind]string{
	Error:          "Error",
	EOF:            "EOF",
	Word:           "Word",
	ConstantNumber: "ConstantNumber",
	ConstantString: "ConstantString",
	Enum:           "Enum",
	Message:        "Message",
	Service:        "Service",
	Stream:         "Stream",
	Map:            "Map",
	String:         "String",
	Byte:           "Byte",
	Bool:           "Bool",
	Int8:           "Int8",
	Int16:          "Int16",
	Int32:          "Int32",
	Int64:          "Int64",
	Uint8:          "Uint8",
	Uint16:         "Uint16",
	Uint32:         "Uint32",
	Uint64:         "Uint64",
	Float32:        "Float32",
	Float64:        "Float64",
	Timestamp:      "Timestamp",
	Any:            "Any",
	Assign:         "Assign",
	Colon:          "Colon",
	Comma:          "Comma",
	Dot:            "Dot",
	OpenCurly:      "OpenCurly",
	CloseCurly:     "CloseCurly",
	OpenParen:      "OpenParen",
	CloseParen:     "CloseParen",
	OpenAngle:      "OpenAngle",
	CloseAngle:     "CloseAngle",
	OpenSquare:     "OpenSquare",
	CloseSquare:    "CloseSquare",
	Comment:        "Comment",
}

func (k Kind) String() string {
	if name, ok := names[k]; ok {
		return name
	}
	return "Unknown"
}
