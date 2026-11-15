package interpreter

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var stdinReader = bufio.NewReader(os.Stdin)

func (i *Interpreter) registerNatives() {
	natives := []*NativeFunction{
		NewNativeFunction("clock", 0, clockNative),
		NewNativeFunction("random", 2, randomNative),
		NewNativeFunction("sleep", 1, sleepNative),
		NewNativeFunction("input", 0, inputNative),
		NewNativeFunction("toNumber", 1, toNumberNative),
		NewNativeFunction("clear", 0, clearNative),
		NewNativeFunction("len", 1, lenNative),
		NewNativeFunction("push", 2, pushNative),
		NewNativeFunction("pop", 1, popNative),
		NewNativeFunction("insert", 3, insertNative),
		NewNativeFunction("remove", 2, removeNative),
	}

	for _, nat := range natives {
		i.globals.Define(nat.name, nat)
	}
}

func clockNative(_ *Interpreter, _ []any) (any, error) {
	return float64(time.Now().UnixNano()) / 1e9, nil
}

func randomNative(_ *Interpreter, args []any) (any, error) {
	minVal, okMin := args[0].(float64)
	maxVal, okMax := args[1].(float64)
	if !okMin || !okMax {
		return nil, fmt.Errorf("random(min, max) expects two numbers")
	}
	lo, hi := int(minVal), int(maxVal)
	if lo > hi {
		return nil, fmt.Errorf("random(min, max) expects min <= max")
	}
	return float64(lo + rand.Intn(hi-lo+1)), nil
}

func sleepNative(_ *Interpreter, args []any) (any, error) {
	ms, ok := args[0].(float64)
	if !ok {
		return nil, fmt.Errorf("sleep(ms) expects a number")
	}
	if ms > 0 {
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
	return nil, nil
}

func inputNative(_ *Interpreter, _ []any) (any, error) {
	line, err := stdinReader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return nil, nil
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func toNumberNative(_ *Interpreter, args []any) (any, error) {
	if n, ok := args[0].(float64); ok {
		return n, nil
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, nil
	}
	n, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return nil, nil
	}
	return n, nil
}

func clearNative(i *Interpreter, _ []any) (any, error) {
	w := i.out
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprint(w, "\x1b[2J\x1b[H")
	return nil, nil
}

func lenNative(_ *Interpreter, args []any) (any, error) {
	switch v := args[0].(type) {
	case *LoxArray:
		return float64(len(v.Elements)), nil
	case string:
		return float64(len(v)), nil
	default:
		return nil, fmt.Errorf("len(x) expects an array or string")
	}
}

func pushNative(_ *Interpreter, args []any) (any, error) {
	arr, ok := args[0].(*LoxArray)
	if !ok {
		return nil, fmt.Errorf("push(arr, val) expects an array as first argument")
	}
	arr.Elements = append(arr.Elements, args[1])
	return float64(len(arr.Elements)), nil
}

func popNative(_ *Interpreter, args []any) (any, error) {
	arr, ok := args[0].(*LoxArray)
	if !ok {
		return nil, fmt.Errorf("pop(arr) expects an array")
	}
	if len(arr.Elements) == 0 {
		return nil, fmt.Errorf("pop(arr) cannot pop from empty array")
	}
	last := arr.Elements[len(arr.Elements)-1]
	arr.Elements = arr.Elements[:len(arr.Elements)-1]
	return last, nil
}

func insertNative(_ *Interpreter, args []any) (any, error) {
	arr, ok := args[0].(*LoxArray)
	if !ok {
		return nil, fmt.Errorf("insert(arr, idx, val) expects an array as first argument")
	}
	idx, err := parseArrayIndex(args[1], len(arr.Elements), true)
	if err != nil {
		return nil, err
	}
	arr.Elements = append(arr.Elements, nil)
	copy(arr.Elements[idx+1:], arr.Elements[idx:])
	arr.Elements[idx] = args[2]
	return float64(len(arr.Elements)), nil
}

func removeNative(_ *Interpreter, args []any) (any, error) {
	arr, ok := args[0].(*LoxArray)
	if !ok {
		return nil, fmt.Errorf("remove(arr, idx) expects an array as first argument")
	}
	idx, err := parseArrayIndex(args[1], len(arr.Elements), false)
	if err != nil {
		return nil, err
	}
	removed := arr.Elements[idx]
	arr.Elements = append(arr.Elements[:idx], arr.Elements[idx+1:]...)
	return removed, nil
}
