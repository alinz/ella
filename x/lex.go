package lex

import (
	"github.com/alinz/rpc.go/pkg/lexer"
	"github.com/alinz/rpc.go/schema/lex/token"
)

func errorf(l *lexer.Lexer, format string, args ...interface{}) {
	l.Errorf(token.Error, format, args...)
}

func Stmt(l *lexer.Lexer) lexer.State {
	lexer.IgnoreWhiteSpace(l)

	val := l.Next()

	if val == 0 {
		l.Emit(token.EOF)
		return nil
	}

	l.AcceptRunUntil(" =\t\n\r")

	value := l.Current()

	switch value {
	case "enum":
		l.Emit(token.Enum)
		return Enum(l)
	case "message":
		l.Emit(token.Message)
		return Message(l)
	case "service":
		l.Emit(token.Service)
		return Service(l)
	default:
		l.Emit(token.Identifier)
		return Constant(l)
	}
}

func Constant(l *lexer.Lexer) lexer.State {
	lexer.IgnoreWhiteSpace(l)

	if value := l.Next(); value != '=' {
		errorf(l, "expected '=' but got %s", string(value))
		return nil
	}
	l.Emit(token.Assign)

	lexer.IgnoreWhiteSpace(l)

	l.AcceptRunUntil(" \t\n\r")

	if l.Current() == "" {
		errorf(l, "expected value for constant but got nothing")
		return nil
	}
	l.Emit(token.Value)

	return Stmt
}

func Enum(l *lexer.Lexer) lexer.State {
	lexer.IgnoreWhiteSpace(l)

	l.AcceptRunUntil(" \t\n\r")
	value := l.Current()
	if value == "" {
		errorf(l, "expected enum name but got nothing")
		return nil
	}
	l.Emit(token.Identifier)

	lexer.IgnoreWhiteSpace(l)

	l.AcceptRunUntil(" {\t\n\r")
	value = l.Current()
	if value == "" {
		errorf(l, "expect enum type but got nothing")
		return nil
	}
	l.Emit(token.Type)

	lexer.IgnoreWhiteSpace(l)

	l.AcceptRunUntil(" \t\n\r")
	value = l.Current()
	if value != "{" {
		errorf(l, "expected '{' but got %s", value)
		return nil
	}
	l.Emit(token.OpenCurl)

	return EnumValues
}

func EnumValues(l *lexer.Lexer) lexer.State {
	lexer.IgnoreWhiteSpace(l)

	value := l.PeekRunUntil(" \t\n\r")

	switch value {
	case "}":
		l.Emit(token.CloseCurl)
		return Stmt
	case "":
		errorf(l, "expected enum value but got nothing")
		return nil
	case "=":
		l.Emit(token.Assign)
		return EnumValue
	default:
		l.Emit(token.Identifier)
		return EnumValues
	}
}

func EnumValue(l *lexer.Lexer) lexer.State {
	lexer.IgnoreWhiteSpace(l)

	l.AcceptRunUntil(" \t\n\r")
	value := l.Current()

	switch value {
	case "":
		errorf(l, "expected enum value but got nothing")
		return nil

	case "}":
		errorf(l, "expected enum value but got %s", value)
		return nil
	}

	l.Emit(token.Value)
	return EnumValues
}

func Message(l *lexer.Lexer) lexer.State {
	lexer.IgnoreWhiteSpace(l)

	l.AcceptRunUntil(" {\t\n\r")
	value := l.Current()
	if value == "" {
		errorf(l, "expected message name but got nothing")
		return nil
	}
	l.Emit(token.Identifier)

	lexer.IgnoreWhiteSpace(l)
	if l.Peek() != '{' {
		errorf(l, "expected '{' but got %v", l.Peek())
		return nil
	}
	l.Next()
	l.Emit(token.OpenCurl)

	return MessageFields
}

func MessageFields(l *lexer.Lexer) lexer.State {
	lexer.IgnoreWhiteSpace(l)

	if l.Peek() == '}' {
		l.Next()
		l.Emit(token.CloseCurl)
		return Stmt
	}

	if l.Peek() == '.' {
		return MessageTypes
	}

	l.AcceptRunUntil(" ?:{\t\n\r")
	if l.Current() == "" {
		errorf(l, "expected field name but got nothing")
		return nil
	}
	l.Emit(token.Identifier)

	lexer.IgnoreWhiteSpace(l)

	if l.Peek() == '?' {
		l.Next()
		l.Emit(token.Optional)
	}

	if l.Peek() != ':' {
		errorf(l, "expected ':' but got %s", string(l.Peek()))
		return nil
	}
	l.Next()
	l.Emit(token.Colon)

	return MessageTypes
}

func FieldOptions(l *lexer.Lexer) lexer.State {
	lexer.IgnoreWhiteSpace(l)

	if l.Peek() == '}' {
		l.Next()
		l.Emit(token.CloseCurl)
		return MessageFields
	}

	// identifier -> = -> value
	// identifier -> =

	l.AcceptRunUntil(" \t\n\r")
	if l.Current() == "" {
		errorf(l, "expected field option name but got nothing")
		return nil
	}
	l.Emit(token.Identifier)

	l.AcceptRun(" \t")
	l.Ignore()

	if value := l.Next(); value != '=' {
		errorf(l, "expected '=' but got %s", string(value))
		return nil
	}
	l.Emit(token.Assign)

	l.AcceptRun(" \t}")
	l.Ignore()

	l.AcceptRunUntil(" \t\n\r")
	l.Emit(token.Value)

	return FieldOptions
}

func MessageTypes(l *lexer.Lexer) lexer.State {
	l.AcceptRun(" \t")
	l.Ignore()

	switch l.Peek() {
	case '\n', '\r', 0:
		return MessageFields
	case '[':
		l.Next()
		l.Emit(token.OpenBracket)
		return MessageTypes
	case ']':
		l.Next()
		l.Emit(token.CloseBracket)
		return MessageTypes
	case '<':
		l.Next()
		l.Emit(token.OpenAngle)
		return MessageTypes
	case '>':
		l.Next()
		l.Emit(token.CloseAngle)
		return MessageTypes
	case ',':
		l.Next()
		l.Emit(token.Comma)
		return MessageTypes
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
		return MessageTypes
	case '{':
		l.Next()
		l.Emit(token.OpenCurl)
		return FieldOptions
	default:
		l.AcceptRunUntil(" .<>[]{,\t\n\r")
		l.Emit(token.Type)
		return MessageTypes
	}
}

func Service(l *lexer.Lexer) lexer.State {
	lexer.IgnoreWhiteSpace(l)

	l.AcceptRunUntil(" \t\n\r")
	if l.Current() == "" {
		errorf(l, "expected service name but got nothing")
		return nil
	}
	l.Emit(token.Identifier)

	lexer.IgnoreWhiteSpace(l)
	if l.Peek() != '{' {
		errorf(l, "expected '{' but got %v", l.Peek())
		return nil
	}
	l.Next()
	l.Emit(token.OpenCurl)

	return ServiceMethods
}

func ServiceMethods(l *lexer.Lexer) lexer.State {
	lexer.IgnoreWhiteSpace(l)

	if l.Peek() == '}' {
		l.Next()
		l.Emit(token.CloseCurl)
		return Stmt
	}

	l.AcceptRunUntil("( \t\n\r")
	if l.Current() == "" {
		errorf(l, "expected method name but got nothing")
		return nil
	}
	l.Emit(token.Identifier)

	lexer.IgnoreWhiteSpace(l)

	if l.Peek() != '(' {
		errorf(l, "expected '(' but got %v", l.Peek())
		return nil
	}
	l.Next()
	l.Emit(token.OpenParen)

	return ServiceMethodArgs
}

func ServiceMethodArgs(l *lexer.Lexer) lexer.State {
	lexer.IgnoreWhiteSpace(l)

	if l.Peek() == ')' {
		l.Next()
		l.Emit(token.CloseParen)
		return ServiceMethodReturns
	}

	l.AcceptRunUntil(" \t\n\r")
	if l.Current() == "" {
		errorf(l, "expected argument name but got nothing")
		return nil
	}
	l.Emit(token.Identifier)

	lexer.IgnoreWhiteSpace(l)

	if l.Peek() != ':' {
		errorf(l, "expected ':' but got %v", l.Peek())
		return nil
	}
	l.Next()
	l.Emit(token.Colon)

	return nil
}

func ServiceMethodReturns(l *lexer.Lexer) lexer.State {
	return nil
}
