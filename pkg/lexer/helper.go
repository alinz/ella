package lexer

import (
	"unicode"
)

func IsAlphaNumeric(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c)
}

func IgnoreWhiteSpace(l *Lexer) {
	l.AcceptRun(" \t\n")
	l.Ignore()
}

func IgnoreSpaceTabs(l *Lexer) {
	l.AcceptRun(" \t")
	l.Ignore()
}
