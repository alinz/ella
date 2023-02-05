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

const (
	Error Kind = -1   // Error token type which indicates error
	EOF   Kind = iota // EOF token type which indicates end of input
	Identifier
	Assign // =
	Value  // anything after assign char
	Type
	Enum
	Message
	Service
	Stream
	Comment      // #
	Colon        // :
	Comma        // ,
	Underline    // _
	Optional     // ?
	Ellipsis     // ...
	Return       // =>
	OpenCurl     // {
	CloseCurl    // }
	OpenParen    // (
	CloseParen   // )
	OpenAngle    // <
	CloseAngle   // >
	OpenBracket  // [
	CloseBracket // ]
)

var names = map[Kind]string{
	Error:        "Error",
	EOF:          "EOF",
	Identifier:   "Identifier",
	Assign:       "Assign",
	Value:        "Value",
	Type:         "Type",
	Enum:         "Enum",
	Message:      "Message",
	Service:      "Service",
	Stream:       "Stream",
	Comment:      "Comment",
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

func (k Kind) String() string {
	if name, ok := names[k]; ok {
		return name
	}
	return "Unknown"
}
