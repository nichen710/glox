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

type parseFn func(c *Compiler) error

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
		token.LEFT_PAREN:  {(*Compiler).grouping, nil, PREC_NONE, PREC_NONE},
		token.RIGHT_PAREN: {nil, nil, PREC_NONE, PREC_NONE},
		token.NUMBER:      {(*Compiler).value, nil, PREC_NONE, PREC_NONE},
		token.STRING:      {(*Compiler).value, nil, PREC_NONE, PREC_NONE},
		token.NIL:         {(*Compiler).literal, nil, PREC_NONE, PREC_NONE},
		token.TRUE:        {(*Compiler).literal, nil, PREC_NONE, PREC_NONE},
		token.FALSE:       {(*Compiler).literal, nil, PREC_NONE, PREC_NONE},
		token.EOF:         {nil, nil, PREC_NONE, PREC_NONE},
	}
}

func getRule(tokenType token.TokenType) parseRule {
	return rules[tokenType]
}
