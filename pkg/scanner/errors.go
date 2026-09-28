package scanner

import (
	"fmt"
	"strings"
)

type ScanError struct {
	Line    int
	Message string
}

type ScanErrors []ScanError

func (e ScanError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.Line, e.Message)
}

func NewUnexpectedCharError(line int, c byte) ScanError {
	return ScanError{
		Line:    line,
		Message: fmt.Sprintf("Unexpected character: `%c`", c),
	}
}

func NewUnterminatedStringError(line int, lexeme string) ScanError {
	return ScanError{
		Line:    line,
		Message: fmt.Sprintf("Unterminated string: `%s`", lexeme),
	}
}

func NewInvalidNumberError(line int, text string) ScanError {
	return ScanError{
		Line:    line,
		Message: fmt.Sprintf("Invalid number: `%s`", text),
	}
}

func (e ScanErrors) Error() string {
	var sb strings.Builder
	for i, err := range e {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}
