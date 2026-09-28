package expression

import (
	"fmt"

	"glox/pkg/token"
)

// expression ("and" | "or") expression
type Logical struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (l Logical) String() string {
	return fmt.Sprintf("(%s %s %s)", l.Left, l.Operator, l.Right)
}
