package statement

import (
	"fmt"

	"glox/pkg/expression"
)

// expression ";"
type ExpressionStatement struct {
	Expression expression.Expression
}

func (e ExpressionStatement) String() string {
	return fmt.Sprintf("%s", e.Expression)
}
