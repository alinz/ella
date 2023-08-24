package token

type Type int

type Token struct {
	Filename string
	Literal  string
	Type     Type
	Start    int
	End      int
}

type Emitter interface {
	Emit(token *Token)
}

type EmitterFunc func(token *Token)

func (fn EmitterFunc) Emit(token *Token) {
	fn(token)
}

type Iterator interface {
	NextToken() *Token
}

type EmitterIterator struct {
	tokens chan *Token
	end    *Token
}

var _ Emitter = (*EmitterIterator)(nil)
var _ Iterator = (*EmitterIterator)(nil)

func (e *EmitterIterator) Emit(token *Token) {
	e.tokens <- token
}

func (e *EmitterIterator) NextToken() *Token {
	tok, ok := <-e.tokens
	if !ok {
		return e.end
	} else if tok.Type == EOF {
		e.end = tok
		close(e.tokens)
		e.tokens = nil
	}

	return tok
}

func NewEmitterIterator() *EmitterIterator {
	return &EmitterIterator{
		tokens: make(chan *Token, 2),
	}
}

func OneOfTypes(tok *Token, types ...Type) bool {
	for _, t := range types {
		if tok.Type == t {
			return true
		}
	}

	return false
}
