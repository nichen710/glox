package vm

import (
	"fmt"
	"io"
	"os"

	"glox/pkg/chunk"
	"glox/pkg/interpreter"
)

type VM struct {
	out   io.Writer
	chunk *chunk.Chunk
	ip    int
	stack []any
}

func NewVM() *VM {
	return &VM{
		out:   os.Stdout,
		stack: make([]any, 0),
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
			if _, err := vm.pop(offset); err != nil {
				return err
			}
		case chunk.OP_PRINT:
			val, err := vm.pop(offset)
			if err != nil {
				return err
			}
			fmt.Fprintln(vm.out, interpreter.Stringify(val))
		case chunk.OP_RETURN:
			return nil
		default:
			return NewUnknownOpCodeError(offset, op)
		}
	}
	return nil
}

func (vm *VM) readByte() byte {
	b := vm.chunk.Code[vm.ip]
	vm.ip++
	return b
}

func (vm *VM) push(value any) {
	vm.stack = append(vm.stack, value)
}

func (vm *VM) pop(offset int) (any, error) {
	if len(vm.stack) == 0 {
		return nil, NewStackUnderflowError(offset)
	}
	top := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return top, nil
}
