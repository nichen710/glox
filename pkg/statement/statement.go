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

// "while" "(" expression ")" statement
type While struct {
	Condition expression.Expression
	Body      Statement
}

func (w While) String() string {
	return fmt.Sprintf("WHILE %s DO %s", w.Condition, w.Body)
}

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
	return fmt.Sprintf("FUN fn<%s(%s)> { %s }", f.Name.Lexeme, strings.Join(params, ", "), strings.Join(bodyParts, "; "))
}

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




