package vm

import (
	"bytes"
	"errors"
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

func TestVM_StackIsEmptyAfterStatements(t *testing.T) {
	tokens, _ := scanner.NewScanner("1; print 2; (nil);").Scan()
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
