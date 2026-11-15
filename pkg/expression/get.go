package expression

import (
	"fmt"

	"glox/pkg/token"
)

// expression "." IDENTIFIER
type Get struct {
	Object Expression
	Name   token.Token
}

func (g Get) String() string {
	return fmt.Sprintf("%s.%s", g.Object, g.Name.Lexeme)
}
