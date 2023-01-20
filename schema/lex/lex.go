package lex

import (
	"github.com/alinz/rpc.go/pkg/lexer"
	"github.com/alinz/rpc.go/schema/lex/token"
)

func errorf(l *lexer.Lexer, format string, args ...interface{}) {
	l.Errorf(token.Error, format, args...)
}
