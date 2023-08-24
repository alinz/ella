package scanner

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"ella.to/internal/token"
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

func (l *Lexer) Emit(typ token.Type) {
	token := &token.Token{
		Type:    typ,
		Literal: l.input[l.start:l.pos],
		Start:   l.start,
		End:     l.pos,
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

func (l *Lexer) PeekN(n int) string {
	end := l.pos + n
	if end > len(l.input) {
		end = len(l.input)
	}

	return l.input[l.pos:end]
}

func (l *Lexer) NextN(n int) string {
	end := l.pos + n
	if end > len(l.input) {
		end = len(l.input)
	}

	l.pos = end
	return l.input[l.start:end]
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

func (l *Lexer) AcceptRun(valid string) bool {
	result := false
	for strings.ContainsRune(valid, l.Next()) {
		result = true
	}
	l.Backup()
	return result
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
		Type:    token.Error,
		Literal: fmt.Sprintf(format, args...),
		Start:   l.start,
		End:     l.pos,
	})
}

func (l *Lexer) Run(state State) {
	for state != nil {
		state = state(l)
	}
}

type State func(*Lexer) State

func Start(emitter token.Emitter, inital State, input string) {
	lexer := &Lexer{
		emitter: emitter,
		input:   input,
	}
	for state := inital; state != nil; {
		state = state(lexer)
	}
}

func StartWithFilenames(emitter token.Emitter, inital State, filenames ...string) {
	for i, filename := range filenames {
		b, err := os.ReadFile(filename)
		if err != nil {
			emitter.Emit(&token.Token{
				Type:     token.Error,
				Literal:  err.Error(),
				Filename: filename,
			})
			return
		}
		Start(token.EmitterFunc(func(tok *token.Token) {
			if tok.Type == token.EOF && i != len(filenames)-1 {
				// because we havn't process all the files,
				// we simply ignore EOF of previous procceed files
				return
			}

			tok.Filename = filename
			emitter.Emit(tok)

		}), inital, string(b))
	}
}

func ignoreWhiteSpace(l *Lexer) (newLine bool) {
	l.AcceptRun(" \t\r\n")
	if strings.Contains(l.Current(), "\n") {
		newLine = true
	}
	l.Ignore()
	return newLine
}
