package value

import "fmt"

type (
	NativeFn func(args []Value) (Value, error)

	Function struct {
		Name         string
		Arity        int
		UpvalueCount int
		Chunk        any
	}

	Upvalue struct {
		Slot   int
		Closed Value
		Next   *Upvalue
	}

	Closure struct {
		Fn       *Function
		Upvalues []*Upvalue
	}

	Native struct {
		Name  string
		Arity int
		Fn    NativeFn
	}
)

func (f *Function) String() string {
	if f.Name == "" {
		return "<script>"
	}
	return fmt.Sprintf("<fn %s>", f.Name)
}

func (c *Closure) String() string {
	return c.Fn.String()
}

func (n *Native) String() string {
	return "<native fn>"
}
