package statement

import (
	"fmt"

	"glox/pkg/expression"
	"glox/pkg/token"
)

// "return" expression? ";"
type Return struct {
	Keyword token.Token
	Value   expression.Expression
}

func (r Return) String() string {
	if r.Value == nil {
		return "RETURN NIL"
	}

	return fmt.Sprintf("RETURN %s", r.Value)
}
