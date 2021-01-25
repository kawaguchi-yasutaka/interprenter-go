package ast

import "interpreter-go/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	StatementNode()
}

type Expression interface {
	Node
	ExpressionNode()
}

type LetStatementNode struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls LetStatementNode) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls LetStatementNode) StatementNode() {}

type Identifier struct {
	Token token.Token
	Value string
}

func (i Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i Identifier) ExpressionNode() {}

type Program struct {
	Statements []Statement
}
