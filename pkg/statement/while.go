package statement

import (
	"fmt"

	"glox/pkg/expression"
)

// "while" "(" expression ")" statement
type While struct {
	Condition expression.Expression
	Body      Statement
}

func (w While) String() string {
	return fmt.Sprintf("WHILE %s DO %s", w.Condition, w.Body)
}
