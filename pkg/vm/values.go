package vm

// nil y false son falsos, todo lo demas es verdadero
func isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if b, ok := value.(bool); ok {
		return b
	}
	return true
}

// Valores de distinto tipo nunca son iguales
func isEqual(a, b any) bool {
	return a == b
}
