package interpreter

import (
	"fmt"

	"glox/pkg/token"
)

type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]*LoxFunction
}

type LoxInstance struct {
	klass  *LoxClass
	fields map[string]any
}

func NewLoxClass(name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		name:       name,
		superclass: superclass,
		methods:    methods,
	}
}

func (c *LoxClass) Arity() int {
	if initializer := c.FindMethod("init"); initializer != nil {
		return initializer.Arity()
	}
	return 0
}

func (c *LoxClass) Call(interpreter *Interpreter, arguments []any) (any, error) {
	instance := NewLoxInstance(c)
	if initializer := c.FindMethod("init"); initializer != nil {
		if _, err := initializer.Bind(instance).Call(interpreter, arguments); err != nil {
			return nil, err
		}
	}
	return instance, nil
}

func (c *LoxClass) FindMethod(name string) *LoxFunction {
	if method, ok := c.methods[name]; ok {
		return method
	}
	if c.superclass != nil {
		return c.superclass.FindMethod(name)
	}
	return nil
}

func (c *LoxClass) String() string {
	return c.name
}

func NewLoxInstance(klass *LoxClass) *LoxInstance {
	return &LoxInstance{
		klass:  klass,
		fields: make(map[string]any),
	}
}

func (i *LoxInstance) Get(name token.Token) (any, error) {
	if val, ok := i.fields[name.Lexeme]; ok {
		return val, nil
	}
	if method := i.klass.FindMethod(name.Lexeme); method != nil {
		return method.Bind(i), nil
	}
	return nil, newRuntimeErrorf(name, "Undefined property '%s'.", name.Lexeme)
}

func (i *LoxInstance) Set(name token.Token, value any) {
	i.fields[name.Lexeme] = value
}

func (i *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", i.klass.name)
}
