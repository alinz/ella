package lexer

import (
	"github.com/alinz/rpc.go/schema/token"
)

func Service(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		l.AcceptRunUntil(" \t\r\n#")
		if l.Current() != "service" {
			l.Errorf("expected message keywoard but got %s", l.Current())
			return nil
		}
		l.Emit(token.Service)

		IgnoreWhiteSpace(l)
		checkComment(l)

		l.AcceptRunUntil(" \t\r\n{#")
		if l.Current() == "" {
			l.Errorf("expected name for service but got nothing")
			return nil
		}
		l.Emit(token.Identifier)

		return ServiceMethods(next)
	}
}

func ServiceMethods(next StateFn) StateFn {
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
			l.Errorf("expected '{' for service body but got %s", l.Current())
			return nil
		}
		l.Emit(token.OpenCurl)

		return ServiceMethod(next)
	}
}

func ServiceMethod(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == '}' {
			l.Next()
			l.Emit(token.CloseCurl)
			return next
		}

		l.AcceptRunUntil(" \t\r\n(#")
		if l.Current() == "" {
			l.Errorf("expected method name but got nothing")
			return nil
		}
		l.Emit(token.Identifier)

		IgnoreWhiteSpace(l)
		checkComment(l)

		if value := l.Peek(); value != '(' {
			l.Errorf("expected '(' for method arguments but got %s", string(value))
			return nil
		}
		l.Next()
		l.Emit(token.OpenParen)

		return ServiceMethodArgs(next)
	}
}

func ServiceMethodArgs(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == ')' {
			l.Next()
			l.Emit(token.CloseParen)
			return ServiceMethodReturns(next)
		}

		return ServiceMethodArg(next)
	}
}

func ServiceMethodArg(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == ')' {
			l.Next()
			l.Emit(token.CloseParen)
			return ServiceMethodReturns(next)
		}

		IgnoreWhiteSpace(l)
		checkComment(l)

		l.AcceptRunUntil(" \t\r\n:?#")
		if l.Current() == "" {
			l.Errorf("expected argument name but got nothing")
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
			l.Errorf("expected ':' for argument type but got %s", string(l.Peek()))
			return nil
		}
		l.Next()
		l.Emit(token.Colon)

		return ServiceMethodArgsTypes(next)
	}
}

func ServiceMethodArgsTypes(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		switch l.Peek() {
		case ')':
			return ServiceMethodArg(next)
		case '[':
			l.Next()
			l.Emit(token.OpenBracket)
			return ServiceMethodArgsTypes(next)
		case ']':
			l.Next()
			l.Emit(token.CloseBracket)
			return ServiceMethodArgsTypes(next)
		case '<':
			l.Counter++
			l.Next()
			l.Emit(token.OpenAngle)
			return ServiceMethodArgsTypes(next)
		case '>':
			l.Counter--
			l.Next()
			l.Emit(token.CloseAngle)
			return ServiceMethodArgsTypes(next)
		case ',':
			l.Next()
			l.Emit(token.Comma)
			if l.Counter == 0 {
				return ServiceMethodArg(next)
			}
			return ServiceMethodArgsTypes(next)
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
			return ServiceMethodArgsTypes(next)
		default:
			l.AcceptRunUntil(" .<>[]),\t\n\r")
			l.Emit(token.Type)
			return ServiceMethodArgsTypes(next)
		}
	}
}

func ServiceMethodOptions(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == '}' {
			l.Next()
			l.Emit(token.CloseCurl)
			return ServiceMethod(next)
		}

		// identifier -> = -> value
		// identifier -> =

		l.AcceptRunUntil(" \t\n\r#")
		if l.Current() == "" {
			l.Errorf("expected method option name but got nothing")
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

		switch l.Peek() {
		case '"':
			l.Next()
			l.Ignore()

			l.AcceptRunUntil("\"\n\r")
			if l.Peek() != '"' {
				l.Errorf("expected '\"' but got %s", string(l.Peek()))
				return nil
			}
			l.Emit(token.Value)

			l.Next()
			l.Ignore()

		case '\'':
			l.Next()
			l.Ignore()

			l.AcceptRunUntil("'\n\r")
			if l.Peek() != '\'' {
				l.Errorf("expected ' but got %s", string(l.Peek()))
				return nil
			}
			l.Emit(token.Value)

			l.Next()
			l.Ignore()

		default:
			l.AcceptRunUntil(" \t\n\r#")
			l.Emit(token.Value)
		}

		return ServiceMethodOptions(next)
	}
}

func ServiceMethodReturns(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		value := l.Peek()

		switch value {
		case '=':
		case '{':
			l.Next()
			l.Emit(token.OpenCurl)
			return ServiceMethodOptions(next)
		default:
			// there is no return ( => ) for example Ping(),
			// that's why we check for next method
			return ServiceMethod(next)
		}

		l.Next() // consume '='
		value = l.Next()
		if value != '>' {
			l.Errorf("expected '=>' for method return but got %s", string(value))
			return nil
		}
		l.Emit(token.Return)

		IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() != '(' {
			l.Errorf("expected '(' for method return args but got %s", string(l.Peek()))
			return nil
		}
		l.Next()
		l.Emit(token.OpenParen)

		return ServiceMethodReturn(next)
	}
}

func ServiceMethodReturn(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == ')' {
			l.Next()
			l.Emit(token.CloseParen)
			return ServiceMethod(next)
		}

		IgnoreWhiteSpace(l)
		checkComment(l)

		l.AcceptRunUntil(" \t\r\n:?#")
		if l.Current() == "" {
			l.Errorf("expected return argument name but got nothing")
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
			l.Errorf("expected ':' for argument type but got %s", string(l.Peek()))
			return nil
		}
		l.Next()
		l.Emit(token.Colon)

		return ServiceMethodReturnTypes(next)
	}
}

func ServiceMethodReturnTypes(next StateFn) StateFn {
	return func(l *Lexer) StateFn {
		IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == 's' {
			value := l.PeekRunUntil(" \t\r")
			if value == "stream" {
				l.AcceptRunUntil(" \t\r")
				l.Emit(token.Stream)
				return ServiceMethodReturnTypes(next)
			}
		}

		switch l.Peek() {
		case ')':
			return ServiceMethodArg(next)
		case '[':
			l.Next()
			l.Emit(token.OpenBracket)
			return ServiceMethodReturnTypes(next)
		case ']':
			l.Next()
			l.Emit(token.CloseBracket)
			return ServiceMethodReturnTypes(next)
		case '<':
			l.Counter++
			l.Next()
			l.Emit(token.OpenAngle)
			return ServiceMethodReturnTypes(next)
		case '>':
			l.Counter--
			l.Next()
			l.Emit(token.CloseAngle)
			return ServiceMethodReturnTypes(next)
		case ',':
			l.Next()
			l.Emit(token.Comma)
			if l.Counter == 0 {
				return ServiceMethodReturn(next)
			}
			return ServiceMethodReturnTypes(next)
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
			return ServiceMethodReturnTypes(next)
		default:
			l.AcceptRunUntil(" .<>[]),\t\n\r")
			l.Emit(token.Type)
			return ServiceMethodReturnTypes(next)
		}
	}
}
