package environment

import (
	"fmt"

	"glox/pkg/token"
)

type Environment struct {
	values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name token.Token) (any, error) {
	if val, ok := e.values[name.Lexeme]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("Undefined variable '%s'.", name.Lexeme)
}

func (e *Environment) Assign(name token.Token, value any) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}
	return fmt.Errorf("Undefined variable '%s'.", name.Lexeme)
}
