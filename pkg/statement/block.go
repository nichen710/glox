package statement

import (
	"fmt"
	"strings"
)

// "{" statement* "}"
type Block struct {
	Statements []Statement
}

func (b Block) String() string {
	parts := make([]string, len(b.Statements))
	for i, stmt := range b.Statements {
		parts[i] = stmt.String()
	}

	return fmt.Sprintf("{ %s }", strings.Join(parts, "; "))
}
