package statement

import (
	"fmt"

	"glox/pkg/expression"
)

type Statement interface {
	fmt.Stringer
}

type Print struct {
	Expression expression.Expression
}

func (p Print) String() string {
	return fmt.Sprintf("PRINT %s", p.Expression)
}

type ExpressionStatement struct {
	Expression expression.Expression
}

func (e ExpressionStatement) String() string {
	return fmt.Sprintf("%s", e.Expression)
}
