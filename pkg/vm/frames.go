package vm

import (
	"glox/pkg/chunk"
	"glox/pkg/value"
)

const MaxFrames = 256

type CallFrame struct {
	closure  *value.Closure
	chunk    *chunk.Chunk
	ip       int
	baseSlot int
}

func (vm *VM) pushFrame(closure *value.Closure, baseSlot int) {
	vm.frames = append(vm.frames, CallFrame{
		closure:  closure,
		chunk:    closure.Fn.Chunk.(*chunk.Chunk),
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
	case value.KindClosure:
		return vm.callClosure(offset, callee.Clos, calleeIdx, argCount)
	case value.KindNative:
		return vm.callNative(offset, callee.Nat, calleeIdx, argCount)
	default:
		return newRuntimeErrorf(offset, "Can only call functions.")
	}
}

func (vm *VM) callClosure(offset int, closure *value.Closure, calleeIdx, argCount int) error {
	if argCount != closure.Fn.Arity {
		return newRuntimeErrorf(offset, "Expected %d arguments but got %d.", closure.Fn.Arity, argCount)
	}
	if len(vm.frames) >= MaxFrames {
		return newRuntimeErrorf(offset, "Stack overflow.")
	}
	vm.pushFrame(closure, calleeIdx+1)
	return nil
}

func (vm *VM) callNative(offset int, nat *value.Native, calleeIdx, argCount int) error {
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
}

func (vm *VM) tailCallValue(offset, argCount int) error {
	calleeIdx := len(vm.stack) - 1 - argCount
	if calleeIdx < 0 {
		return NewStackUnderflowError(offset)
	}

	callee := vm.stack[calleeIdx]
	switch callee.Kind {
	case value.KindClosure:
		return vm.tailCallClosure(offset, callee.Clos, calleeIdx, argCount)
	case value.KindNative:
		return vm.callNative(offset, callee.Nat, calleeIdx, argCount)
	default:
		return newRuntimeErrorf(offset, "Can only call functions.")
	}
}

func (vm *VM) tailCallClosure(offset int, closure *value.Closure, calleeIdx, argCount int) error {
	if argCount != closure.Fn.Arity {
		return newRuntimeErrorf(offset, "Expected %d arguments but got %d.", closure.Fn.Arity, argCount)
	}
	vm.closeUpvalues(vm.frame.baseSlot)
	copy(vm.stack[vm.frame.baseSlot:], vm.stack[calleeIdx+1:calleeIdx+1+argCount])
	vm.stack = vm.stack[:vm.frame.baseSlot+argCount]
	vm.frame.closure = closure
	vm.frame.chunk = closure.Fn.Chunk.(*chunk.Chunk)
	vm.frame.ip = 0
	return nil
}

func (vm *VM) returnFromFrame(offset int) (bool, error) {
	if len(vm.frames) == 1 {
		if len(vm.stack) != 0 {
			return false, newRuntimeErrorf(offset, "Inconsistent stack height at return: %d", len(vm.stack))
		}
		return true, nil
	}

	result, err := vm.pop(offset)
	if err != nil {
		return false, err
	}
	vm.closeUpvalues(vm.frame.baseSlot)
	vm.stack = vm.stack[:vm.frame.baseSlot-1]
	vm.popFrame()
	vm.push(result)
	return false, nil
}
