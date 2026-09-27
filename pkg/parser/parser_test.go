package parser

import (
	"reflect"
	"strings"
	"testing"

	"glox/pkg/expression"
	"glox/pkg/scanner"
	"glox/pkg/statement"
	"glox/pkg/token"
)

func mustParse(t *testing.T, source string) expression.Expression {
	t.Helper()
	if !strings.HasSuffix(strings.TrimSpace(source), ";") {
		source += ";"
	}
	tokens, err := scanner.NewScanner(source).Scan()
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
	stmts, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
	exprStmt, ok := stmts[0].(statement.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", stmts[0])
	}
	return exprStmt.Expression
}

func TestParser_Literal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected expression.Expression
	}{
		{
			name:     "number literal",
			input:    "101",
			expected: expression.Literal{Value: 101.0},
		},
		{
			name:     "string literal",
			input:    "\"hola mundo\"",
			expected: expression.Literal{Value: "hola mundo"},
		},
		{
			name:     "boolean literal true",
			input:    "true",
			expected: expression.Literal{Value: true},
		},
		{
			name:     "boolean literal false",
			input:    "false",
			expected: expression.Literal{Value: false},
		},
		{
			name:     "nil literal",
			input:    "nil",
			expected: expression.Literal{Value: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustParse(t, tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("got %#v, want %#v", got, tt.expected)
			}
		})
	}
}

func TestParser_Unary(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected expression.Expression
	}{
		{
			name:  "unary negation",
			input: "-101",
			expected: expression.Unary{
				Operator: token.Token{Type: token.MINUS, Lexeme: "-", Line: 1},
				Right:    expression.Literal{Value: 101.0},
			},
		},
		{
			name:  "unary logical not",
			input: "!true",
			expected: expression.Unary{
				Operator: token.Token{Type: token.BANG, Lexeme: "!", Line: 1},
				Right:    expression.Literal{Value: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustParse(t, tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("got %#v, want %#v", got, tt.expected)
			}
		})
	}
}

func TestParser_Binary(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected expression.Expression
	}{
		{
			name:  "addition",
			input: "100 + 1",
			expected: expression.Binary{
				Left:     expression.Literal{Value: 100.0},
				Operator: token.Token{Type: token.PLUS, Lexeme: "+", Line: 1},
				Right:    expression.Literal{Value: 1.0},
			},
		},
		{
			name:  "subtraction",
			input: "102 - 2",
			expected: expression.Binary{
				Left:     expression.Literal{Value: 102.0},
				Operator: token.Token{Type: token.MINUS, Lexeme: "-", Line: 1},
				Right:    expression.Literal{Value: 2.0},
			},
		},
		{
			name:  "multiplication",
			input: "6 * 7",
			expected: expression.Binary{
				Left:     expression.Literal{Value: 6.0},
				Operator: token.Token{Type: token.STAR, Lexeme: "*", Line: 1},
				Right:    expression.Literal{Value: 7.0},
			},
		},
		{
			name:  "division",
			input: "64 / 8",
			expected: expression.Binary{
				Left:     expression.Literal{Value: 64.0},
				Operator: token.Token{Type: token.SLASH, Lexeme: "/", Line: 1},
				Right:    expression.Literal{Value: 8.0},
			},
		},
		{
			name:  "modulo",
			input: "10 % 3",
			expected: expression.Binary{
				Left:     expression.Literal{Value: 10.0},
				Operator: token.Token{Type: token.PERCENT, Lexeme: "%", Line: 1},
				Right:    expression.Literal{Value: 3.0},
			},
		},
		{
			name:  "equal equal",
			input: "true == false",
			expected: expression.Binary{
				Left:     expression.Literal{Value: true},
				Operator: token.Token{Type: token.EQUAL_EQUAL, Lexeme: "==", Line: 1},
				Right:    expression.Literal{Value: false},
			},
		},
		{
			name:  "bang equal",
			input: "1 != 2",
			expected: expression.Binary{
				Left:     expression.Literal{Value: 1.0},
				Operator: token.Token{Type: token.BANG_EQUAL, Lexeme: "!=", Line: 1},
				Right:    expression.Literal{Value: 2.0},
			},
		},
		{
			name:  "less and greater",
			input: "1 < 2",
			expected: expression.Binary{
				Left:     expression.Literal{Value: 1.0},
				Operator: token.Token{Type: token.LESS, Lexeme: "<", Line: 1},
				Right:    expression.Literal{Value: 2.0},
			},
		},
		{
			name:  "less equal and greater equal",
			input: "2 >= 1",
			expected: expression.Binary{
				Left:     expression.Literal{Value: 2.0},
				Operator: token.Token{Type: token.GREATER_EQUAL, Lexeme: ">=", Line: 1},
				Right:    expression.Literal{Value: 1.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustParse(t, tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("got %#v, want %#v", got, tt.expected)
			}
		})
	}
}

func TestParser_Grouping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected expression.Expression
	}{
		{
			name:  "grouping with parentheses",
			input: "(100 + 1)",
			expected: expression.Grouping{
				Expression: expression.Binary{
					Left:     expression.Literal{Value: 100.0},
					Operator: token.Token{Type: token.PLUS, Lexeme: "+", Line: 1},
					Right:    expression.Literal{Value: 1.0},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustParse(t, tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("got %#v, want %#v", got, tt.expected)
			}
		})
	}
}

func TestParser_Precedence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected expression.Expression
	}{
		{
			name:  "multiplication has higher precedence than addition",
			input: "2 + 3 * 4",
			expected: expression.Binary{
				Left:     expression.Literal{Value: 2.0},
				Operator: token.Token{Type: token.PLUS, Lexeme: "+", Line: 1},
				Right: expression.Binary{
					Left:     expression.Literal{Value: 3.0},
					Operator: token.Token{Type: token.STAR, Lexeme: "*", Line: 1},
					Right:    expression.Literal{Value: 4.0},
				},
			},
		},
		{
			name:  "subtraction is left associative",
			input: "1 - 2 - 3",
			expected: expression.Binary{
				Left: expression.Binary{
					Left:     expression.Literal{Value: 1.0},
					Operator: token.Token{Type: token.MINUS, Lexeme: "-", Line: 1},
					Right:    expression.Literal{Value: 2.0},
				},
				Operator: token.Token{Type: token.MINUS, Lexeme: "-", Line: 1},
				Right:    expression.Literal{Value: 3.0},
			},
		},
		{
			name:  "parentheses override precedence",
			input: "(1 + 2) * 3",
			expected: expression.Binary{
				Left: expression.Grouping{
					Expression: expression.Binary{
						Left:     expression.Literal{Value: 1.0},
						Operator: token.Token{Type: token.PLUS, Lexeme: "+", Line: 1},
						Right:    expression.Literal{Value: 2.0},
					},
				},
				Operator: token.Token{Type: token.STAR, Lexeme: "*", Line: 1},
				Right:    expression.Literal{Value: 3.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustParse(t, tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("got %#v, want %#v", got, tt.expected)
			}
		})
	}
}

func TestParser_ParseErrors(t *testing.T) {
	errorCases := []struct {
		name  string
		input string
	}{
		{
			name:  "only semicolon without expression",
			input: ";",
		},
		{
			name:  "missing right parenthesis",
			input: "(100 + 1;",
		},
		{
			name:  "unexpected token",
			input: "100 + * 1;",
		},
		{
			name:  "missing semicolon after print",
			input: "print 42",
		},
		{
			name:  "missing semicolon after expression",
			input: "42 + 1",
		},
	}

	for _, tt := range errorCases {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := scanner.NewScanner(tt.input).Scan()
			if err != nil {
				t.Fatalf("unexpected scan error: %v", err)
			}
			parser := NewParser(tokens)
			stmts, err := parser.Parse()
			if err == nil {
				t.Fatalf("expected parse error for %q, but got none (stmts: %v)", tt.input, stmts)
			}
		})
	}
}

func mustParseStatements(t *testing.T, source string) []statement.Statement {
	t.Helper()
	tokens, err := scanner.NewScanner(source).Scan()
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
	stmts, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	return stmts
}

func TestParser_Statements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []statement.Statement
	}{
		{
			name:  "print statement",
			input: "print 42;",
			expected: []statement.Statement{
				statement.Print{
					Expression: expression.Literal{Value: 42.0},
				},
			},
		},
		{
			name:  "print statement with expression",
			input: "print 1 + 2;",
			expected: []statement.Statement{
				statement.Print{
					Expression: expression.Binary{
						Left:     expression.Literal{Value: 1.0},
						Operator: token.Token{Type: token.PLUS, Lexeme: "+", Line: 1},
						Right:    expression.Literal{Value: 2.0},
					},
				},
			},
		},
		{
			name:  "expression statement",
			input: "100 + 1;",
			expected: []statement.Statement{
				statement.ExpressionStatement{
					Expression: expression.Binary{
						Left:     expression.Literal{Value: 100.0},
						Operator: token.Token{Type: token.PLUS, Lexeme: "+", Line: 1},
						Right:    expression.Literal{Value: 1.0},
					},
				},
			},
		},
		{
			name:  "multiple statements",
			input: "print 1; 2; print 3;",
			expected: []statement.Statement{
				statement.Print{
					Expression: expression.Literal{Value: 1.0},
				},
				statement.ExpressionStatement{
					Expression: expression.Literal{Value: 2.0},
				},
				statement.Print{
					Expression: expression.Literal{Value: 3.0},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustParseStatements(t, tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("got %#v, want %#v", got, tt.expected)
			}
		})
	}
}
