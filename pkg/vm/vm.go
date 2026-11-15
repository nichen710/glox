package vm

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"time"

	"glox/pkg/chunk"
	"glox/pkg/value"
)

const MaxFrames = 256

type (
	CallFrame struct {
		fn       *value.Function
		chunk    *chunk.Chunk
		ip       int
		baseSlot int
	}

	VM struct {
		out     io.Writer
		frames  []CallFrame
		frame   *CallFrame
		stack   []value.Value
		globals map[string]value.Value
	}
)

func NewVM() *VM {
	globals := make(map[string]value.Value)
	globals["clock"] = value.FromNative(&value.Native{
		Name:  "clock",
		Arity: 0,
		Fn: func(_ []value.Value) (value.Value, error) {
			return value.Number(float64(time.Now().UnixNano()) / 1e9), nil
		},
	})

	return &VM{
		out:     os.Stdout,
		frames:  make([]CallFrame, 0, MaxFrames),
		stack:   make([]value.Value, 0),
		globals: globals,
	}
}

func (vm *VM) SetWriter(w io.Writer) {
	vm.out = w
}

// Lee la instruccion apuntada por ip, la ejecuta y avanza hasta OP_RETURN
func (vm *VM) Run(c *chunk.Chunk) error {
	vm.stack = vm.stack[:0]
	vm.frames = vm.frames[:0]
	vm.pushFrame(&value.Function{Chunk: c}, 0)

	for vm.frame.ip < len(vm.frame.chunk.Code) {
		offset := vm.frame.ip
		op := chunk.OpCode(vm.readByte())

		var err error
		switch op {
		case chunk.OP_CONSTANT:
			index := vm.readByte()
			vm.push(vm.frame.chunk.Constants[index])
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
			slot := int(vm.readByte())
			vm.push(vm.stack[vm.frame.baseSlot+slot])
		case chunk.OP_SET_LOCAL:
			var val value.Value
			slot := int(vm.readByte())
			if val, err = vm.peek(offset); err == nil {
				vm.stack[vm.frame.baseSlot+slot] = val
			}
		case chunk.OP_JUMP:
			jump := vm.readShort()
			vm.frame.ip += jump
		case chunk.OP_JUMP_IF_FALSE:
			var cond value.Value
			jump := vm.readShort()
			if cond, err = vm.peek(offset); err == nil && !cond.IsTruthy() {
				vm.frame.ip += jump
			}
		case chunk.OP_LOOP:
			jump := vm.readShort()
			vm.frame.ip -= jump
		case chunk.OP_CALL:
			argCount := int(vm.readByte())
			err = vm.callValue(offset, argCount)
		case chunk.OP_RETURN:
			if len(vm.frames) == 1 {
				if len(vm.stack) != 0 {
					return vm.withLine(newRuntimeErrorf(offset, "Inconsistent stack height at return: %d", len(vm.stack)))
				}
				return nil
			}
			var result value.Value
			if result, err = vm.pop(offset); err == nil {
				vm.stack = vm.stack[:vm.frame.baseSlot-1]
				vm.popFrame()
				vm.push(result)
			}
		default:
			err = NewUnknownOpCodeError(offset, op)
		}

		if err != nil {
			return vm.withLine(err)
		}
	}
	return nil
}

func (vm *VM) withLine(err error) error {
	var runtimeErr *RuntimeError
	if errors.As(err, &runtimeErr) && runtimeErr.Line == 0 && vm.frame != nil {
		runtimeErr.Line = vm.frame.chunk.LineAt(runtimeErr.Offset)
	}
	return err
}

func (vm *VM) pushFrame(fn *value.Function, baseSlot int) {
	vm.frames = append(vm.frames, CallFrame{
		fn:       fn,
		chunk:    fn.Chunk.(*chunk.Chunk),
		ip:       0,
		baseSlot: baseSlot,
	})
	vm.frame = &vm.frames[len(vm.frames)-1]
}

func (vm *VM) popFrame() {
	vm.frames = vm.frames[:len(vm.frames)-1]
	if len(vm.frames) > 0 {
		vm.frame = &vm.frames[len(vm.frames)-1]
	} else {
		vm.frame = nil
	}
}

func (vm *VM) callValue(offset, argCount int) error {
	calleeIdx := len(vm.stack) - 1 - argCount
	if calleeIdx < 0 {
		return NewStackUnderflowError(offset)
	}

	callee := vm.stack[calleeIdx]
	switch callee.Kind {
	case value.KindFunction:
		fn := callee.Fn
		if argCount != fn.Arity {
			return newRuntimeErrorf(offset, "Expected %d arguments but got %d.", fn.Arity, argCount)
		}
		if len(vm.frames) >= MaxFrames {
			return newRuntimeErrorf(offset, "Stack overflow.")
		}
		vm.pushFrame(fn, calleeIdx+1)
		return nil
	case value.KindNative:
		nat := callee.Nat
		if argCount != nat.Arity {
			return newRuntimeErrorf(offset, "Expected %d arguments but got %d.", nat.Arity, argCount)
		}
		res, err := nat.Fn(vm.stack[calleeIdx+1:])
		if err != nil {
			return newRuntimeErrorf(offset, "%v", err)
		}
		vm.stack = vm.stack[:calleeIdx]
		vm.push(res)
		return nil
	default:
		return newRuntimeErrorf(offset, "Can only call functions.")
	}
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
