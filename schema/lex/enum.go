package lex

import (
	"github.com/alinz/rpc.go/pkg/lexer"
	"github.com/alinz/rpc.go/schema/lex/token"
)

func Enum(next lexer.State) lexer.State {
	return func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)

		l.AcceptRunUntil(" \t\r\n#")
		if l.Current() == "" {
			errorf(l, "expected enum for enum but got nothing")
			return nil
		}
		l.Emit(token.Enum)

		lexer.IgnoreWhiteSpace(l)

		checkComment(l)

		l.AcceptRunUntil(" \t\r\n#")
		if l.Current() == "" {
			errorf(l, "expected name for enum but got nothing")
			return nil
		}
		l.Emit(token.Identifier)

		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		l.AcceptRunUntil(" \t\r\n{#")

		if l.Current() == "" {
			errorf(l, "expected type of enum but got %s", l.Current())
			return nil
		}
		l.Emit(token.Type)

		checkComment(l)

		return EnumValues(next)
	}
}

func EnumValues(next lexer.State) lexer.State {
	return func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		if value := l.Peek(); value != '{' {
			errorf(l, "expected '{' but got %s", string(value))
			return nil
		}
		l.Next()
		l.Emit(token.OpenCurl)

		return EnumValue(next)
	}
}

func EnumValue(next lexer.State) lexer.State {
	return func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		l.AcceptRunUntil(" \t\n\r#")
		value := l.Current()

		switch value {
		case "}":
			l.Emit(token.CloseCurl)
			return next
		case "=":
			l.Emit(token.Assign)

			lexer.IgnoreWhiteSpace(l)
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
