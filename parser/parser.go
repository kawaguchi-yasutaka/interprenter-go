package parser

import (
	"fmt"
	"interpreter-go/ast"
	"interpreter-go/lexer"
	"interpreter-go/token"
)

type Parser struct {
	lexer     *lexer.Lexer
	errors    []string
	curToken  token.Token
	peekToken token.Token
}

func New(lexer lexer.Lexer) Parser {
	p := Parser{
		lexer:  &lexer,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()
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
		return nil
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
	if p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return letSmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatementNode {
	returnSmt := &ast.ReturnStatementNode{}
	returnSmt.Token = p.curToken

	if p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return returnSmt
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
