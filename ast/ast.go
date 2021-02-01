package ast

import (
	"bytes"
	"interpreter-go/token"
)

type Node interface {
	String() string
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

func (ls LetStatementNode) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i Identifier) ExpressionNode() {}

func (i Identifier) String() string {
	return i.Value
}

type Program struct {
	Statements []Statement
}

func (p Program) String() string {
	var out bytes.Buffer

	for _, statement := range p.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
}

type ReturnStatementNode struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs ReturnStatementNode) TokenLiteral() string {
	return rs.TokenLiteral()
}

func (rs ReturnStatementNode) StatementNode() {}

func (rs ReturnStatementNode) String() string {
	var out bytes.Buffer

	out.WriteString(rs.Token.Literal + "")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es ExpressionStatement) StatementNode() {}

func (es ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il IntegerLiteral) ExpressionNode() {}

func (il IntegerLiteral) String() string {
	return il.TokenLiteral()
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe PrefixExpression) ExpressionNode() {}

func (pe PrefixExpression) String() string {
	if pe.Right != nil {
		return pe.Operator + pe.Right.String()
	}
	return ""
}
