package expression

import (
	"fmt"

	"glox/pkg/token"
)

// IDENTIFIER "=" expression
type Assign struct {
	ID    uint64
	Name  token.Token
	Value Expression
}

func NewAssign(name token.Token, value Expression) Assign {
	return Assign{
		ID:    NextNodeID(),
		Name:  name,
		Value: value,
	}
}

func (a Assign) String() string {
	return fmt.Sprintf("%s = %s", a.Name.Lexeme, a.Value)
}
