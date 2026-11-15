package vm

import (
	"glox/pkg/chunk"
	"glox/pkg/value"
)

func (vm *VM) negate(offset int) error {
	val, err := vm.pop(offset)
	if err != nil {
		return err
	}
	num, ok := val.AsNumber()
	if !ok {
		return newRuntimeErrorf(offset, "Operand of OP_NEGATE must be a number, got: `%v`", val)
	}
	vm.push(value.Number(-num))
	return nil
}

// + suma numeros o concatena strings
func (vm *VM) add(offset int) error {
	a, b, err := vm.popOperands(offset)
	if err != nil {
		return err
	}

	if l, okL := a.AsNumber(); okL {
		if r, okR := b.AsNumber(); okR {
			vm.push(value.Number(l + r))
			return nil
		}
	}
	if l, okL := a.AsString(); okL {
		if r, okR := b.AsString(); okR {
			vm.push(value.String(l + r))
			return nil
		}
	}
	return newRuntimeErrorf(offset, "Operands of OP_ADD must be either numbers or strings, got: `%v, %v`", a, b)
}

// Solo aceptan numeros, popean los dos operandos y pushean el resultado
func (vm *VM) numberOp(offset int, op chunk.OpCode, operation func(a, b float64) value.Value) error {
	a, b, err := vm.popNumbers(offset, op.String())
	if err != nil {
		return err
	}
	vm.push(operation(a, b))
	return nil
}

// No permite dividir por cero
func (vm *VM) divide(offset int, op chunk.OpCode, operation func(a, b float64) value.Value) error {
	a, b, err := vm.popNumbers(offset, op.String())
	if err != nil {
		return err
	}
	if b == 0 {
		return newRuntimeErrorf(offset, "Division by %v is not allowed in %s", b, op)
	}
	vm.push(operation(a, b))
	return nil
}

func (vm *VM) readByte() byte {
	b := vm.frame.chunk.Code[vm.frame.ip]
	vm.frame.ip++
	return b
}

func (vm *VM) readShort() int {
	val := vm.frame.chunk.ReadShort(vm.frame.ip)
	vm.frame.ip += 2
	return val
}

func (vm *VM) readName() string {
	name, _ := vm.frame.chunk.Constants[vm.readByte()].AsString()
	return name
}
