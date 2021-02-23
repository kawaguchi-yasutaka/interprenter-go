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
