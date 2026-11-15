package expression

import (
	"fmt"

	"glox/pkg/token"
)

// "super" "." IDENTIFIER
type Super struct {
	ID      uint64
	Keyword token.Token
	Method  token.Token
}

func NewSuper(keyword, method token.Token) Super {
	return Super{
		ID:      NextNodeID(),
		Keyword: keyword,
		Method:  method,
	}
}

func (s Super) String() string {
	return fmt.Sprintf("super.%s", s.Method.Lexeme)
}
