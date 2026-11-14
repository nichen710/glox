package vm

import (
	"fmt"
	"io"
	"math"
	"os"

	"glox/pkg/chunk"
	"glox/pkg/interpreter"
)

type VM struct {
	out     io.Writer
	chunk   *chunk.Chunk
	ip      int
	stack   []any
	globals map[string]any
}

func NewVM() *VM {
	return &VM{
		out:     os.Stdout,
		stack:   make([]any, 0),
		globals: make(map[string]any),
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
			vm.push(nil)
		case chunk.OP_TRUE:
			vm.push(true)
		case chunk.OP_FALSE:
			vm.push(false)
		case chunk.OP_POP:
			_, err = vm.pop(offset)
		case chunk.OP_PRINT:
			var val any
			if val, err = vm.pop(offset); err == nil {
				fmt.Fprintln(vm.out, interpreter.Stringify(val))
			}
		case chunk.OP_NEGATE:
			err = vm.negate(offset)
		case chunk.OP_NOT:
			var val any
			if val, err = vm.pop(offset); err == nil {
				vm.push(!isTruthy(val))
			}
		case chunk.OP_ADD:
			err = vm.add(offset)
		case chunk.OP_SUBTRACT:
			err = vm.numberOp(offset, op, func(a, b float64) any { return a - b })
		case chunk.OP_MULTIPLY:
			err = vm.numberOp(offset, op, func(a, b float64) any { return a * b })
		case chunk.OP_DIVIDE:
			err = vm.divide(offset, op, func(a, b float64) any { return a / b })
		case chunk.OP_MODULO:
			err = vm.divide(offset, op, func(a, b float64) any { return math.Mod(a, b) })
		case chunk.OP_GREATER:
			err = vm.numberOp(offset, op, func(a, b float64) any { return a > b })
		case chunk.OP_LESS:
			err = vm.numberOp(offset, op, func(a, b float64) any { return a < b })
		case chunk.OP_EQUAL:
			var a, b any
			if a, b, err = vm.popOperands(offset); err == nil {
				vm.push(isEqual(a, b))
			}
		case chunk.OP_DEFINE_GLOBAL:
			var val any
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
			var val any
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
			var val any
			slot := vm.readByte()
			if val, err = vm.peek(offset); err == nil {
				vm.stack[slot] = val
			}
		case chunk.OP_JUMP:
			jump := vm.readShort()
			vm.ip += jump
		case chunk.OP_JUMP_IF_FALSE:
			var cond any
			jump := vm.readShort()
			if cond, err = vm.peek(offset); err == nil && !isTruthy(cond) {
				vm.ip += jump
			}
		case chunk.OP_LOOP:
			jump := vm.readShort()
			vm.ip -= jump
		case chunk.OP_RETURN:
			return nil
		default:
			err = NewUnknownOpCodeError(offset, op)
		}

		if err != nil {
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
	num, ok := val.(float64)
	if !ok {
		return newRuntimeErrorf(offset, "Operand of OP_NEGATE must be a number, got: `%v`", val)
	}
	vm.push(-num)
	return nil
}

// + suma numeros o concatena strings
func (vm *VM) add(offset int) error {
	a, b, err := vm.popOperands(offset)
	if err != nil {
		return err
	}

	switch l := a.(type) {
	case float64:
		if r, ok := b.(float64); ok {
			vm.push(l + r)
			return nil
		}
	case string:
		if r, ok := b.(string); ok {
			vm.push(l + r)
			return nil
		}
	}
	return newRuntimeErrorf(offset, "Operands of OP_ADD must be either numbers or strings, got: `%v, %v`", a, b)
}

// Solo aceptan numeros, popean los dos operandos y pushean el resultado
func (vm *VM) numberOp(offset int, op chunk.OpCode, operation func(a, b float64) any) error {
	a, b, err := vm.popNumbers(offset, op)
	if err != nil {
		return err
	}
	vm.push(operation(a, b))
	return nil
}

// No permite dividir por cero
func (vm *VM) divide(offset int, op chunk.OpCode, operation func(a, b float64) any) error {
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
	l, okL := a.(float64)
	r, okR := b.(float64)
	if !okL || !okR {
		return 0, 0, newRuntimeErrorf(offset, "Operands of %s must be numbers, got: `%v, %v`", op, a, b)
	}
	return l, r, nil
}

// El operando derecho esta en el tope del stack, el izquierdo justo abajo
func (vm *VM) popOperands(offset int) (any, any, error) {
	b, err := vm.pop(offset)
	if err != nil {
		return nil, nil, err
	}
	a, err := vm.pop(offset)
	if err != nil {
		return nil, nil, err
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
	return vm.chunk.Constants[vm.readByte()].(string)
}

func (vm *VM) push(value any) {
	vm.stack = append(vm.stack, value)
}

// Devuelve el tope del stack sin sacarlo
func (vm *VM) peek(offset int) (any, error) {
	if len(vm.stack) == 0 {
		return nil, NewStackUnderflowError(offset)
	}
	return vm.stack[len(vm.stack)-1], nil
}

func (vm *VM) pop(offset int) (any, error) {
	if len(vm.stack) == 0 {
		return nil, NewStackUnderflowError(offset)
	}
	top := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return top, nil
}
