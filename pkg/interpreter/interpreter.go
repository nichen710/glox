package interpreter

import (
	"fmt"
	"io"
	"math"
	"os"
	"time"

	"glox/pkg/environment"
	"glox/pkg/expression"
	"glox/pkg/statement"
	"glox/pkg/token"
)

type Interpreter struct {
	out         io.Writer
	globals     *environment.Environment
	environment *environment.Environment
}

func NewInterpreter() *Interpreter {
	globals := environment.NewEnvironment()
	globals.Define("clock", NewNativeFunction("clock", 0, func(interpreter *Interpreter, arguments []any) (any, error) {
		return float64(time.Now().UnixNano()) / 1e9, nil
	}))

	return &Interpreter{
		out:         os.Stdout,
		globals:     globals,
		environment: globals,
	}
}

func (i *Interpreter) SetWriter(w io.Writer) {
	i.out = w
}

func (i *Interpreter) Interpret(statements []statement.Statement) (any, error) {
	var last any
	for _, stmt := range statements {
		val, err := i.Execute(stmt)
		if err != nil {
			if ret, ok := err.(*ReturnSignal); ok {
				return ret.Value, nil
			}
			return nil, err
		}
		last = val
	}
	return last, nil
}

func (i *Interpreter) Execute(stmt statement.Statement) (any, error) {
	switch s := stmt.(type) {
	case statement.Print:
		val, err := i.Evaluate(s.Expression)
		if err != nil {
			return nil, err
		}
		out := i.out
		if out == nil {
			out = os.Stdout
		}
		fmt.Fprintln(out, Stringify(val))
		return nil, nil
	case statement.ExpressionStatement:
		return i.Evaluate(s.Expression)
	case statement.Var:
		var val any
		var err error
		if s.Initializer != nil {
			val, err = i.Evaluate(s.Initializer)
			if err != nil {
				return nil, err
			}
		}
		i.environment.Define(s.Name.Lexeme, val)
		return nil, nil
	case statement.Block:
		return i.executeBlock(s.Statements, environment.NewEnclosedEnvironment(i.environment))
	case statement.If:
		condition, err := i.Evaluate(s.Condition)
		if err != nil {
			return nil, err
		}
		if isTruthy(condition) {
			return i.Execute(s.ThenBranch)
		} else if s.ElseBranch != nil {
			return i.Execute(s.ElseBranch)
		}
		return nil, nil
	case statement.While:
		for {
			condition, err := i.Evaluate(s.Condition)
			if err != nil {
				return nil, err
			}
			if !isTruthy(condition) {
				break
			}
			_, err = i.Execute(s.Body)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	case statement.Function:
		function := NewLoxFunction(s, i.environment)
		i.environment.Define(s.Name.Lexeme, function)
		return nil, nil
	case statement.Return:
		var val any
		var err error
		if s.Value != nil {
			val, err = i.Evaluate(s.Value)
			if err != nil {
				return nil, err
			}
		}
		return nil, &ReturnSignal{Value: val}
	default:
		return nil, fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func (i *Interpreter) executeBlock(statements []statement.Statement, env *environment.Environment) (any, error) {
	previous := i.environment
	i.environment = env
	defer func() {
		i.environment = previous
	}()

	var last any
	for _, stmt := range statements {
		val, err := i.Execute(stmt)
		if err != nil {
			return nil, err
		}
		last = val
	}
	return last, nil
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
	case expression.Variable:
		val, err := i.environment.Get(e.Name)
		if err != nil {
			return nil, &RuntimeError{
				Token:   e.Name,
				Message: err.Error(),
			}
		}
		return val, nil
	case expression.Assign:
		val, err := i.Evaluate(e.Value)
		if err != nil {
			return nil, err
		}
		err = i.environment.Assign(e.Name, val)
		if err != nil {
			return nil, &RuntimeError{
				Token:   e.Name,
				Message: err.Error(),
			}
		}
		return val, nil
	case expression.Logical:
		left, err := i.Evaluate(e.Left)
		if err != nil {
			return nil, err
		}
		switch e.Operator.Type {
		case token.OR:
			if isTruthy(left) {
				return left, nil
			}
		case token.AND:
			if !isTruthy(left) {
				return left, nil
			}
		}
		return i.Evaluate(e.Right)
	case expression.Call:
		callee, err := i.Evaluate(e.Callee)
		if err != nil {
			return nil, err
		}

		arguments := make([]any, 0, len(e.Arguments))
		for _, argExpr := range e.Arguments {
			argVal, err := i.Evaluate(argExpr)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, argVal)
		}

		function, ok := callee.(Callable)
		if !ok {
			return nil, &RuntimeError{
				Token:   e.Paren,
				Message: "Can only call functions and classes.",
			}
		}

		if len(arguments) != function.Arity() {
			return nil, &RuntimeError{
				Token:   e.Paren,
				Message: fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments)),
			}
		}

		return function.Call(i, arguments)
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
