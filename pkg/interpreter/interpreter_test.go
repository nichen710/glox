package interpreter

import (
	"bytes"
	"strings"
	"testing"

	"glox/pkg/parser"
	"glox/pkg/resolver"
	"glox/pkg/scanner"
)

func evaluateSource(t *testing.T, source string, out *bytes.Buffer) (any, error) {
	t.Helper()
	trimmed := strings.TrimSpace(source)
	if !strings.HasSuffix(trimmed, ";") && !strings.HasSuffix(trimmed, "}") {
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

func TestInterpreter_LogicalOperators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{"or returns first truthy", `"hi" or 2`, "hi"},
		{"or returns right if left falsey", `nil or "yes"`, "yes"},
		{"or returns right if both falsey", `false or nil`, nil},
		{"and returns left if falsey", `false and "no"`, false},
		{"and returns left if nil", `nil and "no"`, nil},
		{"and returns right if left truthy", `"hi" and "there"`, "there"},
		{"and with numbers", `1 and 2`, 2.0},
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
		{"accessing inner block variable from outer scope", "{ var inner = 123; } print inner;", "Undefined variable 'inner'."},
		{"calling non-callable", `"hello"();`, "Can only call functions and classes."},
		{"call arity too few", `fun add(a, b) { return a + b; } add(1);`, "Expected 2 arguments but got 1."},
		{"call arity too many", `fun add(a, b) { return a + b; } add(1, 2, 3);`, "Expected 2 arguments but got 3."},
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
		{
			name:           "block statement execution",
			source:         "{ var a = 1; print a; }",
			expectedOutput: "1\n",
		},
		{
			name: "block variable shadowing",
			source: `
var a = "global";
{
    var a = "local";
    print a;
}
print a;
`,
			expectedOutput: "local\nglobal\n",
		},
		{
			name: "nested scope chains",
			source: `
var a = "global";
{
    var b = "outer";
    {
        var c = "inner";
        print a;
        print b;
        print c;
    }
}
`,
			expectedOutput: "global\nouter\ninner\n",
		},
		{
			name: "mutate outer variable from inner scope",
			source: `
var x = 1;
{
    x = 2;
    print x;
}
print x;
`,
			expectedOutput: "2\n2\n",
		},
		{
			name: "if true executes then branch",
			source: `
if (true) print "yes";
`,
			expectedOutput: "yes\n",
		},
		{
			name: "if false without else does not execute then branch",
			source: `
if (false) print "yes";
`,
			expectedOutput: "",
		},
		{
			name: "if false with else executes else branch",
			source: `
if (false) print "yes"; else print "no";
`,
			expectedOutput: "no\n",
		},
		{
			name: "if with block branches and scopes",
			source: `
var x = "global";
if (true) {
    var x = "local";
    print x;
}
print x;
`,
			expectedOutput: "local\nglobal\n",
		},
		{
			name: "dangling else binds to nearest if",
			source: `
if (true) if (false) print "bad"; else print "good";
`,
			expectedOutput: "good\n",
		},
		{
			name: "logical or short circuits side effects",
			source: `
var a = "safe";
true or (a = "modified");
print a;
`,
			expectedOutput: "safe\n",
		},
		{
			name: "logical and short circuits side effects",
			source: `
var a = "safe";
false and (a = "modified");
print a;
`,
			expectedOutput: "safe\n",
		},
		{
			name: "while loop counter",
			source: `
var i = 0;
while (i < 3) {
    print i;
    i = i + 1;
}
`,
			expectedOutput: "0\n1\n2\n",
		},
		{
			name: "while loop with false condition never executes body",
			source: `
var ran = false;
while (false) {
    ran = true;
}
print ran;
`,
			expectedOutput: "false\n",
		},
		{
			name: "while loop block scope variables",
			source: `
var i = 0;
while (i < 2) {
    var msg = "item";
    print msg;
    i = i + 1;
}
`,
			expectedOutput: "item\nitem\n",
		},
		{
			name: "for loop standard iteration",
			source: `
for (var i = 0; i < 3; i = i + 1) {
    print i;
}
`,
			expectedOutput: "0\n1\n2\n",
		},
		{
			name: "for loop without initializer or increment",
			source: `
var i = 0;
for (; i < 3;) {
    print i;
    i = i + 1;
}
print i;
`,
			expectedOutput: "0\n1\n2\n3\n",
		},
		{
			name: "for loop variable scope isolation",
			source: `
var i = "outer";
for (var i = 0; i < 1; i = i + 1) {
    print i;
}
print i;
`,
			expectedOutput: "0\nouter\n",
		},
		{
			name: "function declaration and call with print",
			source: `
fun sayHello(name) {
    print "hello " + name;
}
sayHello("world");
`,
			expectedOutput: "hello world\n",
		},
		{
			name: "function without return defaults to nil",
			source: `
fun noReturn() {}
print noReturn();
`,
			expectedOutput: "nil\n",
		},
		{
			name: "function return with value",
			source: `
fun add(a, b) {
    return a + b;
}
print add(3, 4);
`,
			expectedOutput: "7\n",
		},
		{
			name: "function early return in conditional",
			source: `
fun check(n) {
    if (n < 0) return "negative";
    return "non-negative";
}
print check(-5);
print check(5);
`,
			expectedOutput: "negative\nnon-negative\n",
		},
		{
			name: "recursive fibonacci",
			source: `
fun fib(n) {
    if (n <= 1) return n;
    return fib(n - 2) + fib(n - 1);
}
print fib(7);
`,
			expectedOutput: "13\n",
		},
		{
			name: "nested function closure preserves state",
			source: `
fun makeCounter() {
    var count = 0;
    fun increment() {
        count = count + 1;
        return count;
    }
    return increment;
}
var counter1 = makeCounter();
var counter2 = makeCounter();
print counter1();
print counter1();
print counter2();
print counter1();
`,
			expectedOutput: "1\n2\n1\n3\n",
		},
		{
			name: "function as first-class value",
			source: `
fun apply(fn, x) {
    return fn(x);
}
fun double(x) {
    return x * 2;
}
print apply(double, 21);
`,
			expectedOutput: "42\n",
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

func TestInterpreter_NativeClock(t *testing.T) {
	val := mustEvaluate(t, "clock()")
	num, ok := val.(float64)
	if !ok {
		t.Fatalf("expected clock() to return float64, got %T (%v)", val, val)
	}
	if num <= 0 {
		t.Errorf("expected positive timestamp, got %v", num)
	}
}

func TestInterpreter_WithResolver(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		expectedOutput string
	}{
		{
			name: "closure retains static scope instead of dynamic mutable scope",
			source: `
var a = "global";
{
    fun ret_a() {
        return a;
    }

    print ret_a();
    var a = "block";
    print ret_a();
}
`,
			expectedOutput: "global\nglobal\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := scanner.NewScanner(tt.source).Scan()
			if err != nil {
				t.Fatalf("scan error: %v", err)
			}
			stmts, err := parser.NewParser(tokens).Parse()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			table := resolver.NewBindingTable()
			rslvr := resolver.NewResolverBuilder().WithBinding(table).Build()
			if err := rslvr.Resolve(stmts); err != nil {
				t.Fatalf("resolver error: %v", err)
			}

			var buf bytes.Buffer
			interp := NewInterpreter()
			interp.SetWriter(&buf)
			interp.SetBindings(table)

			_, err = interp.Interpret(stmts)
			if err != nil {
				t.Fatalf("interpret error: %v", err)
			}

			if buf.String() != tt.expectedOutput {
				t.Errorf("got %q, want %q", buf.String(), tt.expectedOutput)
			}
		})
	}
}



