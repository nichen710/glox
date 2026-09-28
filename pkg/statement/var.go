package statement

import (
	"fmt"

	"glox/pkg/expression"
	"glox/pkg/token"
)

// "var" IDENTIFIER ( "=" expression )? ";"
type Var struct {
	Name        token.Token
	Initializer expression.Expression
}

func (v Var) String() string {
	if v.Initializer == nil {
		return fmt.Sprintf("VAR %s = None", v.Name.Lexeme)
	}

	return fmt.Sprintf("VAR %s = %s", v.Name.Lexeme, v.Initializer)
}
