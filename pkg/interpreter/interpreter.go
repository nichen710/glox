package interpreter

import (
	"fmt"
	"math"

	"glox/pkg/expression"
	"glox/pkg/token"
)

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(expr expression.Expression) (any, error) {
	return i.Evaluate(expr)
}

func (i *Interpreter) Evaluate(expr expression.Expression) (any, error) {
	switch e := expr.(type) {
	case expression.Literal:
		return e.Value, nil
	case expression.Grouping:
		return i.Evaluate(e.Expression)
	case expression.Unary:
		return i.evaluateUnary(e)
	case expression.Binary:
		return i.evaluateBinary(e)
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (i *Interpreter) evaluateUnary(u expression.Unary) (any, error) {
	right, err := i.Evaluate(u.Right)
	if err != nil {
		return nil, err
	}

	switch u.Operator.Type {
	case token.MINUS:
		val, ok := right.(float64)
		if !ok {
			return nil, &RuntimeError{
				Token:   u.Operator,
				Message: fmt.Sprintf("Operand of - must be a number, got: `-%v`", right),
			}
		}
		return -val, nil
	case token.BANG:
		return !isTruthy(right), nil
	default:
		return nil, &RuntimeError{
			Token:   u.Operator,
			Message: fmt.Sprintf("Unknown unary operator: `%s`", u.Operator.Lexeme),
		}
	}
}

func (i *Interpreter) evaluateBinary(b expression.Binary) (any, error) {
	left, err := i.Evaluate(b.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.Evaluate(b.Right)
	if err != nil {
		return nil, err
	}

	switch b.Operator.Type {
	case token.PLUS:
		switch l := left.(type) {
		case float64:
			if r, ok := right.(float64); ok {
				return l + r, nil
			}
		case string:
			if r, ok := right.(string); ok {
				return l + r, nil
			}
		}
		return nil, &RuntimeError{
			Token:   b.Operator,
			Message: fmt.Sprintf("Operands of + must be either numbers or strings, got: `%v + %v`", left, right),
		}
	case token.MINUS:
		l, r, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l - r, nil
	case token.STAR:
		l, r, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l * r, nil
	case token.SLASH:
		l, r, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil, err
		}
		if r == 0 {
			return nil, &RuntimeError{
				Token:   b.Operator,
				Message: fmt.Sprintf("Division by %v is not allowed", r),
			}
		}
		return l / r, nil
	case token.PERCENT:
		l, r, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil, err
		}
		if r == 0 {
			return nil, &RuntimeError{
				Token:   b.Operator,
				Message: fmt.Sprintf("Modulo by %v is not allowed", r),
			}
		}
		return math.Mod(l, r), nil
	case token.GREATER:
		l, r, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l > r, nil
	case token.GREATER_EQUAL:
		l, r, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l >= r, nil
	case token.LESS:
		l, r, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l < r, nil
	case token.LESS_EQUAL:
		l, r, err := i.checkNumberOperands(b.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l <= r, nil
	case token.EQUAL_EQUAL:
		return isEqual(left, right), nil
	case token.BANG_EQUAL:
		return !isEqual(left, right), nil
	default:
		return nil, &RuntimeError{
			Token:   b.Operator,
			Message: fmt.Sprintf("Unknown binary operator: `%s`", b.Operator.Lexeme),
		}
	}
}

func (i *Interpreter) checkNumberOperands(op token.Token, left, right any) (float64, float64, error) {
	l, okL := left.(float64)
	r, okR := right.(float64)
	if !okL || !okR {
		return 0, 0, &RuntimeError{
			Token:   op,
			Message: fmt.Sprintf("Operands of %s must be numbers, got: `%v %s %v`", op.Lexeme, left, op.Lexeme, right),
		}
	}
	return l, r, nil
}
