package environment

import (
	"strings"
	"testing"

	"glox/pkg/token"
)

func TestEnvironment_DefineAndGet(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		value    any
		expected any
	}{
		{"define and get number", "a", 42.0, 42.0},
		{"define and get string", "b", "hello", "hello"},
		{"define and get nil", "c", nil, nil},
		{"define and get boolean", "d", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			env.Define(tt.varName, tt.value)

			tok := token.Token{Type: token.IDENTIFIER, Lexeme: tt.varName, Line: 1}
			got, err := env.Get(tok)
			if err != nil {
				t.Fatalf("unexpected error getting variable %s: %v", tt.varName, err)
			}
			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEnvironment_Assign(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		initial  any
		updated  any
		expected any
	}{
		{"assign number", "x", 1.0, 2.0, 2.0},
		{"assign string", "name", "alice", "bob", "bob"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			env.Define(tt.varName, tt.initial)

			tok := token.Token{Type: token.IDENTIFIER, Lexeme: tt.varName, Line: 1}
			err := env.Assign(tok, tt.updated)
			if err != nil {
				t.Fatalf("unexpected error assigning variable %s: %v", tt.varName, err)
			}

			got, err := env.Get(tok)
			if err != nil {
				t.Fatalf("unexpected error getting variable %s: %v", tt.varName, err)
			}
			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEnvironment_UndefinedErrors(t *testing.T) {
	tests := []struct {
		name        string
		action      func(env *Environment, tok token.Token) error
		varName     string
		errContains string
	}{
		{
			name: "get undefined variable",
			action: func(env *Environment, tok token.Token) error {
				_, err := env.Get(tok)
				return err
			},
			varName:     "unknown",
			errContains: "Undefined variable 'unknown'.",
		},
		{
			name: "assign undefined variable",
			action: func(env *Environment, tok token.Token) error {
				return env.Assign(tok, 123)
			},
			varName:     "missing",
			errContains: "Undefined variable 'missing'.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()
			tok := token.Token{Type: token.IDENTIFIER, Lexeme: tt.varName, Line: 1}
			err := tt.action(env, tok)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
			}
		})
	}
}

func TestEnvironment_EnclosingScopes(t *testing.T) {
	tests := []struct {
		name      string
		varName   string
		parentVal any
		childVal  any
		assignVal any
		lookupIn  string // "child" or "parent"
		expected  any
		expectErr bool
	}{
		{
			name:      "read variable from enclosing environment",
			varName:   "x",
			parentVal: 10.0,
			lookupIn:  "child",
			expected:  10.0,
		},
		{
			name:      "shadowing variable in child returns child value",
			varName:   "x",
			parentVal: "global",
			childVal:  "local",
			lookupIn:  "child",
			expected:  "local",
		},
		{
			name:      "shadowing variable in child preserves parent value",
			varName:   "x",
			parentVal: "global",
			childVal:  "local",
			lookupIn:  "parent",
			expected:  "global",
		},
		{
			name:      "assigning from child updates parent variable",
			varName:   "x",
			parentVal: 1.0,
			assignVal: 2.0,
			lookupIn:  "parent",
			expected:  2.0,
		},
		{
			name:      "child variable is not visible in parent",
			varName:   "inner",
			childVal:  "value",
			lookupIn:  "parent",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent := NewEnvironment()
			if tt.parentVal != nil {
				parent.Define(tt.varName, tt.parentVal)
			}

			child := NewEnclosedEnvironment(parent)
			if tt.childVal != nil {
				child.Define(tt.varName, tt.childVal)
			}

			tok := token.Token{Type: token.IDENTIFIER, Lexeme: tt.varName, Line: 1}

			if tt.assignVal != nil {
				if err := child.Assign(tok, tt.assignVal); err != nil {
					t.Fatalf("unexpected assign error: %v", err)
				}
			}

			env := child
			if tt.lookupIn == "parent" {
				env = parent
			}

			got, err := env.Get(tok)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil (got %v)", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error getting variable %s: %v", tt.varName, err)
			}
			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}
