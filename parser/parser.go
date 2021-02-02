package parser

import (
	"fmt"
	"interpreter-go/ast"
	"interpreter-go/lexer"
	"interpreter-go/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGRATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type Parser struct {
	lexer          *lexer.Lexer
	errors         []string
	curToken       token.Token
	peekToken      token.Token
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(lexer lexer.Lexer) Parser {
	p := Parser{
		lexer:          &lexer,
		errors:         []string{},
		prefixParseFns: map[token.TokenType]prefixParseFn{},
	}
	p.nextToken()
	p.nextToken()
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for {
		if p.peekToken.Type == token.EOF {
			break
		}
		statement := p.ParseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) ParseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatementNode {
	letSmt := &ast.LetStatementNode{}
	letSmt.Token = p.curToken

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	letSmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//式に対応したら変更
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return letSmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatementNode {
	returnSmt := &ast.ReturnStatementNode{}
	returnSmt.Token = p.curToken

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return returnSmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expressionSmt := &ast.ExpressionStatement{}
	expressionSmt.Token = p.curToken

	expressionSmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return expressionSmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn, ok := p.prefixParseFns[p.curToken.Type]
	if !ok {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	return prefixFn()
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	return &ast.IntegerLiteral{Token: p.curToken, Value: value}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)
	return &expression
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next oken to b %s,got %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p Parser) Errors() []string {
	return p.errors
}

func (p Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
