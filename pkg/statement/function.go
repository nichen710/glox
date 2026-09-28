package statement

import (
	"fmt"
	"strings"

	"glox/pkg/token"
)

// "fun" IDENTIFIER "(" parameters? ")" block
type Function struct {
	Name   token.Token
	Params []token.Token
	Body   []Statement
}

func (f Function) String() string {
	params := make([]string, len(f.Params))
	for i, param := range f.Params {
		params[i] = param.Lexeme
	}

	bodyParts := make([]string, len(f.Body))
	for i, stmt := range f.Body {
		bodyParts[i] = stmt.String()
	}

	return fmt.Sprintf("FUN fn<%s(%s)> { %s }",
		f.Name.Lexeme, strings.Join(params, ", "),
		strings.Join(bodyParts, "; "))
}
