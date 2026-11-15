package vm

import (
	"glox/pkg/chunk"
	"glox/pkg/value"
)

const MaxFrames = 256

type CallFrame struct {
	closure    *value.Closure
	chunk      *chunk.Chunk
	ip         int
	baseSlot   int
	returnSlot int
}

func (vm *VM) pushFrame(closure *value.Closure, baseSlot, returnSlot int) {
	vm.frames = append(vm.frames, CallFrame{
		closure:    closure,
		chunk:      closure.Fn.Chunk.(*chunk.Chunk),
		ip:         0,
		baseSlot:   baseSlot,
		returnSlot: returnSlot,
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
	case value.KindClass:
		return vm.callClass(offset, callee.Class, calleeIdx, argCount)
	case value.KindBoundMethod:
		vm.stack[calleeIdx] = callee.Bound.Receiver
		return vm.callMethod(offset, callee.Bound.Method, calleeIdx, argCount)
	default:
		return newRuntimeErrorf(offset, "Can only call functions and classes.")
	}
}

func (vm *VM) callClosure(offset int, closure *value.Closure, calleeIdx, argCount int) error {
	if argCount != closure.Fn.Arity {
		return newRuntimeErrorf(offset, "Expected %d arguments but got %d.", closure.Fn.Arity, argCount)
	}
	if len(vm.frames) >= MaxFrames {
		return newRuntimeErrorf(offset, "Stack overflow.")
	}
	vm.pushFrame(closure, calleeIdx+1, calleeIdx)
	return nil
}

func (vm *VM) callMethod(offset int, closure *value.Closure, calleeIdx, argCount int) error {
	if argCount != closure.Fn.Arity {
		return newRuntimeErrorf(offset, "Expected %d arguments but got %d.", closure.Fn.Arity, argCount)
	}
	if len(vm.frames) >= MaxFrames {
		return newRuntimeErrorf(offset, "Stack overflow.")
	}
	vm.pushFrame(closure, calleeIdx, calleeIdx)
	return nil
}

func (vm *VM) callClass(offset int, klass *value.Class, calleeIdx, argCount int) error {
	inst := value.NewInstance(klass)
	vm.stack[calleeIdx] = value.FromInstance(inst)
	if initMethod, ok := klass.Methods["init"]; ok {
		return vm.callMethod(offset, initMethod, calleeIdx, argCount)
	}
	if argCount != 0 {
		return newRuntimeErrorf(offset, "Expected 0 arguments but got %d.", argCount)
	}
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
	case value.KindClass:
		return vm.callClass(offset, callee.Class, calleeIdx, argCount)
	case value.KindBoundMethod:
		return vm.tailCallBoundMethod(offset, callee.Bound, calleeIdx, argCount)
	default:
		return newRuntimeErrorf(offset, "Can only call functions and classes.")
	}
}

func (vm *VM) tailCallClosure(offset int, closure *value.Closure, calleeIdx, argCount int) error {
	if argCount != closure.Fn.Arity {
		return newRuntimeErrorf(offset, "Expected %d arguments but got %d.", closure.Fn.Arity, argCount)
	}
	vm.closeUpvalues(vm.frame.baseSlot)
	copy(vm.stack[vm.frame.returnSlot+1:], vm.stack[calleeIdx+1:calleeIdx+1+argCount])
	vm.stack = vm.stack[:vm.frame.returnSlot+1+argCount]
	vm.frame.baseSlot = vm.frame.returnSlot + 1
	vm.frame.closure = closure
	vm.frame.chunk = closure.Fn.Chunk.(*chunk.Chunk)
	vm.frame.ip = 0
	return nil
}

func (vm *VM) tailCallBoundMethod(offset int, bound *value.BoundMethod, calleeIdx, argCount int) error {
	if argCount != bound.Method.Fn.Arity {
		return newRuntimeErrorf(offset, "Expected %d arguments but got %d.", bound.Method.Fn.Arity, argCount)
	}
	vm.closeUpvalues(vm.frame.baseSlot)
	vm.stack[vm.frame.returnSlot] = bound.Receiver
	copy(vm.stack[vm.frame.returnSlot+1:], vm.stack[calleeIdx+1:calleeIdx+1+argCount])
	vm.stack = vm.stack[:vm.frame.returnSlot+1+argCount]
	vm.frame.baseSlot = vm.frame.returnSlot
	vm.frame.closure = bound.Method
	vm.frame.chunk = bound.Method.Fn.Chunk.(*chunk.Chunk)
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
	vm.stack = vm.stack[:vm.frame.returnSlot]
	vm.popFrame()
	vm.push(result)
	return false, nil
}
