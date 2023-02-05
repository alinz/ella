package lexer

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

	Counter int // this is a hack to get the argument comma position
	// detected in the lexer, for example (a: map<string, string>, b:int)
	// there is no way for us to detect the comma position whether it's inside
	// map or ouside. This counter increment when we see < and decrement when we
	// see >. When we see , and counter is 0, we know that it's a comma outside
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
	fmt.Println(token)
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

// PeekN is a function that returns the next runes in the input
// based on given number without advancing the position
// returns two values, the runes and the number of iteration <= n
func (l *Lexer) PeekN(n int) (string, int) {
	sb := strings.Builder{}

	total := 0
	i := 0

	for i < n {
		value := l.Next()
		if value == 0 {
			break
		}

		total += l.width
		sb.WriteRune(value)
		i++
	}

	l.pos -= total

	return sb.String(), i
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

func (l *Lexer) PeekRunUntil(invalid string) (value string) {
	pos := l.pos
	start := l.start
	width := l.width

	l.AcceptRunUntil(invalid)
	value = l.Current()

	l.pos = pos
	l.start = start
	l.width = width

	return value
}

func (l *Lexer) Errorf(format string, args ...interface{}) {
	l.emitter.Emit(&token.Token{
		Kind:  token.Error,
		Val:   fmt.Sprintf(format, args...),
		Start: l.start,
		End:   l.pos,
	})
}

func (l *Lexer) Run(state StateFn) {
	for state != nil {
		state = state(l)
	}
}

type StateFn func(*Lexer) StateFn

func Start(input string, emitter token.Emitter, inital StateFn) {
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

func IgnoreSpaceTabs(l *Lexer) {
	l.AcceptRun(" \t")
	l.Ignore()
}
