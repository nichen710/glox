package interpreter

import (
	"testing"

	"glox/pkg/parser"
	"glox/pkg/scanner"
)

func evaluateSource(t *testing.T, source string) (any, error) {
	t.Helper()
	tokens, err := scanner.NewScanner(source).Scan()
	if err != nil {
		t.Fatalf("unexpected scan error on %q: %v", source, err)
	}

	expr, err := parser.NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("unexpected parse error on %q: %v", source, err)
	}

	return NewInterpreter().Interpret(expr)
}

func mustEvaluate(t *testing.T, source string) any {
	t.Helper()
	val, err := evaluateSource(t, source)
	if err != nil {
		t.Fatalf("unexpected evaluation error on %q: %v", source, err)
	}
	return val
}

func TestInterpreter_LiteralsAndGrouping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{"number", "42", 42.0},
		{"float", "3.14", 3.14},
		{"string", `"hola mundo"`, "hola mundo"},
		{"true", "true", true},
		{"false", "false", false},
		{"nil", "nil", nil},
		{"grouping", "((10))", 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustEvaluate(t, tt.input)
			if got != tt.expected {
				t.Errorf("got %v (%T), want %v (%T)", got, got, tt.expected, tt.expected)
			}
		})
	}
}

func TestInterpreter_Unary(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{"negation", "-5", -5.0},
		{"double negation", "-(-10)", 10.0},
		{"not true", "!true", false},
		{"not false", "!false", true},
		{"not nil", "!nil", true},
		{"not number is false", "!0", false},
		{"not string is false", `!""`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustEvaluate(t, tt.input)
			if got != tt.expected {
				t.Errorf("got %v (%T), want %v (%T)", got, got, tt.expected, tt.expected)
			}
		})
	}
}

func TestInterpreter_BinaryArithmetic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{"addition", "2 + 3", 5.0},
		{"subtraction", "10 - 4", 6.0},
		{"multiplication", "6 * 7", 42.0},
		{"division", "15 / 3", 5.0},
		{"modulo", "10 % 3", 1.0},
		{"string concat", `"hola " + "mundo"`, "hola mundo"},
		{"precedence", "2 + 3 * 4", 14.0},
		{"parenthesized precedence", "(2 + 3) * 4", 20.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustEvaluate(t, tt.input)
			if got != tt.expected {
				t.Errorf("got %v (%T), want %v (%T)", got, got, tt.expected, tt.expected)
			}
		})
	}
}

func TestInterpreter_ComparisonAndEquality(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"less true", "1 < 2", true},
		{"less false", "2 < 1", false},
		{"less equal true", "2 <= 2", true},
		{"greater true", "5 > 3", true},
		{"greater equal true", "5 >= 5", true},
		{"equality numbers", "42 == 42", true},
		{"equality strings", `"abc" == "abc"`, true},
		{"equality different strings", `"abc" == "def"`, false},
		{"equality different types", `"42" == 42`, false},
		{"equality nil", "nil == nil", true},
		{"equality nil with false", "nil == false", false},
		{"inequality", "1 != 2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustEvaluate(t, tt.input)
			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestInterpreter_RuntimeErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		errContains string
	}{
		{"negation on string", `-"hola"`, "Operand of - must be a number"},
		{"addition string and number", `"hola" + 5`, "Operands of + must be either numbers or strings"},
		{"subtraction string and number", `"hola" - 5`, "Operands of - must be numbers"},
		{"multiplication on boolean", "true * 2", "Operands of * must be numbers"},
		{"division by zero", "5 / 0", "Division by 0 is not allowed"},
		{"modulo by zero", "10 % 0", "Modulo by 0 is not allowed"},
		{"relational on strings", `"a" < "b"`, "Operands of < must be numbers"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evaluateSource(t, tt.input)
			if err == nil {
				t.Fatalf("expected runtime error on %q, got none", tt.input)
			}
			rErr, ok := err.(*RuntimeError)
			if !ok {
				t.Fatalf("expected *RuntimeError, got %T: %v", err, err)
			}
			t.Logf("Got expected runtime error: %v", rErr)
		})
	}
}

func TestStringify(t *testing.T) {
	tests := []struct {
		input    any
		expected string
	}{
		{nil, "nil"},
		{true, "true"},
		{false, "false"},
		{42.0, "42"},
		{3.14, "3.14"},
		{"texto", "texto"},
	}

	for _, tt := range tests {
		got := Stringify(tt.input)
		if got != tt.expected {
			t.Errorf("Stringify(%v) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
