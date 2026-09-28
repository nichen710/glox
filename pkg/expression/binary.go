package expression

import (
	"fmt"

	"glox/pkg/token"
)

// expression ["+" | "-" | "*" | "/" | "==" | "!=" | "<" | "<=" | ">" | ">="] expression
type Binary struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (b Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left, b.Operator, b.Right)
}
