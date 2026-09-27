package expression

import (
	"fmt"

	"glox/pkg/token"
)

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
	Name token.Token
}

func (v Variable) String() string {
	return fmt.Sprintf("<%s>", v.Name.Lexeme)
}

// IDENTIFIER "=" expression
type Assign struct {
	Name  token.Token
	Value Expression
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

