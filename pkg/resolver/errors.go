package resolver

import (
	"fmt"

	"glox/pkg/token"
)

type ResolveError struct {
	Token   token.Token
	Message string
}

func (e ResolveError) Error() string {
	return fmt.Sprintf("[line %d] Error at '%s': %s", e.Token.Line, e.Token.Lexeme, e.Message)
}

func NewResolveError(tok token.Token, message string) ResolveError {
	return ResolveError{
		Token:   tok,
		Message: message,
	}
}
