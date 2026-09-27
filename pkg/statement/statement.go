package statement

import (
	"fmt"
	"strings"

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


