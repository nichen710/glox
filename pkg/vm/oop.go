package vm

import "glox/pkg/value"

func (vm *VM) defineMethod(offset int, name string) error {
	methodVal, err := vm.pop(offset)
	if err != nil {
		return err
	}
	klassVal, err := vm.peek(offset)
	if err != nil {
		return err
	}
	klassVal.Class.Methods[name] = methodVal.Clos
	return nil
}

func (vm *VM) getProperty(offset int, name string) error {
	obj, err := vm.peek(offset)
	if err != nil {
		return err
	}
	if obj.Kind != value.KindInstance {
		return newRuntimeErrorf(offset, "Only instances have properties.")
	}
	if val, ok := obj.Inst.Fields[name]; ok {
		_, _ = vm.pop(offset)
		vm.push(val)
		return nil
	}
	return vm.bindMethod(offset, obj.Inst.Klass, name)
}

func (vm *VM) setProperty(offset int, name string) error {
	if len(vm.stack) < 2 {
		return NewStackUnderflowError(offset)
	}
	obj := vm.stack[len(vm.stack)-2]
	if obj.Kind != value.KindInstance {
		return newRuntimeErrorf(offset, "Only instances have fields.")
	}
	val, err := vm.pop(offset)
	if err != nil {
		return err
	}
	obj.Inst.Fields[name] = val
	_, _ = vm.pop(offset)
	vm.push(val)
	return nil
}

func (vm *VM) inherit(offset int) error {
	if len(vm.stack) < 2 {
		return NewStackUnderflowError(offset)
	}
	subVal := vm.stack[len(vm.stack)-1]
	superVal := vm.stack[len(vm.stack)-2]
	if superVal.Kind != value.KindClass {
		return newRuntimeErrorf(offset, "Superclass must be a class.")
	}
	for name, method := range superVal.Class.Methods {
		subVal.Class.Methods[name] = method
	}
	_, _ = vm.pop(offset)
	return nil
}

func (vm *VM) bindMethod(offset int, klass *value.Class, name string) error {
	method, ok := klass.Methods[name]
	if !ok {
		return newRuntimeErrorf(offset, "Undefined property '%s'.", name)
	}
	receiver, err := vm.pop(offset)
	if err != nil {
		return err
	}
	vm.push(value.FromBoundMethod(value.NewBoundMethod(receiver, method)))
	return nil
}

func (vm *VM) invoke(offset int, name string, argCount int) error {
	calleeIdx := len(vm.stack) - 1 - argCount
	if calleeIdx < 0 {
		return NewStackUnderflowError(offset)
	}
	receiver := vm.stack[calleeIdx]
	if receiver.Kind != value.KindInstance {
		return newRuntimeErrorf(offset, "Only instances have methods.")
	}
	if field, ok := receiver.Inst.Fields[name]; ok {
		vm.stack[calleeIdx] = field
		return vm.callValue(offset, argCount)
	}
	return vm.invokeFromClass(offset, receiver.Inst.Klass, name, calleeIdx, argCount)
}

func (vm *VM) superInvoke(offset int, name string, argCount int) error {
	superVal, err := vm.pop(offset)
	if err != nil {
		return err
	}
	calleeIdx := len(vm.stack) - 1 - argCount
	if calleeIdx < 0 {
		return NewStackUnderflowError(offset)
	}
	return vm.invokeFromClass(offset, superVal.Class, name, calleeIdx, argCount)
}

func (vm *VM) invokeFromClass(offset int, klass *value.Class, name string, calleeIdx, argCount int) error {
	method, ok := klass.Methods[name]
	if !ok {
		return newRuntimeErrorf(offset, "Undefined property '%s'.", name)
	}
	return vm.callMethod(offset, method, calleeIdx, argCount)
}
