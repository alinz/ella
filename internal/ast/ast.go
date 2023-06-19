package ast

type Node interface {
	nodeLiteral()
	String() string
}
