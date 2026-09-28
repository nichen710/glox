package expression

import (
	"fmt"

	"glox/pkg/token"
)

// IDENTIFIER
type Variable struct {
	ID   uint64
	Name token.Token
}

func NewVariable(name token.Token) Variable {
	return Variable{
		ID:   NextNodeID(),
		Name: name,
	}
}

func (v Variable) String() string {
	return fmt.Sprintf("<%s>", v.Name.Lexeme)
}
