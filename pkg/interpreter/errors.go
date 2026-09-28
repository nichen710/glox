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

func NewRuntimeError(tok token.Token, message string) *RuntimeError {
	return &RuntimeError{
		Token:   tok,
		Message: message,
	}
}

func newRuntimeErrorf(tok token.Token, format string, args ...any) *RuntimeError {
	return &RuntimeError{
		Token:   tok,
		Message: fmt.Sprintf(format, args...),
	}
}
