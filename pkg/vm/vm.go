package vm

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"

	"glox/pkg/chunk"
	"glox/pkg/value"
)

type VM struct {
	out     io.Writer
	chunk   *chunk.Chunk
	ip      int
	stack   []value.Value
	globals map[string]value.Value
}

func NewVM() *VM {
	return &VM{
		out:     os.Stdout,
		stack:   make([]value.Value, 0),
		globals: make(map[string]value.Value),
	}
}

func (vm *VM) SetWriter(w io.Writer) {
	vm.out = w
}

// Lee la instruccion apuntada por ip, la ejecuta y avanza hasta OP_RETURN
func (vm *VM) Run(c *chunk.Chunk) error {
	vm.chunk = c
	vm.ip = 0
	vm.stack = vm.stack[:0]

	for vm.ip < len(vm.chunk.Code) {
		offset := vm.ip
		op := chunk.OpCode(vm.readByte())

		var err error
		switch op {
		case chunk.OP_CONSTANT:
			index := vm.readByte()
			vm.push(vm.chunk.Constants[index])
		case chunk.OP_NIL:
			vm.push(value.Nil())
		case chunk.OP_TRUE:
			vm.push(value.Bool(true))
		case chunk.OP_FALSE:
			vm.push(value.Bool(false))
		case chunk.OP_POP:
			_, err = vm.pop(offset)
		case chunk.OP_PRINT:
			var val value.Value
			if val, err = vm.pop(offset); err == nil {
				fmt.Fprintln(vm.out, val.String())
			}
		case chunk.OP_NEGATE:
			err = vm.negate(offset)
		case chunk.OP_NOT:
			var val value.Value
			if val, err = vm.pop(offset); err == nil {
				vm.push(value.Bool(!val.IsTruthy()))
			}
		case chunk.OP_ADD:
			err = vm.add(offset)
		case chunk.OP_SUBTRACT:
			err = vm.numberOp(offset, op, func(a, b float64) value.Value { return value.Number(a - b) })
		case chunk.OP_MULTIPLY:
			err = vm.numberOp(offset, op, func(a, b float64) value.Value { return value.Number(a * b) })
		case chunk.OP_DIVIDE:
			err = vm.divide(offset, op, func(a, b float64) value.Value { return value.Number(a / b) })
		case chunk.OP_MODULO:
			err = vm.divide(offset, op, func(a, b float64) value.Value { return value.Number(math.Mod(a, b)) })
		case chunk.OP_GREATER:
			err = vm.numberOp(offset, op, func(a, b float64) value.Value { return value.Bool(a > b) })
		case chunk.OP_LESS:
			err = vm.numberOp(offset, op, func(a, b float64) value.Value { return value.Bool(a < b) })
		case chunk.OP_EQUAL:
			var a, b value.Value
			if a, b, err = vm.popOperands(offset); err == nil {
				vm.push(value.Bool(a.Equals(b)))
			}
		case chunk.OP_DEFINE_GLOBAL:
			var val value.Value
			name := vm.readName()
			if val, err = vm.pop(offset); err == nil {
				vm.globals[name] = val
			}
		case chunk.OP_GET_GLOBAL:
			name := vm.readName()
			if val, ok := vm.globals[name]; ok {
				vm.push(val)
			} else {
				err = NewUndefinedVariableError(offset, name)
			}
		case chunk.OP_SET_GLOBAL:
			// peek y no pop: la asignacion es una expresion y deja su valor en el stack
			var val value.Value
			name := vm.readName()
			if _, ok := vm.globals[name]; !ok {
				err = NewUndefinedVariableError(offset, name)
			} else if val, err = vm.peek(offset); err == nil {
				vm.globals[name] = val
			}
		case chunk.OP_GET_LOCAL:
			slot := vm.readByte()
			vm.push(vm.stack[slot])
		case chunk.OP_SET_LOCAL:
			var val value.Value
			slot := vm.readByte()
			if val, err = vm.peek(offset); err == nil {
				vm.stack[slot] = val
			}
		case chunk.OP_JUMP:
			jump := vm.readShort()
			vm.ip += jump
		case chunk.OP_JUMP_IF_FALSE:
			var cond value.Value
			jump := vm.readShort()
			if cond, err = vm.peek(offset); err == nil && !cond.IsTruthy() {
				vm.ip += jump
			}
		case chunk.OP_LOOP:
			jump := vm.readShort()
			vm.ip -= jump
		case chunk.OP_RETURN:
			if len(vm.stack) != 0 {
				return newRuntimeErrorf(offset, "Inconsistent stack height at return: %d", len(vm.stack))
			}
			return nil
		default:
			err = NewUnknownOpCodeError(offset, op)
		}

		if err != nil {
			var runtimeErr *RuntimeError
			if errors.As(err, &runtimeErr) && runtimeErr.Line == 0 {
				runtimeErr.Line = vm.chunk.LineAt(runtimeErr.Offset)
			}
			return err
		}
	}
	return nil
}

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
	a, b, err := vm.popNumbers(offset, op)
	if err != nil {
		return err
	}
	vm.push(operation(a, b))
	return nil
}

// No permite dividir por cero
func (vm *VM) divide(offset int, op chunk.OpCode, operation func(a, b float64) value.Value) error {
	a, b, err := vm.popNumbers(offset, op)
	if err != nil {
		return err
	}
	if b == 0 {
		return newRuntimeErrorf(offset, "Division by %v is not allowed in %s", b, op)
	}
	vm.push(operation(a, b))
	return nil
}

func (vm *VM) popNumbers(offset int, op chunk.OpCode) (float64, float64, error) {
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

func (vm *VM) readByte() byte {
	b := vm.chunk.Code[vm.ip]
	vm.ip++
	return b
}

func (vm *VM) readShort() int {
	value := vm.chunk.ReadShort(vm.ip)
	vm.ip += 2
	return value
}

func (vm *VM) readName() string {
	name, _ := vm.chunk.Constants[vm.readByte()].AsString()
	return name
}

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
