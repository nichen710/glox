package resolver

import (
	"strings"
	"testing"

	"glox/pkg/expression"
	"glox/pkg/parser"
	"glox/pkg/scanner"
	"glox/pkg/statement"
)

type mockStep struct {
	name   string
	called bool
}

func (m *mockStep) Name() string { return m.name }
func (m *mockStep) Run(statements []statement.Statement) error {
	m.called = true
	return nil
}

func parseSource(t *testing.T, source string) []statement.Statement {
	t.Helper()
	tokens, err := scanner.NewScanner(source).Scan()
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
	stmts, err := parser.NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	return stmts
}

func TestResolver_VariableDepthResolution(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		validate func(t *testing.T, stmts []statement.Statement, table *BindingTable)
	}{
		{
			name: "local and nested scope variable depths",
			source: `
var global = "g";
{
    var outer = "o";
    {
        var inner = "i";
        print inner;
        print outer;
        print global;
        inner = "new_i";
        outer = "new_o";
    }
}
`,
			validate: func(t *testing.T, stmts []statement.Statement, table *BindingTable) {
				outerBlock := stmts[1].(statement.Block)
				innerBlock := outerBlock.Statements[1].(statement.Block)

				printInner := innerBlock.Statements[1].(statement.Print)
				innerVar := printInner.Expression.(expression.Variable)
				if depth, ok := table.Depth(innerVar.ID); !ok || depth != 0 {
					t.Errorf("expected innerVar depth 0, got depth %d (ok=%v)", depth, ok)
				}

				printOuter := innerBlock.Statements[2].(statement.Print)
				outerVar := printOuter.Expression.(expression.Variable)
				if depth, ok := table.Depth(outerVar.ID); !ok || depth != 1 {
					t.Errorf("expected outerVar depth 1, got depth %d (ok=%v)", depth, ok)
				}

				printGlobal := innerBlock.Statements[3].(statement.Print)
				globalVar := printGlobal.Expression.(expression.Variable)
				if _, ok := table.Depth(globalVar.ID); ok {
					t.Errorf("expected global variable to NOT be in binding table")
				}

				assignInnerStmt := innerBlock.Statements[4].(statement.ExpressionStatement)
				assignInner := assignInnerStmt.Expression.(expression.Assign)
				if depth, ok := table.Depth(assignInner.ID); !ok || depth != 0 {
					t.Errorf("expected assignInner depth 0, got depth %d (ok=%v)", depth, ok)
				}

				assignOuterStmt := innerBlock.Statements[5].(statement.ExpressionStatement)
				assignOuter := assignOuterStmt.Expression.(expression.Assign)
				if depth, ok := table.Depth(assignOuter.ID); !ok || depth != 1 {
					t.Errorf("expected assignOuter depth 1, got depth %d (ok=%v)", depth, ok)
				}
			},
		},
		{
			name: "function parameters and closures",
			source: `
fun outer(x) {
    fun inner(y) {
        return x + y;
    }
    return inner;
}
`,
			validate: func(t *testing.T, stmts []statement.Statement, table *BindingTable) {
				outerFn := stmts[0].(statement.Function)
				innerFn := outerFn.Body[0].(statement.Function)
				returnStmt := innerFn.Body[0].(statement.Return)
				binExpr := returnStmt.Value.(expression.Binary)

				xVar := binExpr.Left.(expression.Variable)
				yVar := binExpr.Right.(expression.Variable)

				if depth, ok := table.Depth(yVar.ID); !ok || depth != 0 {
					t.Errorf("expected param y depth 0, got %d (ok=%v)", depth, ok)
				}

				if depth, ok := table.Depth(xVar.ID); !ok || depth != 1 {
					t.Errorf("expected closure param x depth 1, got %d (ok=%v)", depth, ok)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts := parseSource(t, tt.source)
			table := NewBindingTable()
			rslvr := NewResolverBuilder().WithBinding(table).Build()

			if err := rslvr.Resolve(stmts); err != nil {
				t.Fatalf("unexpected resolver error: %v", err)
			}
			tt.validate(t, stmts, table)
		})
	}
}

func TestResolver_SemanticErrors(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expectedErr string
	}{
		{
			name: "self initialization in local scope",
			source: `
{
    var a = a;
}
`,
			expectedErr: "Can't read local variable in its own initializer",
		},
		{
			name: "duplicate variable in local scope",
			source: `
{
    var x = 1;
    var x = 2;
}
`,
			expectedErr: "Already a variable with this name in this scope",
		},
		{
			name: "top-level return",
			source: `
return "bad";
`,
			expectedErr: "Can't return from top-level code",
		},
		{
			name: "top-level empty return",
			source: `
return;
`,
			expectedErr: "Can't return from top-level code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts := parseSource(t, tt.source)
			table := NewBindingTable()
			rslvr := NewResolverBuilder().WithBinding(table).Build()

			err := rslvr.Resolve(stmts)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.expectedErr)
			}
			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("got %q, want error containing %q", err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestResolver_ResolverBuilder(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		setup    func(table *BindingTable) (*Resolver, *mockStep)
		validate func(t *testing.T, custom *mockStep, rslvr *Resolver)
	}{
		{
			name:   "pipeline with binding and custom step",
			source: "var a = 1;",
			setup: func(table *BindingTable) (*Resolver, *mockStep) {
				custom := &mockStep{name: "custom_analyzer"}
				r := NewResolverBuilder().
					WithBinding(table).
					AddStep(custom).
					Build()
				return r, custom
			},
			validate: func(t *testing.T, custom *mockStep, rslvr *Resolver) {
				if rslvr.Pipeline() == nil || len(rslvr.Pipeline().Steps()) != 2 {
					t.Fatalf("expected 2 steps in resolver pipeline, got %v", rslvr.Pipeline())
				}
				if !custom.called {
					t.Errorf("expected custom step to be called")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts := parseSource(t, tt.source)
			table := NewBindingTable()
			rslvr, custom := tt.setup(table)

			if err := rslvr.Resolve(stmts); err != nil {
				t.Fatalf("unexpected resolver run error: %v", err)
			}
			tt.validate(t, custom, rslvr)
		})
	}
}
