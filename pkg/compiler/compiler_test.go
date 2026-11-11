package compiler

import (
	"errors"
	"strings"
	"testing"

	"glox/pkg/chunk"
	"glox/pkg/scanner"
)

func compileSource(t *testing.T, source string) (*chunk.Chunk, error) {
	t.Helper()
	tokens, err := scanner.NewScanner(source).Scan()
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
	return NewCompiler(tokens).Compile()
}

func TestCompiler_Statements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty program", "", "0000| OP_RETURN"},
		{"print number", "print 2;", "0000| OP_CONSTANT<2>\n0002| OP_PRINT\n0003| OP_RETURN"},
		{"print string", `print "hola";`, "0000| OP_CONSTANT<'hola'>\n0002| OP_PRINT\n0003| OP_RETURN"},
		{"expression statement", "2;", "0000| OP_CONSTANT<2>\n0002| OP_POP\n0003| OP_RETURN"},
		{"many statements", "1; print 2;", "0000| OP_CONSTANT<1>\n0002| OP_POP\n0003| OP_CONSTANT<2>\n0005| OP_PRINT\n0006| OP_RETURN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := compileSource(t, tt.input)
			if err != nil {
				t.Fatalf("unexpected compile error on %q: %v", tt.input, err)
			}
			if got := c.String(); got != tt.expected {
				t.Errorf("got:\n%s\nwant:\n%s", got, tt.expected)
			}
		})
	}
}

func TestCompiler_Literals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"nil", "nil;", "0000| OP_NIL\n0001| OP_POP\n0002| OP_RETURN"},
		{"true", "true;", "0000| OP_TRUE\n0001| OP_POP\n0002| OP_RETURN"},
		{"false", "false;", "0000| OP_FALSE\n0001| OP_POP\n0002| OP_RETURN"},
		{"grouping", "((nil));", "0000| OP_NIL\n0001| OP_POP\n0002| OP_RETURN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := compileSource(t, tt.input)
			if err != nil {
				t.Fatalf("unexpected compile error on %q: %v", tt.input, err)
			}
			if got := c.String(); got != tt.expected {
				t.Errorf("got:\n%s\nwant:\n%s", got, tt.expected)
			}
		})
	}
}

func TestCompiler_Errors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"missing semicolon", "print 1", "Expected ';' after value."},
		{"missing expression semicolon", "1", "Expected ';' after expression."},
		{"unclosed grouping", "(1;", "Expected ')' after expression."},
		{"missing expression", "print;", "Expected expression"},
		{"not supported yet", "1 + 2;", "Expected ';' after expression."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := compileSource(t, tt.input)
			var compileErr *CompileError
			if !errors.As(err, &compileErr) {
				t.Fatalf("expected CompileError, got %v", err)
			}
			if !strings.Contains(err.Error(), tt.expected) {
				t.Errorf("got %q, want it to contain %q", err.Error(), tt.expected)
			}
		})
	}
}

func TestCompiler_TooManyConstants(t *testing.T) {
	source := strings.Repeat("1;", chunk.MaxConstants+1)
	_, err := compileSource(t, source)
	if err == nil || !strings.Contains(err.Error(), "Too many constants") {
		t.Errorf("expected too many constants error, got %v", err)
	}
}
