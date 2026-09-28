package environment

import "fmt"

type UndefinedVariableError struct {
	Name     string
	Distance *int
}

func (e UndefinedVariableError) Error() string {
	if e.Distance != nil {
		return fmt.Sprintf("Undefined variable '%s' at distance %d.", e.Name, *e.Distance)
	}

	return fmt.Sprintf("Undefined variable '%s'.", e.Name)
}

func NewUndefinedVariableError(name string) UndefinedVariableError {
	return UndefinedVariableError{
		Name: name,
	}
}

func NewUndefinedVariableAtDistanceError(name string, distance int) UndefinedVariableError {
	return UndefinedVariableError{
		Name:     name,
		Distance: &distance,
	}
}
