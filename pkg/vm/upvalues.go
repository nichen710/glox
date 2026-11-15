package vm

import "glox/pkg/value"

func (vm *VM) buildClosure() {
	fn := vm.frame.chunk.Constants[vm.readByte()].Fn
	closure := &value.Closure{
		Fn:       fn,
		Upvalues: make([]*value.Upvalue, fn.UpvalueCount),
	}
	for i := 0; i < fn.UpvalueCount; i++ {
		isLocal := vm.readByte() == 1
		index := int(vm.readByte())
		if isLocal {
			closure.Upvalues[i] = vm.captureUpvalue(vm.frame.baseSlot + index)
		} else {
			closure.Upvalues[i] = vm.frame.closure.Upvalues[index]
		}
	}
	vm.push(value.FromClosure(closure))
}

func (vm *VM) getUpvalue() {
	uv := vm.frame.closure.Upvalues[vm.readByte()]
	if uv.Slot >= 0 {
		vm.push(vm.stack[uv.Slot])
	} else {
		vm.push(uv.Closed)
	}
}

func (vm *VM) setUpvalue(offset int) error {
	slot := vm.readByte()
	val, err := vm.peek(offset)
	if err != nil {
		return err
	}
	uv := vm.frame.closure.Upvalues[slot]
	if uv.Slot >= 0 {
		vm.stack[uv.Slot] = val
	} else {
		uv.Closed = val
	}
	return nil
}

func (vm *VM) captureUpvalue(slot int) *value.Upvalue {
	var prev *value.Upvalue
	uv := vm.openUpvalues
	for uv != nil && uv.Slot > slot {
		prev = uv
		uv = uv.Next
	}
	if uv != nil && uv.Slot == slot {
		return uv
	}

	created := &value.Upvalue{Slot: slot, Next: uv}
	if prev == nil {
		vm.openUpvalues = created
	} else {
		prev.Next = created
	}
	return created
}

func (vm *VM) closeUpvalues(lastSlot int) {
	for vm.openUpvalues != nil && vm.openUpvalues.Slot >= lastSlot {
		uv := vm.openUpvalues
		uv.Closed = vm.stack[uv.Slot]
		uv.Slot = -1
		vm.openUpvalues = uv.Next
	}
}
