package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Lexer lexer struct
type Lexer struct {
	emitter Emitter
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

func (l *Lexer) Emit(tokenType Type) {
	token := &Token{
		Type:  tokenType,
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

func (l *Lexer) PeekRunUntil(invalid string) (value string, reset func()) {
	pos := l.pos
	start := l.start
	width := l.width

	l.AcceptRunUntil(invalid)
	value = l.Current()

	return value, func() {
		l.pos = pos
		l.start = start
		l.width = width
	}
}

func (l *Lexer) Errorf(errType Type, format string, args ...interface{}) {
	l.emitter.Emit(&Token{
		Type:  errType,
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

func New(input string, emitter Emitter) *Lexer {
	return &Lexer{
		input:   input,
		emitter: emitter,
	}
}

func Start(input string, emitter Emitter, inital State) {
	lex := New(input, emitter)
	for state := inital; state != nil; {
		state = state(lex)
	}
}
