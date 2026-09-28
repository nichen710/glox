package resolver

import (
	"fmt"

	"glox/pkg/expression"
	"glox/pkg/statement"
	"glox/pkg/token"
)

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunction
)

type FunctionType int

type BindingStep struct {
	table           *BindingTable
	scopes          []map[string]bool
	currentFunction FunctionType
}

func NewBindingStep(table *BindingTable) *BindingStep {
	return &BindingStep{
		table:           table,
		scopes:          make([]map[string]bool, 0),
		currentFunction: FunctionTypeNone,
	}
}

func (s *BindingStep) Name() string {
	return "lexical_binding"
}

func (s *BindingStep) Run(statements []statement.Statement) error {
	for _, stmt := range statements {
		if err := s.resolveStatement(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (s *BindingStep) resolveStatement(stmt statement.Statement) error {
	if stmt == nil {
		return nil
	}
	switch st := stmt.(type) {
	case statement.Block:
		s.beginScope()
		for _, item := range st.Statements {
			if err := s.resolveStatement(item); err != nil {
				return err
			}
		}
		s.endScope()
		return nil
	case statement.Var:
		if err := s.declare(st.Name); err != nil {
			return err
		}
		if st.Initializer != nil {
			if err := s.resolveExpression(st.Initializer); err != nil {
				return err
			}
		}
		s.define(st.Name)
		return nil
	case statement.Function:
		if err := s.declare(st.Name); err != nil {
			return err
		}
		s.define(st.Name)
		return s.resolveFunction(st, FunctionTypeFunction)
	case statement.ExpressionStatement:
		return s.resolveExpression(st.Expression)
	case statement.If:
		if err := s.resolveExpression(st.Condition); err != nil {
			return err
		}
		if err := s.resolveStatement(st.ThenBranch); err != nil {
			return err
		}
		if st.ElseBranch != nil {
			return s.resolveStatement(st.ElseBranch)
		}
		return nil
	case statement.Print:
		return s.resolveExpression(st.Expression)
	case statement.Return:
		if s.currentFunction == FunctionTypeNone {
			return fmt.Errorf("[line %d] Error at '%s': Can't return from top-level code.", st.Keyword.Line, st.Keyword.Lexeme)
		}
		if st.Value != nil {
			return s.resolveExpression(st.Value)
		}
		return nil
	case statement.While:
		if err := s.resolveExpression(st.Condition); err != nil {
			return err
		}
		return s.resolveStatement(st.Body)
	default:
		return fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func (s *BindingStep) resolveExpression(expr expression.Expression) error {
	if expr == nil {
		return nil
	}
	switch e := expr.(type) {
	case expression.Variable:
		if len(s.scopes) > 0 {
			if initialized, ok := s.scopes[len(s.scopes)-1][e.Name.Lexeme]; ok && !initialized {
				return fmt.Errorf("[line %d] Error at '%s': Can't read local variable in its own initializer.", e.Name.Line, e.Name.Lexeme)
			}
		}
		s.resolveLocal(e.ID, e.Name)
		return nil
	case expression.Assign:
		if err := s.resolveExpression(e.Value); err != nil {
			return err
		}
		s.resolveLocal(e.ID, e.Name)
		return nil
	case expression.Binary:
		if err := s.resolveExpression(e.Left); err != nil {
			return err
		}
		return s.resolveExpression(e.Right)
	case expression.Call:
		if err := s.resolveExpression(e.Callee); err != nil {
			return err
		}
		for _, arg := range e.Arguments {
			if err := s.resolveExpression(arg); err != nil {
				return err
			}
		}
		return nil
	case expression.Grouping:
		return s.resolveExpression(e.Expression)
	case expression.Literal:
		return nil
	case expression.Logical:
		if err := s.resolveExpression(e.Left); err != nil {
			return err
		}
		return s.resolveExpression(e.Right)
	case expression.Unary:
		return s.resolveExpression(e.Right)
	default:
		return fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (s *BindingStep) resolveFunction(fn statement.Function, fnType FunctionType) error {
	enclosingFunction := s.currentFunction
	s.currentFunction = fnType
	defer func() {
		s.currentFunction = enclosingFunction
	}()

	s.beginScope()
	for _, param := range fn.Params {
		if err := s.declare(param); err != nil {
			return err
		}
		s.define(param)
	}
	for _, stmt := range fn.Body {
		if err := s.resolveStatement(stmt); err != nil {
			return err
		}
	}
	s.endScope()
	return nil
}

func (s *BindingStep) resolveLocal(nodeID uint64, name token.Token) {
	for i := len(s.scopes) - 1; i >= 0; i-- {
		if _, ok := s.scopes[i][name.Lexeme]; ok {
			s.table.ResolveNamed(nodeID, name.Lexeme, len(s.scopes)-1-i)
			return
		}
	}
}

func (s *BindingStep) beginScope() {
	s.scopes = append(s.scopes, make(map[string]bool))
}

func (s *BindingStep) endScope() {
	s.scopes = s.scopes[:len(s.scopes)-1]
}

func (s *BindingStep) declare(name token.Token) error {
	if len(s.scopes) == 0 {
		return nil
	}
	scope := s.scopes[len(s.scopes)-1]
	if _, exists := scope[name.Lexeme]; exists {
		return fmt.Errorf("[line %d] Error at '%s': Already a variable with this name in this scope.", name.Line, name.Lexeme)
	}
	scope[name.Lexeme] = false
	return nil
}

func (s *BindingStep) define(name token.Token) {
	if len(s.scopes) == 0 {
		return
	}
	s.scopes[len(s.scopes)-1][name.Lexeme] = true
}
