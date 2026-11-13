package vm

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"glox/pkg/chunk"
	"glox/pkg/compiler"
	"glox/pkg/scanner"
)

func runSource(t *testing.T, source string) (string, error) {
	t.Helper()
	tokens, err := scanner.NewScanner(source).Scan()
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
	c, err := compiler.NewCompiler(tokens).Compile()
	if err != nil {
		t.Fatalf("unexpected compile error: %v", err)
	}

	var out bytes.Buffer
	machine := NewVM()
	machine.SetWriter(&out)
	err = machine.Run(c)
	return out.String(), err
}

func TestVM_Print(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"number", "print 42;", "42\n"},
		{"float", "print 3.14;", "3.14\n"},
		{"string", `print "hola fiuba";`, "hola fiuba\n"},
		{"true", "print true;", "true\n"},
		{"false", "print false;", "false\n"},
		{"nil", "print nil;", "nil\n"},
		{"grouping", "print (1);", "1\n"},
		{"expression statement prints nothing", "1;", ""},
		{"many statements", "print 1; 2; print 3;", "1\n3\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runSource(t, tt.input)
			if err != nil {
				t.Fatalf("unexpected runtime error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestVM_Operators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"negate", "print -3;", "-3"},
		{"not", "print !true;", "false"},
		{"not nil", "print !nil;", "true"},
		{"not number", "print !0;", "false"},
		{"add", "print 1 + 2;", "3"},
		{"concat", `print "ab" + "cd";`, "abcd"},
		{"subtract", "print 5 - 7;", "-2"},
		{"multiply", "print 4 * 2.5;", "10"},
		{"divide", "print 7 / 2;", "3.5"},
		{"modulo", "print 7 % 3;", "1"},
		{"precedence", "print 1 + 2 * 3;", "7"},
		{"associativity", "print 10 - 4 - 3;", "3"},
		{"grouping", "print (1 + 2) * 3;", "9"},
		{"equal", "print 1 == 1;", "true"},
		{"equal strings", `print "a" == "a";`, "true"},
		{"equal different types", `print 1 == "1";`, "false"},
		{"equal nil", "print nil == nil;", "true"},
		{"not equal", "print 1 != 2;", "true"},
		{"greater", "print 2 > 1;", "true"},
		{"greater equal", "print 1 >= 1;", "true"},
		{"less", "print 2 < 1;", "false"},
		{"less equal", "print 2 <= 1;", "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runSource(t, tt.input)
			if err != nil {
				t.Fatalf("unexpected runtime error: %v", err)
			}
			if got != tt.expected+"\n" {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestVM_OperatorErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"negate string", `-"a";`, "Operand of OP_NEGATE must be a number"},
		{"add number and string", `1 + "a";`, "Operands of OP_ADD must be either numbers or strings"},
		{"subtract bools", "true - false;", "Operands of OP_SUBTRACT must be numbers"},
		{"compare strings", `"a" < "b";`, "Operands of OP_LESS must be numbers"},
		{"divide by zero", "1 / 0;", "Division by 0 is not allowed"},
		{"modulo by zero", "1 % 0;", "Division by 0 is not allowed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runSource(t, tt.input)
			var runtimeErr *RuntimeError
			if !errors.As(err, &runtimeErr) {
				t.Fatalf("expected RuntimeError, got %v", err)
			}
			if !strings.Contains(err.Error(), tt.expected) {
				t.Errorf("got %q, want it to contain %q", err.Error(), tt.expected)
			}
		})
	}
}

func TestVM_GlobalVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"declare and print", "var a = 1; print a;", "1\n"},
		{"declare without value", "var a; print a;", "nil\n"},
		{"assign", "var a = 1; a = 2; print a;", "2\n"},
		{"assignment is an expression", "var a; print a = 3;", "3\n"},
		{"chained assignment", "var a; var b; a = b = 4; print a; print b;", "4\n4\n"},
		{"use in expressions", `var a = "ho"; var b = a + "la"; print b;`, "hola\n"},
		{"redeclare", "var a = 1; var a = 2; print a;", "2\n"},
		{"self reference", "var a = 1; a = a + 1; print a;", "2\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runSource(t, tt.input)
			if err != nil {
				t.Fatalf("unexpected runtime error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestVM_UndefinedVariables(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"reference", "print a;"},
		{"assignment", "a = 1;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runSource(t, tt.input)
			if err == nil || !strings.Contains(err.Error(), "Undefined variable 'a'.") {
				t.Errorf("expected undefined variable error, got %v", err)
			}
		})
	}
}

// Las globales sobreviven entre ejecuciones, como entre lineas del REPL
func TestVM_GlobalsPersistBetweenRuns(t *testing.T) {
	machine := NewVM()
	var out bytes.Buffer
	machine.SetWriter(&out)

	for _, source := range []string{"var a = 1;", "a = a + 1;", "print a;"} {
		tokens, _ := scanner.NewScanner(source).Scan()
		c, err := compiler.NewCompiler(tokens).Compile()
		if err != nil {
			t.Fatalf("unexpected compile error: %v", err)
		}
		if err := machine.Run(c); err != nil {
			t.Fatalf("unexpected runtime error: %v", err)
		}
	}

	if out.String() != "2\n" {
		t.Errorf("got %q, want %q", out.String(), "2\n")
	}
}

func TestVM_LocalVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"declare and print", "{ var a = 1; print a; }", "1\n"},
		{"declare without value", "{ var a; print a; }", "nil\n"},
		{"assign", "{ var a = 1; a = 2; print a; }", "2\n"},
		{"assignment is an expression", "{ var a; print a = 3; }", "3\n"},
		{"many locals", "{ var a = 1; var b = 2; var c = a + b; print c; }", "3\n"},
		{"nested scopes", "{ var a = 1; { var b = 2; print a + b; } print a; }", "3\n1\n"},
		{"shadowing", "{ var a = 1; { var a = 2; print a; } print a; }", "2\n1\n"},
		{"shadowing a global", "var a = 1; { var a = 2; print a; } print a;", "2\n1\n"},
		{"assign global from block", "var a = 1; { a = 2; } print a;", "2\n"},
		{"assign outer local", "{ var a = 1; { a = 2; } print a; }", "2\n"},
		{"initialize with outer variable", "var a = 1; { var b = a + 1; print b; }", "2\n"},
		{"slot reused after block", "{ { var a = 1; } var b = 2; print b; }", "2\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runSource(t, tt.input)
			if err != nil {
				t.Fatalf("unexpected runtime error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestVM_StackIsEmptyAfterStatements(t *testing.T) {
	tokens, _ := scanner.NewScanner("1; print 2; (nil); { var a = 1; { var b = a; } }").Scan()
	c, _ := compiler.NewCompiler(tokens).Compile()

	machine := NewVM()
	machine.SetWriter(&bytes.Buffer{})
	if err := machine.Run(c); err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if len(machine.stack) != 0 {
		t.Errorf("expected empty stack, got %v", machine.stack)
	}
}

func TestVM_RuntimeErrors(t *testing.T) {
	underflow := chunk.NewChunk()
	underflow.Write(byte(chunk.OP_POP))

	unknown := chunk.NewChunk()
	unknown.Write(255)

	tests := []struct {
		name  string
		chunk *chunk.Chunk
	}{
		{"stack underflow", underflow},
		{"unknown opcode", unknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewVM().Run(tt.chunk)
			var runtimeErr *RuntimeError
			if !errors.As(err, &runtimeErr) {
				t.Errorf("expected RuntimeError, got %v", err)
			}
		})
	}
}
