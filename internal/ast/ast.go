package ast

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementLiteral()
}
