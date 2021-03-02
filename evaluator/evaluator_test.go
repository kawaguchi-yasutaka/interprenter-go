package evaluator

import (
	"interpreter-go/lexer"
	"interpreter-go/object"
	"interpreter-go/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			input:    "5",
			expected: 5,
		},
		{
			input:    "10",
			expected: 10,
		},
		{
			input:    "-5",
			expected: -5,
		},
		{
			input:    "-10",
			expected: -10,
		},
		{
			input:    "5 + 10",
			expected: 15,
		},
		{
			input:    "10 - 5",
			expected: 5,
		},
		{
			input:    "10 * 5",
			expected: 50,
		},
		{
			input:    "10 / 5",
			expected: 2,
		},
		{
			input:    "10 - 10",
			expected: 0,
		},
		{
			input:    "(10 - (10 - 4) * 1) / 2",
			expected: 3,
		},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testIntegerObject(t, obj, tt.expected)
	}
}

func TestEvaluateBooleanExpresion(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			input:    "true",
			expected: true,
		},
		{
			input:    "false",
			expected: false,
		},
		{
			input:    "1 < 2",
			expected: true,
		},
		{
			input:    "2 < 1",
			expected: false,
		},
		{
			input:    "1 > 2",
			expected: false,
		},
		{
			input:    "1 < 2",
			expected: true,
		},
		{
			input:    "1 == 2",
			expected: false,
		},
		{
			input:    "1 == 1",
			expected: true,
		},
		{
			input:    "1 != 1",
			expected: false,
		},
		{
			input:    "1 != 2",
			expected: true,
		},
		{
			input:    "true == true",
			expected: true,
		},
		{
			input:    "true == false",
			expected: false,
		},
		{
			input:    "true != true",
			expected: false,
		},
		{
			input:    "true != false",
			expected: true,
		},
		{
			input:    "(1 < 2) != false",
			expected: true,
		},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testBooleanObject(t, obj, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			input:    "!true",
			expected: false,
		},
		{
			input:    "!false",
			expected: true,
		},
		{
			input:    "!5",
			expected: false,
		},
		{
			input:    "!!true",
			expected: true,
		},
		{
			input:    "!!false",
			expected: false,
		},
		{
			input:    "!!5",
			expected: true,
		},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testBooleanObject(t, obj, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			input:    "if (true) {10}",
			expected: 10,
		},
		{
			input:    "if (false) {10}",
			expected: nil,
		},
		{
			input:    "if (1) {10}",
			expected: 10,
		},
		{
			input:    "if (1 < 2) {10}",
			expected: 10,
		},
		{
			input:    "if (1 > 2) {10}",
			expected: nil,
		},
		{
			input:    "if (1 < 2) {10} else {20}",
			expected: 10,
		},
		{
			input:    "if (1 > 2) {10} else {20}",
			expected: 20,
		},
	}
	for _, tt := range tests {
		obj := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, obj, int64(integer))
		} else {
			testNullObject(t, obj)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			input:    "return 10",
			expected: 10,
		},
		{
			input:    "return 10;9;",
			expected: 10,
		},
		{
			input:    "9;return 2 * 5;9;",
			expected: 10,
		},
		{
			input:    "if (10 > 1) { if (10 > 1) {return 10;} return 1;}",
			expected: 10,
		},
	}
	for _, tt := range tests {
		obj := testEval(tt.input)
		testIntegerObject(t, obj, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			input:           "5 + true;",
			expectedMessage: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			input:           "5 + true;5;",
			expectedMessage: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			input:           "-true",
			expectedMessage: "unknown operator: -BOOLEAN",
		},
		{
			input:           "true + false",
			expectedMessage: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input:           "5;true + false;5;",
			expectedMessage: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input:           "if(10 > 1) { true + false}",
			expectedMessage: "unknown operator: INTEGER + BOOLEAN",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not integer got=%T(%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object wrong value got=%d,expected=%d", result.Value, expected)
		return false
	}
	return true

}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not boolean got=%T(%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object wrong value got=%t,expected=%t", result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not null got=%T(%+v)", obj, obj)
		return false
	}
	return true
}
