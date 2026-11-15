package expression

import (
	"fmt"
	"strings"
)

// "[" (expression ("," expression)*)? "]"
type Array struct {
	Elements []Expression
}

func (a Array) String() string {
	parts := make([]string, len(a.Elements))
	for i, el := range a.Elements {
		parts[i] = fmt.Sprintf("%v", el)
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
}

