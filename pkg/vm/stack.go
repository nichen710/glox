package vm

import "glox/pkg/value"

func (vm *VM) push(val value.Value) {
	vm.stack = append(vm.stack, val)
}

// Devuelve el tope del stack sin sacarlo
func (vm *VM) peek(offset int) (value.Value, error) {
	if len(vm.stack) == 0 {
		return value.Nil(), NewStackUnderflowError(offset)
	}
	return vm.stack[len(vm.stack)-1], nil
}

func (vm *VM) pop(offset int) (value.Value, error) {
	if len(vm.stack) == 0 {
		return value.Nil(), NewStackUnderflowError(offset)
	}
	top := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return top, nil
}

// El operando derecho esta en el tope del stack, el izquierdo justo abajo
func (vm *VM) popOperands(offset int) (value.Value, value.Value, error) {
	b, err := vm.pop(offset)
	if err != nil {
		return value.Nil(), value.Nil(), err
	}
	a, err := vm.pop(offset)
	if err != nil {
		return value.Nil(), value.Nil(), err
	}
	return a, b, nil
}

func (vm *VM) popNumbers(offset int, op string) (float64, float64, error) {
	a, b, err := vm.popOperands(offset)
	if err != nil {
		return 0, 0, err
	}
	l, okL := a.AsNumber()
	r, okR := b.AsNumber()
	if !okL || !okR {
		return 0, 0, newRuntimeErrorf(offset, "Operands of %s must be numbers, got: `%v, %v`", op, a, b)
	}
	return l, r, nil
}
