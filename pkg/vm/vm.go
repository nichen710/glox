package vm

import (
	"errors"
	"fmt"
	"io"
	"os"

	"glox/pkg/chunk"
	"glox/pkg/value"
)

type VM struct {
	out          io.Writer
	frames       []CallFrame
	frame        *CallFrame
	stack        []value.Value
	globals      map[string]value.Value
	openUpvalues *value.Upvalue
}

func NewVM() *VM {
	vm := &VM{
		out:     os.Stdout,
		frames:  make([]CallFrame, 0, MaxFrames),
		stack:   make([]value.Value, 0, 256),
		globals: make(map[string]value.Value),
	}
	vm.registerNatives()
	return vm
}

func (vm *VM) SetWriter(w io.Writer) {
	vm.out = w
}

// Lee la instruccion apuntada por ip, la ejecuta y avanza hasta OP_RETURN
func (vm *VM) Run(c *chunk.Chunk) error {
	vm.stack = vm.stack[:0]
	vm.frames = vm.frames[:0]
	vm.openUpvalues = nil
	vm.pushFrame(&value.Closure{Fn: &value.Function{Chunk: c}}, 0, 0)

	for vm.frame.ip < len(vm.frame.chunk.Code) {
		offset := vm.frame.ip
		op := chunk.OpCode(vm.readByte())

		var err error
		switch op {
		case chunk.OP_CONSTANT:
			vm.push(vm.frame.chunk.Constants[vm.readByte()])
		case chunk.OP_CONSTANT_LONG:
			vm.push(vm.frame.chunk.Constants[vm.readShort()])
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
			err = vm.sub(offset)
		case chunk.OP_MULTIPLY:
			err = vm.mul(offset)
		case chunk.OP_DIVIDE, chunk.OP_MODULO:
			err = vm.div(offset, op)
		case chunk.OP_GREATER:
			err = vm.greater(offset)
		case chunk.OP_GREATER_EQUAL:
			err = vm.greaterEqual(offset)
		case chunk.OP_LESS:
			err = vm.less(offset)
		case chunk.OP_LESS_EQUAL:
			err = vm.lessEqual(offset)
		case chunk.OP_EQUAL:
			var a, b value.Value
			if a, b, err = vm.popOperands(offset); err == nil {
				vm.push(value.Bool(a.Equals(b)))
			}
		case chunk.OP_NOT_EQUAL:
			var a, b value.Value
			if a, b, err = vm.popOperands(offset); err == nil {
				vm.push(value.Bool(!a.Equals(b)))
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
			vm.push(vm.stack[vm.frame.baseSlot+int(vm.readByte())])
		case chunk.OP_SET_LOCAL:
			var val value.Value
			slot := int(vm.readByte())
			if val, err = vm.peek(offset); err == nil {
				vm.stack[vm.frame.baseSlot+slot] = val
			}
		case chunk.OP_JUMP:
			vm.frame.ip += vm.readShort()
		case chunk.OP_JUMP_IF_FALSE:
			var cond value.Value
			jump := vm.readShort()
			if cond, err = vm.peek(offset); err == nil && !cond.IsTruthy() {
				vm.frame.ip += jump
			}
		case chunk.OP_LOOP:
			vm.frame.ip -= vm.readShort()
		case chunk.OP_CALL:
			err = vm.callValue(offset, int(vm.readByte()))
		case chunk.OP_TAIL_CALL:
			err = vm.tailCallValue(offset, int(vm.readByte()))
		case chunk.OP_CLOSURE:
			vm.buildClosure()
		case chunk.OP_GET_UPVALUE:
			vm.getUpvalue()
		case chunk.OP_SET_UPVALUE:
			err = vm.setUpvalue(offset)
		case chunk.OP_CLOSE_UPVALUE:
			vm.closeUpvalues(len(vm.stack) - 1)
			_, err = vm.pop(offset)
		case chunk.OP_CLASS:
			vm.push(value.FromClass(value.NewClass(vm.readName())))
		case chunk.OP_METHOD:
			err = vm.defineMethod(offset, vm.readName())
		case chunk.OP_GET_PROPERTY:
			err = vm.getProperty(offset, vm.readName())
		case chunk.OP_SET_PROPERTY:
			err = vm.setProperty(offset, vm.readName())
		case chunk.OP_INVOKE:
			name := vm.readName()
			argCount := int(vm.readByte())
			err = vm.invoke(offset, name, argCount)
		case chunk.OP_INHERIT:
			err = vm.inherit(offset)
		case chunk.OP_GET_SUPER:
			name := vm.readName()
			var superVal value.Value
			if superVal, err = vm.pop(offset); err == nil {
				err = vm.bindMethod(offset, superVal.Class, name)
			}
		case chunk.OP_SUPER_INVOKE:
			name := vm.readName()
			argCount := int(vm.readByte())
			err = vm.superInvoke(offset, name, argCount)
		case chunk.OP_RETURN:
			var done bool
			if done, err = vm.returnFromFrame(offset); err == nil && done {
				return nil
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
