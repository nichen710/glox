package value

import "fmt"

type (
	Class struct {
		Name    string
		Methods map[string]*Closure
	}

	Instance struct {
		Klass  *Class
		Fields map[string]Value
	}

	BoundMethod struct {
		Receiver Value
		Method   *Closure
	}
)

func NewClass(name string) *Class {
	return &Class{
		Name:    name,
		Methods: make(map[string]*Closure),
	}
}

func (c *Class) String() string {
	return c.Name
}

func NewInstance(klass *Class) *Instance {
	return &Instance{
		Klass:  klass,
		Fields: make(map[string]Value),
	}
}

func (i *Instance) String() string {
	return fmt.Sprintf("%s instance", i.Klass.Name)
}

func NewBoundMethod(receiver Value, method *Closure) *BoundMethod {
	return &BoundMethod{
		Receiver: receiver,
		Method:   method,
	}
}

func (b *BoundMethod) String() string {
	return b.Method.String()
}

