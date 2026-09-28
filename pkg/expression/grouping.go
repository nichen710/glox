package expression

import "fmt"

// "(" expression ")"
type Grouping struct {
	Expression Expression
}

func (g Grouping) String() string {
	return fmt.Sprintf("(%s)", g.Expression)
}
