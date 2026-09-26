package parser

import (
	"fmt"

	"glox/pkg/expression"
	"glox/pkg/token"
)

type Parser struct {
	tokens  []token.Token
	current int
	_error  error // Parsing Phase Detected Error
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() (expression.Expression, error) {
	expr := p.expression()
	if p._error != nil {
		return nil, p._error
	}

	return expr, nil
}

func (p *Parser) expression() expression.Expression {
	return p.equality()
}

func (p *Parser) equality() expression.Expression {
	expr := p.comparison()
	if p._error != nil {
		return nil
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		if p._error != nil {
			return nil
		}
		expr = expression.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) comparison() expression.Expression {
	expr := p.term()
	if p._error != nil {
		return nil
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		if p._error != nil {
			return nil
		}
		expr = expression.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) term() expression.Expression {
	expr := p.factor()
	if p._error != nil {
		return nil
	}

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		if p._error != nil {
			return nil
		}
		expr = expression.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) factor() expression.Expression {
	expr := p.unary()
	if p._error != nil {
		return nil
	}

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		if p._error != nil {
			return nil
		}
		expr = expression.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) unary() expression.Expression {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		if p._error != nil {
			return nil
		}
		return expression.Unary{
			Operator: operator,
			Right:    right,
		}
	}

	return p.primary()
}

func (p *Parser) primary() expression.Expression {
	if p.match(token.FALSE) {
		return expression.Literal{Value: false}
	}
	if p.match(token.TRUE) {
		return expression.Literal{Value: true}
	}
	if p.match(token.NIL) {
		return expression.Literal{Value: nil}
	}

	if p.match(token.NUMBER, token.STRING) {
		return expression.Literal{Value: p.previous().Literal}
	}

	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		if p._error != nil {
			return nil
		}
		if !p.match(token.RIGHT_PAREN) {
			p.error(p.peek(), "Expected ')' after grouping expression. Got "+p.peek().Lexeme+".")
			return nil
		}
		return expression.Grouping{Expression: expr}
	}

	p.error(p.peek(), "Expected expression. Got "+p.peek().Lexeme+".")
	return nil
}

func (p *Parser) error(tok token.Token, message string) {
	p._error = fmt.Errorf("[line %d] Error: %s\n", tok.Line, message)
}

// helper methods

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}
