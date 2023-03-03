package scanner

import (
	"unicode"

	"github.com/alinz/rpc.go/schema/token"
)

func Lex(l *Lexer) State {
	for {
		IgnoreWhiteSpace(l)

		switch l.Peek() {
		case 0:
			l.Emit(token.EOF)
			return nil
		case '=':
			l.Next()
			l.Emit(token.Assign)
		case ':':
			l.Next()
			l.Emit(token.Colon)
		case ',':
			l.Next()
			l.Emit(token.Comma)
		case '.':
			l.Next()
			l.Emit(token.Dot)
		case '{':
			l.Next()
			l.Emit(token.OpenCurly)
		case '}':
			l.Next()
			l.Emit(token.CloseCurly)
		case '(':
			l.Next()
			l.Emit(token.OpenParen)
		case ')':
			l.Next()
			l.Emit(token.CloseParen)
		case '<':
			l.Next()
			l.Emit(token.OpenAngle)
		case '>':
			l.Next()
			l.Emit(token.CloseAngle)
		case '[':
			l.Next()
			l.Emit(token.OpenSquare)
		case ']':
			l.Next()
			l.Emit(token.CloseSquare)
		case '#':
			l.AcceptRunUntil("\n\r")
			l.Emit(token.Comment)
		case '\'':
			l.Next()
			l.Ignore()
			l.AcceptRunUntil("'\n\r")
			if l.Peek() != '\'' {
				l.Errorf("expect ' to close single quote")
				return nil
			}
			l.Emit(token.ConstantString)
			l.Next()
			l.Ignore()
		case '"':
			l.Next()
			l.Ignore()
			l.AcceptRunUntil("\"\n\r")
			if l.Peek() != '"' {
				l.Errorf("expect \" to close double quote")
				return nil
			}
			l.Emit(token.ConstantString)
			l.Next()
			l.Ignore()
		default:
			if lexNumber(l) {
				continue
			}

			l.AcceptRunUntil("=,.:?{}()<>[]# \t\n\r")
			if l.Current() == "" {
				l.Errorf("expect something but got nothing")
				return nil
			}
			if !reservedKeywrod(l) {
				l.Emit(token.Word)
			}
		}
	}
}

func lexNumber(l *Lexer) bool {
	l.Accept("+-")
	digits := "0123456789"
	if l.Accept("0") && l.Accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.AcceptRun(digits)
	if l.Accept(".") {
		l.AcceptRun(digits)
	}
	if l.Accept("eE") {
		l.Accept("+-")
		l.AcceptRun("0123456789")
	}
	l.Accept("i")

	peek := l.Peek()

	if peek == 0 || peek == ' ' || peek == '\t' || peek == '\n' || peek == '\r' || peek == '#' {
		l.Emit(token.ConstantNumber)
		return true
	}
	return false
}

func reservedKeywrod(l *Lexer) bool {
	switch l.Current() {
	case "enum":
		l.Emit(token.Enum)
		return true
	case "message":
		l.Emit(token.Message)
		return true
	case "service":
		l.Emit(token.Service)
		return true
	case "stream":
		l.Emit(token.Stream)
		return true
	case "map":
		l.Emit(token.Map)
		return true
	case "string":
		l.Emit(token.String)
		return true
	case "byte":
		l.Emit(token.Byte)
		return true
	case "bool":
		l.Emit(token.Bool)
		return true
	case "int8":
		l.Emit(token.Int8)
		return true
	case "int16":
		l.Emit(token.Int16)
		return true
	case "int32":
		l.Emit(token.Int32)
		return true
	case "int64":
		l.Emit(token.Int64)
		return true
	case "uint8":
		l.Emit(token.Uint8)
		return true
	case "uint16":
		l.Emit(token.Uint16)
		return true
	case "uint32":
		l.Emit(token.Uint32)
		return true
	case "uint64":
		l.Emit(token.Uint64)
		return true
	case "float32":
		l.Emit(token.Float32)
		return true
	case "float64":
		l.Emit(token.Float64)
		return true
	case "timestamp":
		l.Emit(token.Timestamp)
		return true
	case "any":
		l.Emit(token.Any)
		return true
	default:
		return false
	}
}

func isAlphaNumeric(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c)
}
