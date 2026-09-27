package interpreter

import (
	"bytes"
	"strings"
	"testing"

	"glox/pkg/parser"
	"glox/pkg/scanner"
)

func evaluateSource(t *testing.T, source string, out *bytes.Buffer) (any, error) {
	t.Helper()
	if !strings.HasSuffix(strings.TrimSpace(source), ";") {
		source += ";"
	}
	tokens, err := scanner.NewScanner(source).Scan()
	if err != nil {
		t.Fatalf("unexpected scan error on %q: %v", source, err)
	}

	stmts, err := parser.NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("unexpected parse error on %q: %v", source, err)
	}

	interp := NewInterpreter()
	if out != nil {
		interp.SetWriter(out)
	}

	return interp.Interpret(stmts)
}

func mustEvaluate(t *testing.T, source string) any {
	t.Helper()
	val, err := evaluateSource(t, source, nil)
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
		{"statement runtime error", `print 1 + "a";`, "Operands of + must be either numbers or strings"},
		{"reading undefined variable", "print x;", "Undefined variable 'x'."},
		{"assigning undefined variable", "x = 10;", "Undefined variable 'x'."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evaluateSource(t, tt.input, nil)
			if err == nil {
				t.Fatalf("expected runtime error on %q, got none", tt.input)
			}
			rErr, ok := err.(*RuntimeError)
			if !ok {
				t.Fatalf("expected *RuntimeError, got %T: %v", err, err)
			}
			if !strings.Contains(rErr.Error(), tt.errContains) {
				t.Errorf("expected error containing %q, got %q", tt.errContains, rErr.Error())
			}
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
		t.Run(tt.expected, func(t *testing.T) {
			got := Stringify(tt.input)
			if got != tt.expected {
				t.Errorf("Stringify(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestInterpreter_Statements(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		expectedOutput string
	}{
		{
			name:           "print number",
			source:         "print 42;",
			expectedOutput: "42\n",
		},
		{
			name:           "print expression",
			source:         "print 1 + 2;",
			expectedOutput: "3\n",
		},
		{
			name:           "print string",
			source:         `print "hello world";`,
			expectedOutput: "hello world\n",
		},
		{
			name:           "multiple print statements",
			source:         "print 1;\nprint 2;\nprint 3;",
			expectedOutput: "1\n2\n3\n",
		},
		{
			name:           "expression statement produces no print output",
			source:         "100 + 1;",
			expectedOutput: "",
		},
		{
			name:           "var declaration without initializer defaults to nil",
			source:         "var a; print a;",
			expectedOutput: "nil\n",
		},
		{
			name:           "var declaration with initializer",
			source:         "var a = 42; print a;",
			expectedOutput: "42\n",
		},
		{
			name:           "variable assignment and reassignment",
			source:         "var a = 1; a = 2; print a;",
			expectedOutput: "2\n",
		},
		{
			name:           "variables in binary expression",
			source:         "var a = 5; var b = 10; print a + b;",
			expectedOutput: "15\n",
		},
		{
			name:           "assignment expression returns assigned value",
			source:         "var a = 1; print a = 2;",
			expectedOutput: "2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			_, err := evaluateSource(t, tt.source, &buf)
			if err != nil {
				t.Fatalf("unexpected runtime error: %v", err)
			}

			if buf.String() != tt.expectedOutput {
				t.Errorf("got %q, want %q", buf.String(), tt.expectedOutput)
			}
		})
	}
}
