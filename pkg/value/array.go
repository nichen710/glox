package value

import "strings"

type Array struct {
	Elements []Value
}

func NewArray(elements []Value) *Array {
	return &Array{Elements: elements}
}

func (a *Array) String() string {
	parts := make([]string, len(a.Elements))
	for i, el := range a.Elements {
		parts[i] = el.String()
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

