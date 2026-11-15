package vm

import (
	"fmt"

	"glox/pkg/chunk"
)

type RuntimeError struct {
	Line    int
	Offset  int
	Message string
}

func (e *RuntimeError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("[line %d | offset %04d] Error: %s", e.Line, e.Offset, e.Message)
	}
	return fmt.Sprintf("[offset %04d] Error: %s", e.Offset, e.Message)
}

func newRuntimeErrorf(offset int, format string, args ...any) *RuntimeError {
	return &RuntimeError{
		Offset:  offset,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewStackUnderflowError(offset int) *RuntimeError {
	return &RuntimeError{
		Offset:  offset,
		Message: "Stack underflow",
	}
}

func NewUndefinedVariableError(offset int, name string) *RuntimeError {
	return &RuntimeError{
		Offset:  offset,
		Message: fmt.Sprintf("Undefined variable '%s'.", name),
	}
}

func NewUnknownOpCodeError(offset int, op chunk.OpCode) *RuntimeError {
	return &RuntimeError{
		Offset:  offset,
		Message: fmt.Sprintf("Unknown opcode: `%d`", byte(op)),
	}
}
