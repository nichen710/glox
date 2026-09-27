package statement

import (
	"fmt"

	"glox/pkg/expression"
	"glox/pkg/token"
)

type Statement interface {
	fmt.Stringer
}

type Print struct {
	Expression expression.Expression
}

func (p Print) String() string {
	return fmt.Sprintf("PRINT %s", p.Expression)
}

type ExpressionStatement struct {
	Expression expression.Expression
}

func (e ExpressionStatement) String() string {
	return fmt.Sprintf("%s", e.Expression)
}

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

