package interpreter

type ReturnSignal struct {
	Value any
}

func (r *ReturnSignal) Error() string {
	return "return signal"
}
