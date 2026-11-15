package expression

import (
	"fmt"

	"glox/pkg/token"
)

// expression "[" expression "]"
type IndexGet struct {
	Object  Expression
	Bracket token.Token
	Index   Expression
}

func (g IndexGet) String() string {
	return fmt.Sprintf("%s[%s]", g.Object, g.Index)
}

