package interpreter

import (
	"fmt"
	"math"
	"strings"
)

type LoxArray struct {
	Elements []any
}

func NewLoxArray(elements []any) *LoxArray {
	return &LoxArray{Elements: elements}
}

func (a *LoxArray) String() string {
	parts := make([]string, len(a.Elements))
	for i, el := range a.Elements {
		parts[i] = Stringify(el)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func parseArrayIndex(idxVal any, length int, allowEnd bool) (int, error) {
	num, ok := idxVal.(float64)
	if !ok {
		return 0, fmt.Errorf("Array index must be a number.")
	}
	if math.IsNaN(num) || math.IsInf(num, 0) || math.Trunc(num) != num {
		return 0, fmt.Errorf("Array index must be an integer.")
	}
	idx := int(num)
	maxIdx := length - 1
	if allowEnd {
		maxIdx = length
	}
	if idx < 0 || idx > maxIdx {
		return 0, fmt.Errorf("Array index %d out of bounds for length %d.", idx, length)
	}
	return idx, nil
}

