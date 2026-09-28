package environment

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

func (e *Environment) Get(name string) (any, error) {
	for env := e; env != nil; env = env.enclosing {
		if val, ok := env.values[name]; ok {
			return val, nil
		}
	}
	return nil, NewUndefinedVariableError(name)
}

func (e *Environment) Assign(name string, value any) error {
	for env := e; env != nil; env = env.enclosing {
		if _, ok := env.values[name]; ok {
			env.values[name] = value
			return nil
		}
	}
	return NewUndefinedVariableError(name)
}

func (e *Environment) ancestor(distance int) *Environment {
	for range distance {
		if e == nil {
			break
		}
		e = e.enclosing
	}
	return e
}

func (e *Environment) GetAt(distance int, name string) (any, error) {
	anc := e.ancestor(distance)
	if anc == nil {
		return nil, NewUndefinedVariableAtDistanceError(name, distance)
	}
	if val, ok := anc.values[name]; ok {
		return val, nil
	}
	return nil, NewUndefinedVariableError(name)
}

func (e *Environment) AssignAt(distance int, name string, value any) error {
	anc := e.ancestor(distance)
	if anc == nil {
		return NewUndefinedVariableAtDistanceError(name, distance)
	}
	if _, ok := anc.values[name]; ok {
		anc.values[name] = value
		return nil
	}
	return NewUndefinedVariableError(name)
}
