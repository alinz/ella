package lexer

import "fmt"

type Type int

type Token struct {
	Val   string
	Type  Type
	Start int
	End   int
}

func (t Token) String() string {
	return fmt.Sprintf("Type: %d, Val: %s, Start: %d, End: %d", t.Type, t.Val, t.Start, t.End)
}

type Emitter interface {
	Emit(token *Token)
}

type EmitterFunc func(token *Token)

func (fn EmitterFunc) Emit(token *Token) {
	fn(token)
}

type Iterator interface {
	Next() *Token
}

type IteratorFunc func() *Token

func (fn IteratorFunc) Next() *Token {
	return fn()
}

type EmitterIterator struct {
	tokens chan *Token
}

func (e *EmitterIterator) Emit(token *Token) {
	e.tokens <- token
}

func (e EmitterIterator) Next() *Token {
	value, ok := <-e.tokens
	if !ok {
		return nil
	}

	return value
}

func NewEmitterIterator() *EmitterIterator {
	return &EmitterIterator{
		tokens: make(chan *Token, 2),
	}
}
