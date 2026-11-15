package expression

import (
	"fmt"

	"glox/pkg/token"
)

// expression "[" expression "]" "=" expression
type IndexSet struct {
	Object  Expression
	Bracket token.Token
	Index   Expression
	Value   Expression
}

func (s IndexSet) String() string {
	return fmt.Sprintf("%s[%s] = %s", s.Object, s.Index, s.Value)
}

