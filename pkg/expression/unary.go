package expression

import (
	"fmt"

	"glox/pkg/token"
)

// ("!" | "-" ) expression
type Unary struct {
	Operator token.Token
	Right    Expression
}

func (u Unary) String() string {
	return fmt.Sprintf("(%s%s)", u.Operator, u.Right)
}
