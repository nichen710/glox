package interpreter

import (
	"fmt"
	"strconv"
)

func isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if b, ok := value.(bool); ok {
		return b
	}
	return true
}

func isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func Stringify(value any) string {
	if value == nil {
		return "nil"
	}

	switch v := value.(type) {
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
