package lexer

import (
	"github.com/alinz/rpc.go/schema/token"
)

func Stmt(next StateFn) StateFn {
	var internal StateFn

	internal = func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)

		if l.Peek() == 0 {
			l.Emit(token.EOF)
			return next
		}

		if l.Peek() == '#' {
			return Comment(internal)
		}

		peekValue := l.PeekRunUntil(" \t\r\n")

		switch peekValue {
		case "message":
			return Message(internal)
		case "service":
			return Service(internal)
		case "enum":
			return Enum(internal)
		case "":
			return next
		default:
			return Constant(internal)
		}
	}

	return internal
}
