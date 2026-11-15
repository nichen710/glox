package value

import "fmt"

type (
	NativeFn func(args []Value) (Value, error)

	Function struct {
		Name  string
		Arity int
		Chunk any
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

func (n *Native) String() string {
	return "<native fn>"
}

