package ast

import (
	"bytes"
	"interpreter-go/token"
	"strings"
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

//let <Identifier> = <expression>;
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

//return <expression>;
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

//<Expression>;
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

//<prefix operator> <Expression>;
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe PrefixExpression) ExpressionNode() {}

//本の中ではbytes.bufferにwriteStringで書き込んでいる
func (pe PrefixExpression) String() string {
	if pe.Right != nil {
		return "(" + pe.Operator + pe.Right.String() + ")"
	}
	return ""
}

//<expression> <infix operator> <expression>;
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (pe InfixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe InfixExpression) ExpressionNode() {}

//本の中ではbytes.bufferにwriteStringで書き込んでいる
func (pe InfixExpression) String() string {
	if pe.Right != nil {
		return "(" + pe.Left.String() + " " + pe.Operator + " " + pe.Right.String() + ")"
	}
	return ""
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b Boolean) ExpressionNode() {}

func (b Boolean) String() string {
	return b.TokenLiteral()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie IfExpression) ExpressionNode() {}

func (ie IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs BlockStatement) StatementNode() {}

func (bs BlockStatement) String() string {
	var out bytes.Buffer

	for _, statement := range bs.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl FunctionLiteral) ExpressionNode() {}

func (fl FunctionLiteral) String() string {
	var out bytes.Buffer

	parameter := []string{}

	for _, p := range fl.Parameters {
		parameter = append(parameter, p.String())
	}
	out.WriteString(fl.Token.Literal)
	out.WriteString("(")
	out.WriteString(strings.Join(parameter, ","))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}
