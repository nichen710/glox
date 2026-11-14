package compiler

import (
	"errors"
	"fmt"
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

func TestCompiler_Operators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"negate", "-1;", []string{"OP_CONSTANT<1>", "OP_NEGATE"}},
		{"not", "!true;", []string{"OP_TRUE", "OP_NOT"}},
		{"double not", "!!nil;", []string{"OP_NIL", "OP_NOT", "OP_NOT"}},
		{"add", "1 + 2;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_ADD"}},
		{"subtract", "1 - 2;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_SUBTRACT"}},
		{"multiply", "1 * 2;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_MULTIPLY"}},
		{"divide", "1 / 2;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_DIVIDE"}},
		{"modulo", "1 % 2;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_MODULO"}},
		{"equal", "1 == 2;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_EQUAL"}},
		{"not equal", "1 != 2;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_EQUAL", "OP_NOT"}},
		{"greater", "1 > 2;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_GREATER"}},
		{"greater equal", "1 >= 2;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_LESS", "OP_NOT"}},
		{"less", "1 < 2;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_LESS"}},
		{"less equal", "1 <= 2;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_GREATER", "OP_NOT"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := compileSource(t, tt.input)
			if err != nil {
				t.Fatalf("unexpected compile error on %q: %v", tt.input, err)
			}
			assertInstructions(t, c, append(tt.expected, "OP_POP", "OP_RETURN"))
		})
	}
}

func TestCompiler_PrecedenceAndAssociativity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"factor before term", "1 + 2 * 3;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_CONSTANT<3>", "OP_MULTIPLY", "OP_ADD"}},
		{"left associative", "1 - 2 - 3;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_SUBTRACT", "OP_CONSTANT<3>", "OP_SUBTRACT"}},
		{"grouping changes order", "(1 + 2) * 3;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_ADD", "OP_CONSTANT<3>", "OP_MULTIPLY"}},
		{"unary before factor", "-1 * 2;", []string{"OP_CONSTANT<1>", "OP_NEGATE", "OP_CONSTANT<2>", "OP_MULTIPLY"}},
		{"comparison before equality", "1 < 2 == true;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_LESS", "OP_TRUE", "OP_EQUAL"}},
		{"term before comparison", "1 + 2 > 3;", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_ADD", "OP_CONSTANT<3>", "OP_GREATER"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := compileSource(t, tt.input)
			if err != nil {
				t.Fatalf("unexpected compile error on %q: %v", tt.input, err)
			}
			assertInstructions(t, c, append(tt.expected, "OP_POP", "OP_RETURN"))
		})
	}
}

func assertInstructions(t *testing.T, c *chunk.Chunk, expected []string) {
	t.Helper()
	var got []string
	for _, line := range strings.Split(c.String(), "\n") {
		got = append(got, strings.SplitN(line, "| ", 2)[1])
	}
	if strings.Join(got, " ") != strings.Join(expected, " ") {
		t.Errorf("got:\n%v\nwant:\n%v", got, expected)
	}
}

func TestCompiler_GlobalVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"declaration", "var a = 1;", []string{"OP_CONSTANT<1>", "OP_DEFINE_GLOBAL<'a'>"}},
		{"declaration without value", "var a;", []string{"OP_NIL", "OP_DEFINE_GLOBAL<'a'>"}},
		{"reference", "print a;", []string{"OP_GET_GLOBAL<'a'>", "OP_PRINT"}},
		{"assignment", "a = 2;", []string{"OP_CONSTANT<2>", "OP_SET_GLOBAL<'a'>", "OP_POP"}},
		{"chained assignment", "a = b = 2;", []string{"OP_CONSTANT<2>", "OP_SET_GLOBAL<'b'>", "OP_SET_GLOBAL<'a'>", "OP_POP"}},
		{"reference in expression", "print a + 1;", []string{"OP_GET_GLOBAL<'a'>", "OP_CONSTANT<1>", "OP_ADD", "OP_PRINT"}},
		{"assignment with expression", "a = a * 2;", []string{"OP_GET_GLOBAL<'a'>", "OP_CONSTANT<2>", "OP_MULTIPLY", "OP_SET_GLOBAL<'a'>", "OP_POP"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := compileSource(t, tt.input)
			if err != nil {
				t.Fatalf("unexpected compile error on %q: %v", tt.input, err)
			}
			assertInstructions(t, c, append(tt.expected, "OP_RETURN"))
		})
	}
}

func TestCompiler_LocalVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty block", "{}", []string{}},
		{"declaration", "{ var a = 1; }", []string{"OP_CONSTANT<1>", "OP_POP"}},
		{"reference", "{ var a = 1; print a; }", []string{"OP_CONSTANT<1>", "OP_GET_LOCAL<@0>", "OP_PRINT", "OP_POP"}},
		{"assignment", "{ var a; a = 2; }", []string{"OP_NIL", "OP_CONSTANT<2>", "OP_SET_LOCAL<@0>", "OP_POP", "OP_POP"}},
		{"second slot", "{ var a = 1; var b = 2; print b; }", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_GET_LOCAL<@1>", "OP_PRINT", "OP_POP", "OP_POP"}},
		{"nested block", "{ var a = 1; { var b = a; } }", []string{"OP_CONSTANT<1>", "OP_GET_LOCAL<@0>", "OP_POP", "OP_POP"}},
		{"shadowing", "{ var a = 1; { var a = 2; print a; } }", []string{"OP_CONSTANT<1>", "OP_CONSTANT<2>", "OP_GET_LOCAL<@1>", "OP_PRINT", "OP_POP", "OP_POP"}},
		{"global inside block", "{ print a; }", []string{"OP_GET_GLOBAL<'a'>", "OP_PRINT"}},
		{"slot reused after block", "{ { var a = 1; } var b = 2; print b; }", []string{"OP_CONSTANT<1>", "OP_POP", "OP_CONSTANT<2>", "OP_GET_LOCAL<@0>", "OP_PRINT", "OP_POP"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := compileSource(t, tt.input)
			if err != nil {
				t.Fatalf("unexpected compile error on %q: %v", tt.input, err)
			}
			assertInstructions(t, c, append(tt.expected, "OP_RETURN"))
		})
	}
}

func TestCompiler_ControlFlow(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"if", "if (true) print 1;", []string{"OP_TRUE", "OP_JUMP_IF_FALSE<+10>", "OP_POP", "OP_CONSTANT<1>", "OP_PRINT", "OP_JUMP<+4>", "OP_POP"}},
		{"if else", "if (true) print 1; else print 2;", []string{"OP_TRUE", "OP_JUMP_IF_FALSE<+10>", "OP_POP", "OP_CONSTANT<1>", "OP_PRINT", "OP_JUMP<+7>", "OP_POP", "OP_CONSTANT<2>", "OP_PRINT"}},
		{"and", "true and false;", []string{"OP_TRUE", "OP_JUMP_IF_FALSE<+5>", "OP_POP", "OP_FALSE", "OP_POP"}},
		{"or", "true or false;", []string{"OP_TRUE", "OP_JUMP_IF_FALSE<+6>", "OP_JUMP<+5>", "OP_POP", "OP_FALSE", "OP_POP"}},
		{"while", "while (false) print 1;", []string{"OP_FALSE", "OP_JUMP_IF_FALSE<+10>", "OP_POP", "OP_CONSTANT<1>", "OP_PRINT", "OP_LOOP<-8>", "OP_POP"}},
		{"empty for", "for (;;) print 1;", []string{"OP_CONSTANT<1>", "OP_PRINT", "OP_LOOP<-3>"}},
		{"for with local", "for (var i = 0; i < 2; i = i + 1) print i;", []string{
			"OP_CONSTANT<0>",
			"OP_GET_LOCAL<@0>", "OP_CONSTANT<2>", "OP_LESS", "OP_JUMP_IF_FALSE<+24>", "OP_POP",
			"OP_JUMP<+14>",
			"OP_GET_LOCAL<@0>", "OP_CONSTANT<1>", "OP_ADD", "OP_SET_LOCAL<@0>", "OP_POP", "OP_LOOP<-20>",
			"OP_GET_LOCAL<@0>", "OP_PRINT", "OP_LOOP<-14>",
			"OP_POP", "OP_POP",
		}},
		{"for with global initializer", "for (i = 0; false;) {}", []string{"OP_CONSTANT<0>", "OP_SET_GLOBAL<'i'>", "OP_POP", "OP_FALSE", "OP_JUMP_IF_FALSE<+7>", "OP_POP", "OP_LOOP<-5>", "OP_POP"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := compileSource(t, tt.input)
			if err != nil {
				t.Fatalf("unexpected compile error on %q: %v", tt.input, err)
			}
			assertInstructions(t, c, append(tt.expected, "OP_RETURN"))
		})
	}
}

func TestCompiler_TooMuchCodeToJump(t *testing.T) {
	// Cada `print a;` ocupa 3 bytes y no agrega constantes porque `a` es local
	body := "var a;" + strings.Repeat("print a;", chunk.MaxJump/3+1)
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"if", "if (true) {" + body + "}", "Too much code to jump over."},
		{"loop", "for (;;) {" + body + "}", "Loop body too large."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := compileSource(t, tt.input)
			if err == nil || !strings.Contains(err.Error(), tt.expected) {
				t.Errorf("got %v, want it to contain %q", err, tt.expected)
			}
		})
	}
}

func TestCompiler_TooManyLocals(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("{")
	for i := 0; i <= MaxLocals; i++ {
		fmt.Fprintf(&sb, "var a%d;", i)
	}
	sb.WriteString("}")

	_, err := compileSource(t, sb.String())
	if err == nil || !strings.Contains(err.Error(), "Too many local variables") {
		t.Errorf("expected too many locals error, got %v", err)
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
		{"missing right operand", "1 +;", "Expected expression"},
		{"binary without left operand", "* 2;", "Expected expression"},
		{"var without name", "var = 1;", "Expected variable name."},
		{"var without semicolon", "var a = 1", "Expected ';' after variable declaration."},
		{"invalid assignment target", "a + b = 1;", "Invalid assignment target."},
		{"assign to literal", "1 = 2;", "Invalid assignment target."},
		{"assign to grouping", "(a) = 2;", "Invalid assignment target."},
		{"unclosed block", "{ print 1;", "Expected '}' after block."},
		{"redeclare local", "{ var a = 1; var a = 2; }", "Already a variable with this name in this scope."},
		{"read local in own initializer", "{ var a = a; }", "Can't read local variable in its own initializer."},
		{"shadow with own value", "var a = 1; { var a = a + 1; }", "Can't read local variable in its own initializer."},
		{"if without left paren", "if true) print 1;", "Expected '(' after 'if'."},
		{"if without right paren", "if (true print 1;", "Expected ')' after if condition."},
		{"while without left paren", "while true) print 1;", "Expected '(' after 'while'."},
		{"while without right paren", "while (true print 1;", "Expected ')' after while condition."},
		{"for without left paren", "for var i = 0;;) print 1;", "Expected '(' after 'for'."},
		{"for without condition semicolon", "for (;true) print 1;", "Expected ';' after loop condition."},
		{"for without right paren", "for (;; 1 print 1;", "Expected ')' after for clauses."},
		{"var as if body", "if (true) var a = 1;", "Expected expression"},
		{"var as while body", "while (false) var a = 1;", "Expected expression"},
		{"and without right operand", "true and;", "Expected expression"},
		{"assign to logical", "a or b = 1;", "Invalid assignment target."},
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
