package expression

import (
	"fmt"
	"strings"

	"glox/pkg/token"
)

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
