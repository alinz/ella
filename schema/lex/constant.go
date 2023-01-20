package lex

import (
	"github.com/alinz/rpc.go/pkg/lexer"
	"github.com/alinz/rpc.go/schema/lex/token"
)

func Constant(next lexer.State) lexer.State {
	return func(l *lexer.Lexer) lexer.State {
		lexer.IgnoreWhiteSpace(l)

		l.AcceptRunUntil(" =\t")
		if l.Current() == "" {
			errorf(l, "expected name for constant but got nothing")
			return nil
		}
		l.Emit(token.Identifier)

		lexer.IgnoreSpaceTabs(l)

		if value := l.Next(); value != '=' {
			errorf(l, "expected '=' but got %s", string(value))
			return nil
		}
		l.Emit(token.Assign)

		lexer.IgnoreSpaceTabs(l)

		l.AcceptRunUntil(" \r\n")
		l.Emit(token.Value)

		return next
	}
}
