package statement

import (
	"fmt"

	"glox/pkg/expression"
)

// "if" "(" expression ")" statement ( "else" statement )?
type If struct {
	Condition  expression.Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (i If) String() string {
	if i.ElseBranch == nil {
		return fmt.Sprintf("IF %s THEN %s", i.Condition, i.ThenBranch)
	}

	return fmt.Sprintf("IF %s THEN %s ELSE %s", i.Condition, i.ThenBranch, i.ElseBranch)
}
