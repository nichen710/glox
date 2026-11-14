package compiler

import "glox/pkg/token"

// Ordenados de menor a mayor precedencia
type Precedence int

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

type parseFn func(c *Compiler, canAssign bool) error

type parseRule struct {
	prefix           parseFn
	infix            parseFn
	prefixPrecedence Precedence
	infixPrecedence  Precedence
}

// Tabla del Pratt parser
var rules map[token.TokenType]parseRule

func init() {
	rules = map[token.TokenType]parseRule{
		// TokenType: {prefix, infix, prefixPrecedence, infixPrecedence}
		token.LEFT_PAREN:    {(*Compiler).grouping, nil, PREC_NONE, PREC_NONE},
		token.RIGHT_PAREN:   {nil, nil, PREC_NONE, PREC_NONE},
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
		token.IDENTIFIER:    {(*Compiler).variable, nil, PREC_NONE, PREC_NONE},
		token.EOF:           {nil, nil, PREC_NONE, PREC_NONE},
	}
}

func getRule(tokenType token.TokenType) parseRule {
	return rules[tokenType]
}
