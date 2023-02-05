package lexer

import (
	"github.com/alinz/rpc.go/schema/token"
)

func Enum(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)

		l.AcceptRunUntil(" \t\r\n#")
		if l.Current() == "" {
			l.Errorf("expected enum for enum but got nothing")
			return nil
		}
		l.Emit(token.Enum)

		IgnoreWhiteSpace(l)

		checkComment(l)

		l.AcceptRunUntil(" \t\r\n#")
		if l.Current() == "" {
			l.Errorf("expected name for enum but got nothing")
			return nil
		}
		l.Emit(token.Identifier)

		IgnoreWhiteSpace(l)
		checkComment(l)

		l.AcceptRunUntil(" \t\r\n{#")

		if l.Current() == "" {
			l.Errorf("expected type of enum but got %s", l.Current())
			return nil
		}
		l.Emit(token.Type)

		checkComment(l)

		return EnumValues(next)
	}
}

func EnumValues(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		if value := l.Peek(); value != '{' {
			l.Errorf("expected '{' but got %s", string(value))
			return nil
		}
		l.Next()
		l.Emit(token.OpenCurl)

		return EnumValue(next)
	}
}

func EnumValue(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		l.AcceptRunUntil(" \t\n\r#")
		value := l.Current()

		switch value {
		case "}":
			l.Emit(token.CloseCurl)
			return next
		case "=":
			l.Emit(token.Assign)

			IgnoreWhiteSpace(l)
			checkComment(l)

			l.AcceptRunUntil(" \t\n\r}#")
			l.Emit(token.Value)

			return EnumValue(next)
		default:
			if value != "" {
				l.Emit(token.Identifier)
			}
			return EnumValue(next)
		}
	}
}
