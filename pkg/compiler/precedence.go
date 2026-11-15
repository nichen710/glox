package compiler

import "glox/pkg/token"

const (
	PREC_NONE Precedence = iota
	PREC_ASSIGNMENT
	PREC_OR
	PREC_AND
	PREC_EQUALITY
	PREC_COMPARISON
	PREC_TERM
	PREC_FACTOR
	PREC_UNARY
	PREC_CALL
	PREC_PRIMARY
)

// Tabla del Pratt parser
var rules map[token.TokenType]parseRule

type (
	// Ordenados de menor a mayor precedencia
	Precedence int

	parseFn func(c *Compiler, canAssign bool) error

	parseRule struct {
		prefix           parseFn
		infix            parseFn
		prefixPrecedence Precedence
		infixPrecedence  Precedence
	}
)

func init() {
	rules = map[token.TokenType]parseRule{
		// TokenType: {prefix, infix, prefixPrecedence, infixPrecedence}
		token.LEFT_PAREN:    {(*Compiler).grouping, (*Compiler).call, PREC_NONE, PREC_CALL},
		token.RIGHT_PAREN:   {nil, nil, PREC_NONE, PREC_NONE},
		token.LEFT_BRACKET:  {(*Compiler).arrayLiteral, (*Compiler).indexAccess, PREC_NONE, PREC_CALL},
		token.RIGHT_BRACKET: {nil, nil, PREC_NONE, PREC_NONE},
		token.DOT:           {nil, (*Compiler).dot, PREC_NONE, PREC_CALL},
		token.NUMBER:        {(*Compiler).value, nil, PREC_NONE, PREC_NONE},
		token.STRING:        {(*Compiler).value, nil, PREC_NONE, PREC_NONE},
		token.NIL:           {(*Compiler).literal, nil, PREC_NONE, PREC_NONE},
		token.TRUE:          {(*Compiler).literal, nil, PREC_NONE, PREC_NONE},
		token.FALSE:         {(*Compiler).literal, nil, PREC_NONE, PREC_NONE},
		token.MINUS:         {(*Compiler).unary, (*Compiler).binary, PREC_UNARY, PREC_TERM},
		token.PLUS:          {nil, (*Compiler).binary, PREC_NONE, PREC_TERM},
		token.STAR:          {nil, (*Compiler).binary, PREC_NONE, PREC_FACTOR},
		token.SLASH:         {nil, (*Compiler).binary, PREC_NONE, PREC_FACTOR},
		token.PERCENT:       {nil, (*Compiler).binary, PREC_NONE, PREC_FACTOR},
		token.BANG:          {(*Compiler).unary, nil, PREC_UNARY, PREC_NONE},
		token.BANG_EQUAL:    {nil, (*Compiler).binary, PREC_NONE, PREC_EQUALITY},
		token.EQUAL_EQUAL:   {nil, (*Compiler).binary, PREC_NONE, PREC_EQUALITY},
		token.GREATER:       {nil, (*Compiler).binary, PREC_NONE, PREC_COMPARISON},
		token.GREATER_EQUAL: {nil, (*Compiler).binary, PREC_NONE, PREC_COMPARISON},
		token.LESS:          {nil, (*Compiler).binary, PREC_NONE, PREC_COMPARISON},
		token.LESS_EQUAL:    {nil, (*Compiler).binary, PREC_NONE, PREC_COMPARISON},
		token.AND:           {nil, (*Compiler).logicAnd, PREC_NONE, PREC_AND},
		token.OR:            {nil, (*Compiler).logicOr, PREC_NONE, PREC_OR},
		token.THIS:          {(*Compiler).this_, nil, PREC_NONE, PREC_NONE},
		token.SUPER:         {(*Compiler).super_, nil, PREC_NONE, PREC_NONE},
		token.IDENTIFIER:    {(*Compiler).variable, nil, PREC_NONE, PREC_NONE},
		token.EOF:           {nil, nil, PREC_NONE, PREC_NONE},
	}
}

func getRule(tokenType token.TokenType) parseRule {
	return rules[tokenType]
}
