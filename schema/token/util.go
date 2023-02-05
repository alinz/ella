package token

type Emitter interface {
	Emit(token *Token)
}

type Iterator interface {
	NextToken() *Token
}

type EmitterFunc func(token *Token)

func (fn EmitterFunc) Emit(token *Token) {
	fn(token)
}

type EmitterIterator struct {
	tokens chan *Token
}

var _ Emitter = (*EmitterIterator)(nil)
var _ Iterator = (*EmitterIterator)(nil)

func (e *EmitterIterator) Emit(token *Token) {
	e.tokens <- token
}

func (e EmitterIterator) NextToken() *Token {
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
