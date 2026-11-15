package value

import (
	"fmt"
	"strconv"
)

const (
	KindNil Kind = iota
	KindBool
	KindNumber
	KindString
	KindFunction
	KindNative
)

type (
	Kind uint8

	Value struct {
		Kind Kind
		Bool bool
		Num  float64
		Str  string
		Fn   *Function
		Nat  *Native
	}
)

func Nil() Value {
	return Value{Kind: KindNil}
}

func Bool(b bool) Value {
	return Value{Kind: KindBool, Bool: b}
}

func Number(n float64) Value {
	return Value{Kind: KindNumber, Num: n}
}

func String(s string) Value {
	return Value{Kind: KindString, Str: s}
}

func FromFunction(fn *Function) Value {
	return Value{Kind: KindFunction, Fn: fn}
}

func FromNative(nat *Native) Value {
	return Value{Kind: KindNative, Nat: nat}
}

func FromAny(v any) Value {
	switch x := v.(type) {
	case nil:
		return Nil()
	case bool:
		return Bool(x)
	case float64:
		return Number(x)
	case string:
		return String(x)
	default:
		return String(fmt.Sprintf("%v", x))
	}
}

func (v Value) AsNumber() (float64, bool) {
	if v.Kind != KindNumber {
		return 0, false
	}
	return v.Num, true
}

func (v Value) AsString() (string, bool) {
	if v.Kind != KindString {
		return "", false
	}
	return v.Str, true
}

// nil y false son falsos, todo lo demas es verdadero
func (v Value) IsTruthy() bool {
	switch v.Kind {
	case KindNil:
		return false
	case KindBool:
		return v.Bool
	default:
		return true
	}
}

// Valores de distinto tipo nunca son iguales
func (v Value) Equals(other Value) bool {
	if v.Kind != other.Kind {
		return false
	}
	switch v.Kind {
	case KindNil:
		return true
	case KindBool:
		return v.Bool == other.Bool
	case KindNumber:
		return v.Num == other.Num
	case KindString:
		return v.Str == other.Str
	case KindFunction:
		return v.Fn == other.Fn
	case KindNative:
		return v.Nat == other.Nat
	default:
		return false
	}
}

func (v Value) String() string {
	switch v.Kind {
	case KindNil:
		return "nil"
	case KindBool:
		if v.Bool {
			return "true"
		}
		return "false"
	case KindNumber:
		return strconv.FormatFloat(v.Num, 'f', -1, 64)
	case KindString:
		return v.Str
	case KindFunction:
		return v.Fn.String()
	case KindNative:
		return v.Nat.String()
	default:
		return ""
	}
}
