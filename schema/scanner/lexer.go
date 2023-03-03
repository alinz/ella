package scanner

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/alinz/rpc.go/schema/token"
)

// Lexer lexer struct
type Lexer struct {
	emitter token.Emitter
	input   string
	start   int
	pos     int
	width   int
}

func (l *Lexer) Current() string {
	return l.input[l.start:l.pos]
}

func (l *Lexer) Emit(kind token.Kind) {
	token := &token.Token{
		Kind:  kind,
		Val:   l.input[l.start:l.pos],
		Start: l.start,
		End:   l.pos,
	}
	l.emitter.Emit(token)
	l.start = l.pos
}

func (l *Lexer) Next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return 0
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Backup()
	return r
}

func (l *Lexer) Backup() {
	l.pos -= l.width
	l.width = 0
}

func (l *Lexer) Ignore() {
	l.start = l.pos
}

func (l *Lexer) Accept(valid string) bool {
	if strings.ContainsRune(valid, l.Next()) {
		return true
	}
	l.Backup()
	return false
}

func (l *Lexer) AcceptRun(valid string) {
	for strings.ContainsRune(valid, l.Next()) {
	}
	l.Backup()
}

func (l *Lexer) AcceptRunUntil(invalid string) {
	for {
		next := l.Next()
		if next == 0 || strings.ContainsRune(invalid, next) {
			break
		}
	}
	l.Backup()
}

func (l *Lexer) Errorf(format string, args ...interface{}) {
	l.emitter.Emit(&token.Token{
		Kind:  token.Error,
		Val:   fmt.Sprintf(format, args...),
		Start: l.start,
		End:   l.pos,
	})
}

func (l *Lexer) Run(state State) {
	for state != nil {
		state = state(l)
	}
}

type State func(*Lexer) State

func Start(input string, emitter token.Emitter, inital State) {
	lexer := &Lexer{
		emitter: emitter,
		input:   input,
	}
	for state := inital; state != nil; {
		state = state(lexer)
	}
}

func IgnoreWhiteSpace(l *Lexer) {
	l.AcceptRun(" \t\n")
	l.Ignore()
}
