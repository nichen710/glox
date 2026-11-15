package vm

import (
	"fmt"
	"math"

	"glox/pkg/value"
)

func (vm *VM) buildArray(offset int, count int) error {
	if len(vm.stack) < count {
		return NewStackUnderflowError(offset)
	}
	start := len(vm.stack) - count
	elems := make([]value.Value, count)
	copy(elems, vm.stack[start:])
	vm.stack = vm.stack[:start]
	vm.push(value.FromArray(value.NewArray(elems)))
	return nil
}

func (vm *VM) getIndex(offset int) error {
	idxVal, err := vm.pop(offset)
	if err != nil {
		return err
	}
	arrVal, err := vm.pop(offset)
	if err != nil {
		return err
	}
	if arrVal.Kind != value.KindArray {
		return newRuntimeErrorf(offset, "Only arrays can be indexed.")
	}
	idx, err := parseArrayIndex(idxVal, len(arrVal.Arr.Elements), false)
	if err != nil {
		return newRuntimeErrorf(offset, "%v", err)
	}
	vm.push(arrVal.Arr.Elements[idx])
	return nil
}

func (vm *VM) setIndex(offset int) error {
	val, err := vm.pop(offset)
	if err != nil {
		return err
	}
	idxVal, err := vm.pop(offset)
	if err != nil {
		return err
	}
	arrVal, err := vm.pop(offset)
	if err != nil {
		return err
	}
	if arrVal.Kind != value.KindArray {
		return newRuntimeErrorf(offset, "Only arrays can be indexed.")
	}
	idx, err := parseArrayIndex(idxVal, len(arrVal.Arr.Elements), false)
	if err != nil {
		return newRuntimeErrorf(offset, "%v", err)
	}
	arrVal.Arr.Elements[idx] = val
	vm.push(val)
	return nil
}

func parseArrayIndex(idxVal value.Value, length int, allowEnd bool) (int, error) {
	num, ok := idxVal.AsNumber()
	if !ok {
		return 0, fmt.Errorf("Array index must be a number.")
	}
	if math.IsNaN(num) || math.IsInf(num, 0) || math.Trunc(num) != num {
		return 0, fmt.Errorf("Array index must be an integer.")
	}
	idx := int(num)
	maxIdx := length - 1
	if allowEnd {
		maxIdx = length
	}
	if idx < 0 || idx > maxIdx {
		return 0, fmt.Errorf("Array index %d out of bounds for length %d.", idx, length)
	}
	return idx, nil
}

