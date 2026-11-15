package statement

import (
	"fmt"
	"strings"

	"glox/pkg/expression"
	"glox/pkg/token"
)

// "class" IDENTIFIER ("<" IDENTIFIER)? "{" function* "}"
type Class struct {
	Name       token.Token
	Superclass *expression.Variable
	Methods    []Function
}

func (c Class) String() string {
	methods := make([]string, len(c.Methods))
	for i, m := range c.Methods {
		methods[i] = m.String()
	}
	if c.Superclass != nil {
		return fmt.Sprintf("CLASS %s < %s { %s }", c.Name.Lexeme, c.Superclass.Name.Lexeme, strings.Join(methods, "; "))
	}
	return fmt.Sprintf("CLASS %s { %s }", c.Name.Lexeme, strings.Join(methods, "; "))
}
