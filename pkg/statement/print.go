package statement

import (
	"fmt"

	"glox/pkg/expression"
)

// "print" expression ";"
type Print struct {
	Expression expression.Expression
}

func (p Print) String() string {
	return fmt.Sprintf("PRINT %s", p.Expression)
}
