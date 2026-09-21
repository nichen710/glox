package scanner

import (
	"testing"

	"glox/pkg/token"
)

func TestTokenScanner(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token.TokenType
	}{
		{
			name:  "grouping and punctuation",
			input: "(){},;",
			expected: []token.TokenType{
				token.LEFT_PAREN,
				token.RIGHT_PAREN,
				token.LEFT_BRACE,
				token.RIGHT_BRACE,
				token.COMMA,
				token.SEMICOLON,
				token.EOF,
			},
		},
		{
			name:  "arithmetic operators",
			input: "+ - * / %",
			expected: []token.TokenType{
				token.PLUS,
				token.MINUS,
				token.STAR,
				token.SLASH,
				token.PERCENT,
				token.EOF,
			},
		},
		{
			name:  "comparison and equality",
			input: "! != = == < <= > >=",
			expected: []token.TokenType{
				token.BANG,
				token.BANG_EQUAL,
				token.EQUAL,
				token.EQUAL_EQUAL,
				token.LESS,
				token.LESS_EQUAL,
				token.GREATER,
				token.GREATER_EQUAL,
				token.EOF,
			},
		},
		{
			name:     "spaces tabs and carriage returns",
			input:    "   \t\r  ",
			expected: []token.TokenType{token.EOF},
		},
		{
			name:     "single-line comments",
			input:    "// este es un comentario que debe ignorarse\n",
			expected: []token.TokenType{token.EOF},
		},
		{
			name: "tokens with comments and newlines",
			input: `
				// Declarar x
				var x = 10; // asignación
				// fin
			`,
			expected: []token.TokenType{
				token.VAR,
				token.IDENTIFIER,
				token.EQUAL,
				token.NUMBER,
				token.SEMICOLON,
				token.EOF,
			},
		},
		{
			name:  "strings",
			input: `"hola mundo" 'chau mundo'`,
			expected: []token.TokenType{
				token.STRING,
				token.STRING,
				token.EOF,
			},
		},
		{
			name:  "numbers",
			input: "42 3.14 0.5",
			expected: []token.TokenType{
				token.NUMBER,
				token.NUMBER,
				token.NUMBER,
				token.EOF,
			},
		},
		{
			name:  "keywords vs identifiers",
			input: "var variable fun function if iff true trueman while",
			expected: []token.TokenType{
				token.VAR,
				token.IDENTIFIER,
				token.FUN,
				token.IDENTIFIER,
				token.IF,
				token.IDENTIFIER,
				token.TRUE,
				token.IDENTIFIER,
				token.WHILE,
				token.EOF,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScanner(tt.input)
			tokens, err := s.Scan()
			if err != nil {
				t.Fatalf("unexpected error scanning %q: %v", tt.input, err)
			}

			assertTokenTypes(t, tokens, tt.expected)
		})
	}
}

func TestCompleteStatements(t *testing.T) {
	input := `
		fun add(a, b) {
			return a + b;
		}
		var result = add(5, 10);
		print result;
	`
	s := NewScanner(input)
	tokens, err := s.Scan()
	if err != nil {
		t.Fatalf("unexpected error scanning code: %v", err)
	}

	expected := []token.TokenType{
		// fun add(a, b) {
		token.FUN, token.IDENTIFIER, token.LEFT_PAREN, token.IDENTIFIER, token.COMMA, token.IDENTIFIER, token.RIGHT_PAREN, token.LEFT_BRACE,
		// return a + b;
		token.RETURN, token.IDENTIFIER, token.PLUS, token.IDENTIFIER, token.SEMICOLON,
		// }
		token.RIGHT_BRACE,
		// var result = add(5, 10);
		token.VAR, token.IDENTIFIER, token.EQUAL, token.IDENTIFIER, token.LEFT_PAREN, token.NUMBER, token.COMMA, token.NUMBER, token.RIGHT_PAREN, token.SEMICOLON,
		// print result;
		token.PRINT, token.IDENTIFIER, token.SEMICOLON,
		// EOF
		token.EOF,
	}

	assertTokenTypes(t, tokens, expected)
}

func TestScanningErrors(t *testing.T) {
	errorCases := []struct {
		name  string
		input string
	}{
		{
			name:  "unexpected character @",
			input: "var x = @;",
		},
		{
			name:  "unterminated double-quoted string at EOF",
			input: `"cadena sin cerrar`,
		},
	}

	for _, tt := range errorCases {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScanner(tt.input)
			tokens, err := s.Scan()
			if err == nil {
				t.Fatalf("expected scanning error for %q, but got none (tokens: %v)", tt.input, tokens)
			}
			if tokens != nil {
				t.Errorf("expected tokens to be nil on error, got %v", tokens)
			}
		})
	}
}

func assertTokenTypes(t *testing.T, actual []token.Token, expected []token.TokenType) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Fatalf("token count mismatch: expected %d tokens, got %d", len(expected), len(actual))
	}

	for i := range expected {
		if actual[i].Type != expected[i] {
			t.Errorf("token [%d] mismatch: expected %v (%s), got %v (%s) [lexeme: %q]",
				i, expected[i], expected[i], actual[i].Type, actual[i].Type, actual[i].Lexeme)
		}
	}
}
