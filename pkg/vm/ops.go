package vm

import (
	"math"

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

func (vm *VM) sub(offset int) error {
	a, b, err := vm.popNumbers(offset, "OP_SUBTRACT")
	if err != nil {
		return err
	}
	vm.push(value.Number(a - b))
	return nil
}

func (vm *VM) mul(offset int) error {
	a, b, err := vm.popNumbers(offset, "OP_MULTIPLY")
	if err != nil {
		return err
	}
	vm.push(value.Number(a * b))
	return nil
}

func (vm *VM) div(offset int, op chunk.OpCode) error {
	a, b, err := vm.popNumbers(offset, op.String())
	if err != nil {
		return err
	}
	if b == 0 {
		return newRuntimeErrorf(offset, "Division by %v is not allowed in %s", b, op)
	}
	if op == chunk.OP_MODULO {
		vm.push(value.Number(math.Mod(a, b)))
	} else {
		vm.push(value.Number(a / b))
	}
	return nil
}

func (vm *VM) greater(offset int) error {
	a, b, err := vm.popNumbers(offset, "OP_GREATER")
	if err != nil {
		return err
	}
	vm.push(value.Bool(a > b))
	return nil
}

func (vm *VM) greaterEqual(offset int) error {
	a, b, err := vm.popNumbers(offset, "OP_GREATER_EQUAL")
	if err != nil {
		return err
	}
	vm.push(value.Bool(a >= b))
	return nil
}

func (vm *VM) less(offset int) error {
	a, b, err := vm.popNumbers(offset, "OP_LESS")
	if err != nil {
		return err
	}
	vm.push(value.Bool(a < b))
	return nil
}

func (vm *VM) lessEqual(offset int) error {
	a, b, err := vm.popNumbers(offset, "OP_LESS_EQUAL")
	if err != nil {
		return err
	}
	vm.push(value.Bool(a <= b))
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
