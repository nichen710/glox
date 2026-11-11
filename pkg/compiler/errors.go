package compiler

import (
	"fmt"

	"glox/pkg/token"
)

type CompileError struct {
	Token   token.Token
	Message string
}

func (e *CompileError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.Token.Line, e.Message)
}

func NewCompileError(tok token.Token, message string) *CompileError {
	return &CompileError{
		Token:   tok,
		Message: message,
	}
}

func newCompileErrorf(tok token.Token, format string, args ...any) *CompileError {
	return &CompileError{
		Token:   tok,
		Message: fmt.Sprintf(format, args...),
	}
}
