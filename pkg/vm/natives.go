package vm

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"glox/pkg/value"
)

var stdinReader = bufio.NewReader(os.Stdin)

func (vm *VM) registerNatives() {
	natives := []*value.Native{
		{Name: "clock", Arity: 0, Fn: clockNative},
		{Name: "random", Arity: 2, Fn: randomNative},
		{Name: "sleep", Arity: 1, Fn: sleepNative},
		{Name: "input", Arity: 0, Fn: inputNative},
		{Name: "toNumber", Arity: 1, Fn: toNumberNative},
		{Name: "clear", Arity: 0, Fn: vm.clearNative},
	}

	for _, nat := range natives {
		vm.globals[nat.Name] = value.FromNative(nat)
	}
}

func clockNative(_ []value.Value) (value.Value, error) {
	return value.Number(float64(time.Now().UnixNano()) / 1e9), nil
}

func randomNative(args []value.Value) (value.Value, error) {
	minVal, okMin := args[0].AsNumber()
	maxVal, okMax := args[1].AsNumber()
	if !okMin || !okMax {
		return value.Nil(), fmt.Errorf("random(min, max) expects two numbers")
	}
	lo, hi := int(minVal), int(maxVal)
	if lo > hi {
		return value.Nil(), fmt.Errorf("random(min, max) expects min <= max")
	}
	return value.Number(float64(lo + rand.Intn(hi-lo+1))), nil
}

func sleepNative(args []value.Value) (value.Value, error) {
	ms, ok := args[0].AsNumber()
	if !ok {
		return value.Nil(), fmt.Errorf("sleep(ms) expects a number")
	}
	if ms > 0 {
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
	return value.Nil(), nil
}

func inputNative(_ []value.Value) (value.Value, error) {
	line, err := stdinReader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return value.Nil(), nil
	}
	return value.String(strings.TrimRight(line, "\r\n")), nil
}

func toNumberNative(args []value.Value) (value.Value, error) {
	if n, ok := args[0].AsNumber(); ok {
		return value.Number(n), nil
	}
	s, ok := args[0].AsString()
	if !ok {
		return value.Nil(), nil
	}
	n, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return value.Nil(), nil
	}
	return value.Number(n), nil
}

func (vm *VM) clearNative(_ []value.Value) (value.Value, error) {
	fmt.Fprint(vm.out, "\x1b[2J\x1b[H")
	return value.Nil(), nil
}
