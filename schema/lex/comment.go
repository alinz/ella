package lex

import (
	"github.com/alinz/rpc.go/pkg/lexer"
	"github.com/alinz/rpc.go/schema/lex/token"
)

func Comment(next lexer.State) lexer.State {
	return func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)

		value := l.Next()
		if value != '#' {
			errorf(l, "expected '#' for comment but got %s", string(value))
			return nil
		}
		l.Emit(token.Comment)

		l.AcceptRunUntil("\n\r")
		l.Emit(token.Value)

		lexer.IgnoreWhiteSpace(l)

		return next
	}
}

func checkComment(l *lexer.Lexer) {
	if l.Peek() != '#' {
		return
	}

	Comment(nil)(l)
}
