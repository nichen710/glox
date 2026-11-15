package expression

import (
	"fmt"

	"glox/pkg/token"
)

// expression "." IDENTIFIER "=" expression
type Set struct {
	Object Expression
	Name   token.Token
	Value  Expression
}

func (s Set) String() string {
	return fmt.Sprintf("%s.%s = %s", s.Object, s.Name.Lexeme, s.Value)
}

