package interpreter

import (
	"fmt"

	"glox/pkg/token"
)

type RuntimeError struct {
	Token   token.Token
	Message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.Token.Line, e.Message)
}
