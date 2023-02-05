package lexer

import (
	"github.com/alinz/rpc.go/schema/token"
)

func Comment(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)

		value := l.Next()
		if value != '#' {
			l.Errorf("expected '#' for comment but got %s", string(value))
			return nil
		}
		l.Emit(token.Comment)

		l.AcceptRunUntil("\n\r")
		l.Emit(token.Value)

		IgnoreWhiteSpace(l)

		return next
	}
}

func checkComment(l *Lexer) {
	if l.Peek() != '#' {
		return
	}
	Comment(nil)(l)
}
