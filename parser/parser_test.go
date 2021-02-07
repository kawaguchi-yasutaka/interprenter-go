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
	checkParsErrors(t, parser)

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
	checkParsErrors(t, parser)

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

func TestParseExpressionStatement(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParsErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statemsns, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program statemens[0] is not ast.ExpressionStatement got=%T", program.Statements[0])
	}
	if !testIdentifierLiteral(t, stmt.Expression, "foobar") {
		return
	}
}

func TestParseIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParsErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statemsns, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program statemens[0] is not ast.ExpressionStatement got=%T", program.Statements[0])
	}

	if !testIntegerLiteral(t, stmt.Expression, 5) {
		return
	}
}

func TestParseBooleanExpresion(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParsErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statemsns, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program statemens[0] is not ast.ExpressionStatement got=%T", program.Statements[0])
	}

	if !testBooleanExpression(t, stmt.Expression, true) {
		return
	}
}

func TestParsePrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{input: "!5", operator: "!", integerValue: 5},
		{input: "-15", operator: "-", integerValue: 15},
		{input: "!true", operator: "!", integerValue: true},
		{input: "!false", operator: "!", integerValue: false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()
		checkParsErrors(t, parser)

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
			fmt.Errorf("exp.Operator is not %s, got=%s", exp.Operator, tt.operator)
		}
		if !testLiteralExpresion(t, exp.Right, tt.integerValue) {
			return
		}
	}

}

func TestParseInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{input: "5 + 5", leftValue: 5, operator: "+", rightValue: 5},
		{input: "5 - 5", leftValue: 5, operator: "-", rightValue: 5},
		{input: "5 * 5", leftValue: 5, operator: "*", rightValue: 5},
		{input: "5 / 5", leftValue: 5, operator: "/", rightValue: 5},
		{input: "5 > 5", leftValue: 5, operator: ">", rightValue: 5},
		{input: "5 < 5", leftValue: 5, operator: "<", rightValue: 5},
		{input: "5 == 5", leftValue: 5, operator: "==", rightValue: 5},
		{input: "5 != 5", leftValue: 5, operator: "!=", rightValue: 5},
		{input: "true == true", leftValue: true, operator: "==", rightValue: true},
		{input: "true != false", leftValue: true, operator: "!=", rightValue: false},
		{input: "false == false", leftValue: false, operator: "==", rightValue: false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()
		checkParsErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statemsns, got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program statemens[0] is not ast.ExpressionStatement got=%T", program.Statements[0])
		}
		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := "if(x < y) { x }"

	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParsErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statemsns, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program statemens[0] is not ast.ExpressionStatement got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfEpression got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	cStmt, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement got=%T", program.Statements[0])
	}
	sExp, ok := cStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("cStmt.Expression is not ast.Identifier got=%T", program.Statements[0])
	}
	if !testIdentifierLiteral(t, sExp, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternativen is not nil got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if(x < y) { x } else { y }"

	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParsErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statemsns, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program statemens[0] is not ast.ExpressionStatement got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfEpression got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	cStmt, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement got=%T", exp.Consequence.Statements[0])
	}
	sExp, ok := cStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("cStmt.Expression is not ast.Identifier got=%T", cStmt.Expression)
	}
	if !testIdentifierLiteral(t, sExp, "x") {
		return
	}

	aStmt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement got=%T", exp.Alternative.Statements[0])
	}
	aExp, ok := aStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("cStmt.Expression is not ast.Identifier got=%T", cStmt.Expression)
	}
	if !testIdentifierLiteral(t, aExp, "y") {
		return
	}

}

func TestFunctionLiteral(t *testing.T) {
	input := "fn(x, y) {x + y};"
	l := lexer.New(input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParsErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statemsns, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program statemens[0] is not ast.ExpressionStatement got=%T", program.Statements[0])
	}
	funcLiteral, ok := stmt.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not *ast.FunctionLiteral got=%T", stmt.Expression)
	}

	if len(funcLiteral.Parameters) != 2 {
		t.Fatalf("exp.Parameters does not contain 2 statemsns, got=%d", len(funcLiteral.Parameters))
	}

	if !testIdentifierLiteral(t, funcLiteral.Parameters[0], "x") {
		return
	}

	if !testIdentifierLiteral(t, funcLiteral.Parameters[1], "y") {
		return
	}

	if len(funcLiteral.Body.Statements) != 1 {
		t.Fatalf("funcLiteral.Body.Statements does not contain 2 statemsns, got=%d", len(funcLiteral.Parameters))
	}

	blockStatement, ok := funcLiteral.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(" exp.Body.Statements[0] is not *ast.ExpressionStatement got=%T", funcLiteral.Body.Statements[0])
	}

	if !testInfixExpression(t, blockStatement.Expression, "x", "+", "y") {
		return
	}

}

func TestFunctionParameters(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "fn() { 1 };",
			expected: []string{},
		},
		{
			input:    "fn(x) {x};",
			expected: []string{"x"},
		},
		{
			input:    "fn(x,y) {x + y};",
			expected: []string{"x", "y"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()
		checkParsErrors(t, parser)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program statemens[0] is not ast.ExpressionStatement got=%T", program.Statements[0])
		}
		funcLiteral, ok := stmt.Expression.(*ast.FunctionLiteral)

		if !ok {
			t.Fatalf("stmt.Expression is not *ast.FunctionLiteral got=%T", stmt.Expression)
		}

		if len(funcLiteral.Parameters) != len(tt.expected) {
			t.Fatalf("funcLiteral.Parameters does not contain %d statemsns, got=%d", len(tt.expected), len(funcLiteral.Parameters))
		}

		for i, s := range tt.expected {
			testLiteralExpresion(t, funcLiteral.Parameters[i], s)
		}

	}

}
func TestOperatorPrecedencesParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "-a + b",
			expected: "((-a) + b)",
		},
		{
			input:    "!-a",
			expected: "(!(-a))",
		},
		{
			input:    "a + b + c",
			expected: "((a + b) + c)",
		},
		{
			input:    "a + b - c",
			expected: "((a + b) - c)",
		},
		{
			input:    "a * b * c",
			expected: "((a * b) * c)",
		},
		{
			input:    "a * b / c",
			expected: "((a * b) / c)",
		},
		{
			input:    "a + b / c",
			expected: "(a + (b / c))",
		},
		{
			input:    "3 + 4;-5 * 5;",
			expected: "(3 + 4)((-5) * 5)",
		},
		{
			input:    "5 > 4 == 3 < 4",
			expected: "((5 > 4) == (3 < 4))",
		},
		{
			input:    "5 < 4 != 3 > 4",
			expected: "((5 < 4) != (3 > 4))",
		},
		{
			input:    "3 + 4 * 5 == 3 * 1 + 4 * 5",
			expected: "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			input:    "true;",
			expected: "true",
		},
		{
			input:    "false;",
			expected: "false",
		},
		{
			input:    "3 > 5 == false",
			expected: "((3 > 5) == false)",
		},
		{
			input:    "3 > 5 == true",
			expected: "((3 > 5) == true)",
		},
		{
			input:    "1 + (2 + 3) + 4",
			expected: "((1 + (2 + 3)) + 4)",
		},
		{
			input:    "(5 + 5) * 2",
			expected: "((5 + 5) * 2)",
		},
		{
			input:    "-(5 + 5)",
			expected: "(-(5 + 5))",
		},
		{
			input:    "!(true == true)",
			expected: "(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()
		checkParsErrors(t, parser)

		if program.String() != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, program.String())
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

func testIdentifierLiteral(t *testing.T, il ast.Expression, value string) bool {
	identifier, ok := il.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier, got=%T", il)
		return false
	}
	if identifier.Value != value {
		t.Errorf("identifier.Value not %s got=%s", value, identifier.Value)
		return false
	}
	if identifier.TokenLiteral() != value {
		t.Errorf("identifier.TokenLiteral not %s got=%s", value, identifier.TokenLiteral())
		return false
	}
	return true
}

func testBooleanExpression(t *testing.T, il ast.Expression, value bool) bool {
	boolean, ok := il.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp not *ast.Identifier, got=%T", il)
		return false
	}
	if boolean.Value != value {
		t.Errorf("boolena.Value not %t got=%t", value, boolean.Value)
		return false
	}
	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolena.TokenLiteral not %s got=%s", fmt.Sprintf("%t", value), boolean.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpresion(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, exp, int64(v))
	case int64:
		testIntegerLiteral(t, exp, v)
	case string:
		testIdentifierLiteral(t, exp, v)
	case bool:
		testBooleanExpression(t, exp, v)
	default:
		t.Errorf("type of exp not handled got=%T", exp)
	}
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp not *ast.InfixExpression, got=%T", exp)
		return false
	}
	if !testLiteralExpresion(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Fatalf("exp.Operator is not %q got=%q", operator, opExp.Operator)
	}
	if !testLiteralExpresion(t, opExp.Right, right) {
		return false
	}
	return true
}
