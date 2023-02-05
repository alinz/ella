package lexer

import (
	"github.com/alinz/rpc.go/schema/token"
)

func Message(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		l.AcceptRunUntil(" \t\r\n#")
		if l.Current() != "message" {
			l.Errorf("expected message keywoard but got %s", l.Current())
			return nil
		}
		l.Emit(token.Message)

		IgnoreWhiteSpace(l)
		checkComment(l)

		l.AcceptRunUntil(" \t\r\n{#")
		if l.Current() == "" {
			l.Errorf("expected name for message but got nothing")
			return nil
		}
		l.Emit(token.Identifier)

		return MessageFields(next)
	}
}

func MessageFields(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == '}' {
			l.Next()
			l.Emit(token.CloseCurl)
			return next
		}

		l.Next()
		if l.Current() != "{" {
			l.Errorf("expected '{' for message body but got %s", l.Current())
			return nil
		}
		l.Emit(token.OpenCurl)

		return MessageField(next)
	}
}

func MessageField(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == '}' {
			l.Next()
			l.Emit(token.CloseCurl)
			return next
		}

		if l.Peek() == '.' {
			return MessageTypes(next)
		}

		l.AcceptRunUntil(" ?:{\t\n\r#")
		if l.Current() == "" {
			l.Errorf("expected field name but got nothing")
			return nil
		}
		l.Emit(token.Identifier)

		IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == '?' {
			l.Next()
			l.Emit(token.Optional)
		}

		if l.Peek() != ':' {
			l.Errorf("expected ':' but got %s", string(l.Peek()))
			return nil
		}
		l.Next()
		l.Emit(token.Colon)

		return MessageTypes(next)
	}
}

func MessageTypes(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		l.AcceptRun(" \t")
		l.Ignore()

		switch l.Peek() {
		case '\n', '\r', 0:
			return MessageField(next)
		case '[':
			l.Next()
			l.Emit(token.OpenBracket)
			return MessageTypes(next)
		case ']':
			l.Next()
			l.Emit(token.CloseBracket)
			return MessageTypes(next)
		case '<':
			l.Next()
			l.Emit(token.OpenAngle)
			return MessageTypes(next)
		case '>':
			l.Next()
			l.Emit(token.CloseAngle)
			return MessageTypes(next)
		case ',':
			l.Next()
			l.Emit(token.Comma)
			return MessageTypes(next)
		case '.':
			l.Next()
			if l.Peek() != '.' {
				l.Errorf("expected '.' but got %s", string(l.Peek()))
				return nil
			}
			l.Next()
			if l.Peek() != '.' {
				l.Errorf("expected '.' but got %s", string(l.Peek()))
				return nil
			}
			l.Next()
			l.Emit(token.Ellipsis)
			return MessageTypes(next)
		case '{':
			l.Next()
			l.Emit(token.OpenCurl)
			return FieldOptions(next)
		default:
			l.AcceptRunUntil(" .<>[]{,\t\n\r")
			l.Emit(token.Type)
			return MessageTypes(next)
		}
	}
}

func FieldOptions(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == '}' {
			l.Next()
			l.Emit(token.CloseCurl)
			return MessageField(next)
		}

		// identifier -> = -> value
		// identifier -> =

		l.AcceptRunUntil(" \t\n\r#")
		if l.Current() == "" {
			l.Errorf("expected field option name but got nothing")
			return nil
		}
		l.Emit(token.Identifier)

		l.AcceptRun(" \t")
		l.Ignore()
		checkComment(l)

		if value := l.Next(); value != '=' {
			l.Errorf("expected '=' but got %s", string(value))
			return nil
		}
		l.Emit(token.Assign)

		l.AcceptRun(" \t}")
		l.Ignore()

		l.AcceptRunUntil(" \t\n\r#")
		l.Emit(token.Value)

		return FieldOptions(next)
	}
}
