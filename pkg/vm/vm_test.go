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

func TestVM_ControlFlow(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"if true", `if (true) print "si";`, "si\n"},
		{"if false", `if (false) print "si";`, ""},
		{"if else", `if (nil) print "si"; else print "no";`, "no\n"},
		{"zero is truthy", `if (0) print "si";`, "si\n"},
		{"else binds to nearest if", `if (true) if (false) print 1; else print 2;`, "2\n"},
		{"and returns left if falsey", "print false and 1;", "false\n"},
		{"and returns right", "print 1 and 2;", "2\n"},
		{"or returns left if truthy", "print 1 or 2;", "1\n"},
		{"or returns right", "print nil or 2;", "2\n"},
		{"and binds tighter than or", "print false and 1 or 2;", "2\n"},
		{"and short circuits", "var a = 1; false and (a = 2); print a;", "1\n"},
		{"or short circuits", "var a = 1; true or (a = 2); print a;", "1\n"},
		{"while", "var i = 0; while (i < 3) { print i; i = i + 1; }", "0\n1\n2\n"},
		{"while false", `while (false) print "nunca";`, ""},
		{"for", "for (var i = 0; i < 3; i = i + 1) print i;", "0\n1\n2\n"},
		{"for without increment", "for (var i = 0; i < 2;) { print i; i = i + 1; }", "0\n1\n"},
		{"for with global", "var i; for (i = 0; i < 2; i = i + 1) {} print i;", "2\n"},
		{"for local does not leak", "var i = \"global\"; for (var i = 0; i < 1; i = i + 1) {} print i;", "global\n"},
		{"nested loops", "var n = 0; for (var a = 0; a < 3; a = a + 1) for (var b = 0; b < 3; b = b + 1) n = n + 1; print n;", "9\n"},
		{"fibonacci", "var a = 0; var temp; for (var b = 1; a < 20; b = temp + b) { print a; temp = a; a = b; }", "0\n1\n1\n2\n3\n5\n8\n13\n"},
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
	tokens, _ := scanner.NewScanner("1; print 2; (nil); { var a = 1; { var b = a; } }" +
		"if (false) 1; else 2; if (true) 1; 1 and 2; nil or 2; for (var i = 0; i < 2; i = i + 1) { var j = i; }").Scan()
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
	underflow.Write(byte(chunk.OP_POP), 1)

	unknown := chunk.NewChunk()
	unknown.Write(255, 1)

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

func TestVM_RuntimeErrorLine(t *testing.T) {
	source := "var a = 1;\nvar b = 2;\nprint a / 0;"
	_, err := runSource(t, source)
	var runtimeErr *RuntimeError
	if !errors.As(err, &runtimeErr) {
		t.Fatalf("expected RuntimeError, got %v", err)
	}
	if runtimeErr.Line != 3 {
		t.Errorf("got line %d, want 3", runtimeErr.Line)
	}
	if !strings.Contains(err.Error(), "[line 3 |") {
		t.Errorf("expected formatted error with line 3, got %q", err.Error())
	}
}

func TestVM_Functions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"print function", "fun f() {} print f;", "<fn f>\n"},
		{"print native", "print clock;", "<native fn>\n"},
		{"implicit nil return", "fun f() { var x = 1; } print f();", "nil\n"},
		{"bare return", "fun f() { return; print 1; } print f();", "nil\n"},
		{"return value", "fun add(a, b) { return a + b; } print add(2, 3);", "5\n"},
		{"nested calls", "fun sq(x) { return x * x; } fun sumSq(a, b) { return sq(a) + sq(b); } print sumSq(3, 4);", "25\n"},
		{"recursion", "fun fib(n) { if (n <= 1) return n; return fib(n - 1) + fib(n - 2); } print fib(10);", "55\n"},
		{"local function", "{ fun mul(a, b) { return a * b; } print mul(6, 7); }", "42\n"},
		{"clock native", "var t = clock(); print t > 0;", "true\n"},
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

func TestVM_FunctionErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"call non-function", "var x = 1; x();", "Can only call functions."},
		{"wrong arity", "fun f(a, b) {} f(1);", "Expected 2 arguments but got 1."},
		{"native wrong arity", "clock(1);", "Expected 0 arguments but got 1."},
		{"stack overflow", "fun f() { f(); } f();", "Stack overflow."},
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

func TestVM_ClosuresAndUpvalues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"counter",
			`fun makeCounter() {
				var count = 0;
				fun inc() {
					count = count + 1;
					return count;
				}
				return inc;
			}
			var c = makeCounter();
			print c();
			print c();`,
			"1\n2\n",
		},
		{
			"independent counters",
			`fun makeCounter() {
				var count = 0;
				fun inc() {
					count = count + 1;
					return count;
				}
				return inc;
			}
			var a = makeCounter();
			var b = makeCounter();
			print a();
			print a();
			print b();`,
			"1\n2\n1\n",
		},
		{
			"nested upvalues",
			`fun outer() {
				var x = "outside";
				fun middle() {
					fun inner() {
						print x;
					}
					return inner;
				}
				return middle;
			}
			var mid = outer();
			var in = mid();
			in();`,
			"outside\n",
		},
		{
			"shared open upvalue mutated",
			`fun test() {
				var a = 1;
				fun get() { return a; }
				fun set(v) { a = v; }
				print get();
				set(42);
				print get();
			}
			test();`,
			"1\n42\n",
		},
		{"closed on block exit",
			`var f;
			{
				var x = "block";
				fun inner() { print x; }
				f = inner;
			}
			f();`,
			"block\n",
		},
		{
			"deep tail recursion",
			`fun sumTail(n, acc) {
				if (n <= 0) return acc;
				return sumTail(n - 1, acc + n);
			}
			print sumTail(10000, 0);`,
			"50005000\n",
		},
		{
			"tail call closes upvalues",
			`fun make(n) {
				var x = n;
				fun get() { return x; }
				if (n <= 0) return get;
				return make(n - 1);
			}
			var g = make(5);
			print g();`,
			"0\n",
		},
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

func TestVM_ConstantLong(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 300; i++ {
		sb.WriteString(strings.ReplaceAll("N;", "N", string(rune('0'+i/100))+string(rune('0'+(i/10)%10))+string(rune('0'+i%10))))
	}
	sb.WriteString("print 999;")
	got, err := runSource(t, sb.String())
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if got != "999\n" {
		t.Errorf("got %q, want %q", got, "999\n")
	}
}

func TestVM_Natives(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"random range", "var r = random(1, 5); print r >= 1 and r <= 5;", "true\n"},
		{"random single", "print random(7, 7);", "7\n"},
		{"sleep", "sleep(1); print true;", "true\n"},
		{"toNumber valid", `print toNumber("42.5") + 1;`, "43.5\n"},
		{"toNumber invalid", `print toNumber("abc");`, "nil\n"},
		{"clear", "clear();", "\x1b[2J\x1b[H"},
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
