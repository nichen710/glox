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
		{Name: "len", Arity: 1, Fn: lenNative},
		{Name: "push", Arity: 2, Fn: pushNative},
		{Name: "pop", Arity: 1, Fn: popNative},
		{Name: "insert", Arity: 3, Fn: insertNative},
		{Name: "remove", Arity: 2, Fn: removeNative},
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

func lenNative(args []value.Value) (value.Value, error) {
	switch args[0].Kind {
	case value.KindArray:
		return value.Number(float64(len(args[0].Arr.Elements))), nil
	case value.KindString:
		return value.Number(float64(len(args[0].Str))), nil
	default:
		return value.Nil(), fmt.Errorf("len(x) expects an array or string")
	}
}

func pushNative(args []value.Value) (value.Value, error) {
	if args[0].Kind != value.KindArray {
		return value.Nil(), fmt.Errorf("push(arr, val) expects an array as first argument")
	}
	args[0].Arr.Elements = append(args[0].Arr.Elements, args[1])
	return value.Number(float64(len(args[0].Arr.Elements))), nil
}

func popNative(args []value.Value) (value.Value, error) {
	if args[0].Kind != value.KindArray {
		return value.Nil(), fmt.Errorf("pop(arr) expects an array")
	}
	elems := args[0].Arr.Elements
	if len(elems) == 0 {
		return value.Nil(), fmt.Errorf("pop(arr) cannot pop from empty array")
	}
	last := elems[len(elems)-1]
	args[0].Arr.Elements = elems[:len(elems)-1]
	return last, nil
}

func insertNative(args []value.Value) (value.Value, error) {
	if args[0].Kind != value.KindArray {
		return value.Nil(), fmt.Errorf("insert(arr, idx, val) expects an array as first argument")
	}
	elems := args[0].Arr.Elements
	idx, err := parseArrayIndex(args[1], len(elems), true)
	if err != nil {
		return value.Nil(), err
	}
	elems = append(elems, value.Nil())
	copy(elems[idx+1:], elems[idx:])
	elems[idx] = args[2]
	args[0].Arr.Elements = elems
	return value.Number(float64(len(elems))), nil
}

func removeNative(args []value.Value) (value.Value, error) {
	if args[0].Kind != value.KindArray {
		return value.Nil(), fmt.Errorf("remove(arr, idx) expects an array as first argument")
	}
	elems := args[0].Arr.Elements
	idx, err := parseArrayIndex(args[1], len(elems), false)
	if err != nil {
		return value.Nil(), err
	}
	removed := elems[idx]
	args[0].Arr.Elements = append(elems[:idx], elems[idx+1:]...)
	return removed, nil
}
