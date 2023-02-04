package lex

import (
	"github.com/alinz/rpc.go/pkg/lexer"
	"github.com/alinz/rpc.go/schema/lex/token"
)

func Stmt(next lexer.State) lexer.State {
	var internal lexer.State

	internal = func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)

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
