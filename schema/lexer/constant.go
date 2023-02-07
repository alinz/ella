package lexer

import (
	"github.com/alinz/rpc.go/schema/token"
)

func Constant(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)

		l.AcceptRunUntil(" =\t#")
		if l.Current() == "" {
			l.Errorf("expected name for constant but got nothing")
			return nil
		}
		l.Emit(token.Identifier)

		IgnoreSpaceTabs(l)

		checkComment(l)

		if value := l.Next(); value != '=' {
			l.Errorf("expected '=' but got %s", string(value))
			return nil
		}
		l.Emit(token.Assign)

		IgnoreSpaceTabs(l)

		switch l.Peek() {
		case '"':
			l.Next()
			l.Ignore()

			l.AcceptRunUntil("\r\n\"")

			if l.Peek() != '"' {
				l.Errorf("expected \" but got %s", string(l.Peek()))
				return nil
			}

			l.Emit(token.Value)

			l.Next()
			l.Ignore()
		case '\'':
			l.Next()
			l.Ignore()

			l.AcceptRunUntil("\r\n'")

			if l.Peek() != '\'' {
				l.Errorf("expected ' but got %s", string(l.Peek()))
				return nil
			}

			l.Emit(token.Value)

			l.Next()
			l.Ignore()

		default:
			l.AcceptRunUntil(" \r\n#")
			l.Emit(token.Value)
		}

		checkComment(l)

		return next
	}
}
