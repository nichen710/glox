package vm

import (
	"time"

	"glox/pkg/value"
)

func (vm *VM) registerNatives() {
	vm.globals["clock"] = value.FromNative(&value.Native{
		Name:  "clock",
		Arity: 0,
		Fn:    clockNative,
	})
}

func clockNative(_ []value.Value) (value.Value, error) {
	return value.Number(float64(time.Now().UnixNano()) / 1e9), nil
}
