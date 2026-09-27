package interpreter

import (
	"fmt"
	"strings"

	"glox/pkg/environment"
	"glox/pkg/statement"
)

type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []any) (any, error)
	String() string
}

type ReturnSignal struct {
	Value any
}

func (r *ReturnSignal) Error() string {
	return "return signal"
}

type LoxFunction struct {
	declaration statement.Function
	closure     *environment.Environment
}

func NewLoxFunction(declaration statement.Function, closure *environment.Environment) *LoxFunction {
	return &LoxFunction{
		declaration: declaration,
		closure:     closure,
	}
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []any) (any, error) {
	env := environment.NewEnclosedEnvironment(f.closure)
	for i, param := range f.declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	_, err := interpreter.executeBlock(f.declaration.Body, env)
	if err != nil {
		if ret, ok := err.(*ReturnSignal); ok {
			return ret.Value, nil
		}
		return nil, err
	}

	return nil, nil
}

func (f *LoxFunction) String() string {
	params := make([]string, len(f.declaration.Params))
	for i, param := range f.declaration.Params {
		params[i] = param.Lexeme
	}
	return fmt.Sprintf("<fn %s(%s)>", f.declaration.Name.Lexeme, strings.Join(params, ", "))
}

type NativeFunction struct {
	arity int
	name  string
	call  func(interpreter *Interpreter, arguments []any) (any, error)
}

func NewNativeFunction(name string, arity int, call func(interpreter *Interpreter, arguments []any) (any, error)) *NativeFunction {
	return &NativeFunction{
		name:  name,
		arity: arity,
		call:  call,
	}
}

func (n *NativeFunction) Arity() int {
	return n.arity
}

func (n *NativeFunction) Call(interpreter *Interpreter, arguments []any) (any, error) {
	return n.call(interpreter, arguments)
}

func (n *NativeFunction) String() string {
	return fmt.Sprintf("<native fn %s>", n.name)
}
