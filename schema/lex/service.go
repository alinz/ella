package lex

import (
	"github.com/alinz/rpc.go/pkg/lexer"
	"github.com/alinz/rpc.go/schema/lex/token"
)

func Service(next lexer.State) lexer.State {
	return func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		l.AcceptRunUntil(" \t\r\n#")
		if l.Current() != "service" {
			errorf(l, "expected message keywoard but got %s", l.Current())
			return nil
		}
		l.Emit(token.Service)

		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		l.AcceptRunUntil(" \t\r\n{#")
		if l.Current() == "" {
			errorf(l, "expected name for service but got nothing")
			return nil
		}
		l.Emit(token.Identifier)

		return ServiceMethods(next)
	}
}

func ServiceMethods(next lexer.State) lexer.State {
	return func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == '}' {
			l.Next()
			l.Emit(token.CloseCurl)
			return next
		}

		l.Next()
		if l.Current() != "{" {
			errorf(l, "expected '{' for service body but got %s", l.Current())
			return nil
		}
		l.Emit(token.OpenCurl)

		return ServiceMethod(next)
	}
}

func ServiceMethod(next lexer.State) lexer.State {
	return func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == '}' {
			l.Next()
			l.Emit(token.CloseCurl)
			return next
		}

		l.AcceptRunUntil(" \t\r\n(#")
		if l.Current() == "" {
			errorf(l, "expected method name but got nothing")
			return nil
		}
		l.Emit(token.Identifier)

		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		if value := l.Peek(); value != '(' {
			errorf(l, "expected '(' for method arguments but got %s", string(value))
			return nil
		}
		l.Next()
		l.Emit(token.OpenParen)

		return ServiceMethodArgs(next)
	}
}

func ServiceMethodArgs(next lexer.State) lexer.State {
	return func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == ')' {
			l.Next()
			l.Emit(token.CloseParen)
			return ServiceMethodReturns(next)
		}

		return ServiceMethodArg(next)
	}
}

func ServiceMethodArg(next lexer.State) lexer.State {
	return func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == ')' {
			l.Next()
			l.Emit(token.CloseParen)
			return ServiceMethodReturns(next)
		}

		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		l.AcceptRunUntil(" \t\r\n:?#")
		if l.Current() == "" {
			errorf(l, "expected argument name but got nothing")
			return nil
		}
		l.Emit(token.Identifier)

		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == '?' {
			l.Next()
			l.Emit(token.Optional)
		}

		if l.Peek() != ':' {
			errorf(l, "expected ':' for argument type but got %s", string(l.Peek()))
			return nil
		}
		l.Next()
		l.Emit(token.Colon)

		return ServiceMethodArgsTypes(next)
	}
}

func ServiceMethodArgsTypes(next lexer.State) lexer.State {
	return func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)
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
				errorf(l, "expected '.' but got %s", string(l.Peek()))
				return nil
			}
			l.Next()
			if l.Peek() != '.' {
				errorf(l, "expected '.' but got %s", string(l.Peek()))
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

func ServiceMethodReturns(next lexer.State) lexer.State {
	return func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		value := l.Peek()
		if value != '=' {
			// there is no return ( => ) for example Ping(),
			// that's why we check for next method
			return ServiceMethod(next)
		}
		l.Next() // consume '='
		value = l.Next()
		if value != '>' {
			errorf(l, "expected '=>' for method return but got %s", string(value))
			return nil
		}
		l.Emit(token.Return)

		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() != '(' {
			errorf(l, "expected '(' for method return args but got %s", string(l.Peek()))
			return nil
		}
		l.Next()
		l.Emit(token.OpenParen)

		return ServiceMethodReturn(next)
	}
}

func ServiceMethodReturn(next lexer.State) lexer.State {
	return func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == ')' {
			l.Next()
			l.Emit(token.CloseParen)
			return ServiceMethod(next)
		}

		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		l.AcceptRunUntil(" \t\r\n:?#")
		if l.Current() == "" {
			errorf(l, "expected return argument name but got nothing")
			return nil
		}
		l.Emit(token.Identifier)

		lexer.IgnoreWhiteSpace(l)
		checkComment(l)

		if l.Peek() == '?' {
			l.Next()
			l.Emit(token.Optional)
		}

		if l.Peek() != ':' {
			errorf(l, "expected ':' for argument type but got %s", string(l.Peek()))
			return nil
		}
		l.Next()
		l.Emit(token.Colon)

		return ServiceMethodReturnTypes(next)
	}
}

func ServiceMethodReturnTypes(next lexer.State) lexer.State {
	return func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)
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
				errorf(l, "expected '.' but got %s", string(l.Peek()))
				return nil
			}
			l.Next()
			if l.Peek() != '.' {
				errorf(l, "expected '.' but got %s", string(l.Peek()))
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
