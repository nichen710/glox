package parser

import (
	"fmt"

	"glox/pkg/token"
)

type ParseError struct {
	Token   token.Token
	Message string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s\n", e.Token.Line, e.Message)
}

func NewParseError(tok token.Token, message string) ParseError {
	return ParseError{
		Token:   tok,
		Message: message,
	}
}
