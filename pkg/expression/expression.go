package expression

import (
	"fmt"
	"strings"
	"sync/atomic"

	"glox/pkg/token"
)

var nextNodeID uint64

func NextNodeID() uint64 {
	return atomic.AddUint64(&nextNodeID, 1)
}

func ResetNodeID() {
	atomic.StoreUint64(&nextNodeID, 0)
}

type Expression interface{}

// ("!" | "-" ) expression
type Unary struct {
	Operator token.Token
	Right    Expression
}

func (u Unary) String() string {
	return fmt.Sprintf("(%s%s)", u.Operator, u.Right)
}

// expression ["+" | "-" | "*" | "/" | "==" | "!=" | "<" | "<=" | ">" | ">="] expression
type Binary struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (b Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left, b.Operator, b.Right)
}

// "(" expression ")"
type Grouping struct {
	Expression Expression
}

func (g Grouping) String() string {
	return fmt.Sprintf("(%s)", g.Expression)
}

// NUMBER | STRING | TRUE | FALSE | NIL
type Literal struct {
	Value any
}

func (l Literal) String() string {
	switch v := l.Value.(type) {
	case string:
		return fmt.Sprintf("<\"%s\">", v)
	case float64:
		return fmt.Sprintf("<%f>", v)
	case bool:
		if v {
			return "<TRUE>"
		}
		return "<FALSE>"
	case nil:
		return "<NIL>"
	default:
		return fmt.Sprintf("<%v>", v)
	}
}

// IDENTIFIER
type Variable struct {
	ID   uint64
	Name token.Token
}

func NewVariable(name token.Token) Variable {
	return Variable{
		ID:   NextNodeID(),
		Name: name,
	}
}

func (v Variable) String() string {
	return fmt.Sprintf("<%s>", v.Name.Lexeme)
}

// IDENTIFIER "=" expression
type Assign struct {
	ID    uint64
	Name  token.Token
	Value Expression
}

func NewAssign(name token.Token, value Expression) Assign {
	return Assign{
		ID:    NextNodeID(),
		Name:  name,
		Value: value,
	}
}

func (a Assign) String() string {
	return fmt.Sprintf("%s = %s", a.Name.Lexeme, a.Value)
}

// expression ("and" | "or") expression
type Logical struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (l Logical) String() string {
	return fmt.Sprintf("(%s %s %s)", l.Left, l.Operator, l.Right)
}

// expression "(" arguments? ")"
type Call struct {
	Callee    Expression
	Paren     token.Token
	Arguments []Expression
}

func (c Call) String() string {
	args := make([]string, len(c.Arguments))
	for i, arg := range c.Arguments {
		args[i] = fmt.Sprintf("%v", arg)
	}
	return fmt.Sprintf("fn<%s(%s)>", c.Callee, strings.Join(args, ", "))
}


