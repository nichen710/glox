package environment

import (
	"fmt"

	"glox/pkg/token"
)

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		values:    make(map[string]any),
	}
}

func NewEnclosedEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name token.Token) (any, error) {
	if val, ok := e.values[name.Lexeme]; ok {
		return val, nil
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, fmt.Errorf("Undefined variable '%s'.", name.Lexeme)
}

func (e *Environment) Assign(name token.Token, value any) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return fmt.Errorf("Undefined variable '%s'.", name.Lexeme)
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		if env == nil {
			return nil
		}
		env = env.enclosing
	}
	return env
}

func (e *Environment) GetAt(distance int, name string) (any, error) {
	anc := e.ancestor(distance)
	if anc == nil {
		return nil, fmt.Errorf("Undefined variable '%s' at distance %d.", name, distance)
	}
	if val, ok := anc.values[name]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("Undefined variable '%s'.", name)
}

func (e *Environment) AssignAt(distance int, name string, value any) error {
	anc := e.ancestor(distance)
	if anc == nil {
		return fmt.Errorf("Undefined variable '%s' at distance %d.", name, distance)
	}
	if _, ok := anc.values[name]; ok {
		anc.values[name] = value
		return nil
	}
	return fmt.Errorf("Undefined variable '%s'.", name)
}

