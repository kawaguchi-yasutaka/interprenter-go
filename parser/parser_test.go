package parser

import (
	"fmt"
	"interpreter-go/ast"
	"interpreter-go/lexer"
	"interpreter-go/token"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 12222222;
`
	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParsErrors(t, &parser)

	if program == nil {
		t.Fatal("ParseProgram() return nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statemsns, got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{expectedIdentifier: "x"},
		{expectedIdentifier: "y"},
		{expectedIdentifier: "foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokneLiternal not 'let' gote=%q", s.TokenLiteral())
		return false
	}

	letSmt, ok := s.(*ast.LetStatementNode)
	if !ok {
		t.Errorf("s not *ast.LetStatemeNode got=%T", s)
		return false
	}

	if letSmt.Token.Type != token.LET {
		t.Errorf("letSmt.Token.Type not '%s' got=%s", token.LET, letSmt.Token.Type)
		return false
	}

	if letSmt.Name.Value != name {
		t.Errorf("letSmt.Name.Value not '%s' got=%s", name, letSmt.Name.Value)
		return false
	}

	if letSmt.Name.Token.Type != token.IDENT {
		t.Errorf("letSmt.Name.Token not '%s' got=%s", token.IDENT, letSmt.Name.Token)
		return false
	}

	if letSmt.Name.TokenLiteral() != name {
		t.Errorf("letSmt.Name.TokenLiteral not '%s' got=%s", name, letSmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 1;
return 10;
return 993322;
`

	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParsErrors(t, &parser)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statemsns, got=%d", len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatementNode)
		if !ok {
			t.Errorf("s not *ast.ReturnStatementNode got=%T", returnStatement)
		}
		if returnStatement.Token.Type != token.RETURN {
			t.Errorf("returnStatement.Token.Type not '%s' got=%s", token.RETURN, returnStatement.Token.Type)
		}
	}
}

func checkParsErrors(t *testing.T, parse *Parser) {
	errors := parse.Errors()

	if len(errors) == 0 {
		return
	}

	for _, err := range errors {
		t.Errorf("parse error: %q", err)
	}
	t.FailNow()
}

func TestParseExpressionStatement(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParsErrors(t, &parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statemsns, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program statemens[0] is not ast.ExpressionStatement got=%T", program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier, got=%s", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() not %s got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestParseIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParsErrors(t, &parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statemsns, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program statemens[0] is not ast.ExpressionStatement got=%T", program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.Identifier, got=%s", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %s got=%d", "5", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() not %s got=%s", "5", literal.TokenLiteral())
	}
}

func TestParsePrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{input: "!5", operator: "!", integerValue: 5},
		{input: "-15", operator: "1", integerValue: 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()
		checkParsErrors(t, &parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statemsns, got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program statemens[0] is not ast.ExpressionStatement got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp not *ast.PrefixExpression, got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {

		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}

	}

}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral, got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("literal.Value not %d got=%d", value, integ.Value)
		return false

	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("literal.TokenLiteral() not %d got=%s", value, integ.TokenLiteral())
		return false
	}
	return true
}
