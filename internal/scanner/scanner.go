package scanner

import (
	"strings"

	"ella.to/internal/token"
)

func Lex(l *Lexer) State {
	newLine := ignoreWhiteSpace(l)

	switch l.Peek() {
	case 0:
		l.Emit(token.EOF)
		return nil
	case '=':
		l.Next()
		if l.Peek() == '>' {
			l.Next()
			l.Emit(token.Return)
			return Lex
		}
		l.Emit(token.Assign)
		return Lex
	case ':':
		l.Next()
		l.Emit(token.Colon)
		return Lex
	case ',':
		l.Next()
		l.Emit(token.Comma)
		return Lex
	case '.':
		l.Next()
		if l.Next() != '.' {
			l.Errorf("extend requires 3 consecutive dots")
			return nil
		}
		if l.Next() != '.' {
			l.Errorf("extend requires 3 consecutive dots")
			return nil
		}
		l.Emit(token.Extend)
		return Lex
	case '{':
		l.Next()
		l.Emit(token.OpenCurly)
		return Lex
	case '}':
		l.Next()
		l.Emit(token.CloseCurly)
		return Lex
	case '(':
		l.Next()
		l.Emit(token.OpenParen)
		return Lex
	case ')':
		l.Next()
		l.Emit(token.CloseParen)
		return Lex
	case '<':
		l.Next()
		l.Emit(token.OpenAngle)
		return Lex
	case '>':
		l.Next()
		l.Emit(token.CloseAngle)
		return Lex
	case '[':
		l.Next()
		if l.Peek() != ']' {
			l.Errorf("expect ] to close array")
			return nil
		}
		l.Next()
		l.Emit(token.Array)
		return Lex
	case '#':
		l.Next()
		l.Ignore()
		l.AcceptRunUntil("\n\r")
		if newLine {
			l.Emit(token.TopComment)
		} else {
			l.Emit(token.RightComment)
		}
	case '\'':
		l.Next()
		l.Ignore()
		l.AcceptRunUntil("'\n\r")
		if l.Peek() != '\'' {
			l.Errorf("expect ' to close single quote")
			return nil
		}
		l.Emit(token.ConstStringSingleQuote)
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
		l.Emit(token.ConstStringDoubleQuote)
		l.Next()
		l.Ignore()
	case '`':
		l.Next()
		l.Ignore()
		l.AcceptRunUntil("`")
		if l.Peek() != '`' {
			l.Errorf("expect ` to close back multi line quote")
			return nil
		}
		l.Emit(token.ConstStringBacktickQoute)
		l.Next()
		l.Ignore()
	default:
		ok, found := parseNumber(l)
		if found {
			return Lex
		} else if !ok {
			return nil
		}

		l.AcceptRunUntil("=,.:?{}()<>[]# \t\n\r")
		if l.Current() == "" {
			l.Errorf("expect something but got nothing")
			return nil
		}
		if !reservedKeywrod(l) {
			l.Emit(token.Identifier)
		}
	}

	return Lex
}

func Number(l *Lexer) State {
	parseNumber(l)
	return nil
}

func parseNumber(l *Lexer) (ok bool, found bool) {
	isFloat := false

	l.Accept("+-")

	digits := "0123456789"
	if l.Accept("0") && l.Accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}

	digits += "_"

	l.AcceptRun(digits)

	if len(l.Current()) == 0 || strings.HasPrefix(l.Current(), "_") {
		return true, false // not founding number but no error
	}

	if l.Accept(".") {
		isFloat = true
		if !l.AcceptRun(digits) {
			l.Errorf("expected digit after decimal point")
			return false, false // not founding number and with error
		}
	}

	if l.Accept("eE") {
		l.Accept("+-")
		l.AcceptRun("0123456789_")
	}

	if strings.HasSuffix(l.Current(), "_") {
		l.Errorf("expected digit after each underscore")
		return false, false // not founding number and with error
	}

	l.Accept("i")

	isDuration := false
	isBytes := isBytesTypeNum(l)
	if !isBytes && isDurationTypeNum(l) {
		isDuration = true
	}

	peek := l.Peek()

	if peek == 0 || peek == ' ' || peek == '\t' || peek == '\n' || peek == '\r' || peek == '#' {
		if strings.Contains(l.Current(), "__") {
			l.Errorf("expected digit after each underscore")
			return false, false // not founding number and with error
		}

		if isFloat && isBytes {
			l.Errorf("bytes number can't be presented as float")
			return false, false
		} else if isFloat && isDuration {
			l.Errorf("duration number can't be presented as float")
			return false, false
		} else if isFloat {
			l.Emit(token.ConstFloat)
		} else if isBytes {
			l.Emit(token.ConstBytes)
		} else if isDuration {
			l.Emit(token.ConstDuration)
		} else {
			l.Emit(token.ConstInt)
		}

		return true, true // founding number and no error
	}

	l.Errorf("unexpected character after number: %c", peek)

	return false, false // not founding number and with error
}

// checking if there is any B, KB, MB, GB, TB, PB, EB, ZB, YB
func isBytesTypeNum(l *Lexer) bool {
	isBytes := false
	if l.Accept("b") {
		isBytes = true
	} else {
		value := l.PeekN(2)
		if value == "kb" || // kilobyte
			value == "mb" || // megabyte
			value == "gb" || // gigabyte
			value == "tb" || // terabyte
			value == "pb" || // petabyte
			value == "eb" { // exabyte
			isBytes = true
			l.Next()
			l.Next()
		}
	}
	return isBytes
}

// checking if there is any ms, s, m, h, which represent millisecond, second, minute, hour
func isDurationTypeNum(l *Lexer) bool {
	value := l.PeekN(2)

	if value == "ns" || value == "us" || value == "ms" { // microsecond
		l.Next()
		l.Next()
		return true
	} else {
		return l.Accept("smh")
	}
}

func reservedKeywrod(l *Lexer) bool {
	switch l.Current() {
	case "const":
		l.Emit(token.Const)
		return true
	case "enum":
		l.Emit(token.Enum)
		return true
	case "message":
		l.Emit(token.Message)
		return true
	case "http":
		l.Emit(token.Http)
		return true
	case "rpc":
		l.Emit(token.Rpc)
		return true
	case "service":
		l.Emit(token.Service)
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
	case "string":
		l.Emit(token.String)
		return true
	case "map":
		l.Emit(token.Map)
		return true
	case "any":
		l.Emit(token.Any)
		return true
	case "file":
		l.Emit(token.File)
		return true
	case "stream":
		l.Emit(token.Stream)
		return true
	case "true", "false":
		l.Emit(token.ConstBool)
		return true
	case "null":
		l.Emit(token.ConstNull)
		return true
	default:
		return false
	}
}
