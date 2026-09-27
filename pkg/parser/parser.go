package parser

import (
	"fmt"

	"glox/pkg/expression"
	"glox/pkg/statement"
	"glox/pkg/token"
)

type Parser struct {
	tokens    []token.Token
	current   int
	_error    error // Parsing Phase Detected Error
	factories []StatementFactory
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
		factories: []StatementFactory{
			&PrintFactory{},
			&VarFactory{},
			&BlockFactory{},
			&IfFactory{},
		},
	}
}

// Parser Public Methods
func (p *Parser) Parse() ([]statement.Statement, error) {
	var stmts []statement.Statement
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func (p *Parser) statement() (statement.Statement, error) {
	for _, factory := range p.factories {
		if factory.CanParse(p) {
			return factory.Parse(p)
		}
	}
	return p.expressionStatement()
}

func (p *Parser) expressionStatement() (statement.Statement, error) {
	expr := p.expression()
	if p._error != nil {
		return nil, p._error
	}

	if !p.match(token.SEMICOLON) {
		p.error(p.peek(), "Expected ';' after expression. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	return statement.ExpressionStatement{Expression: expr}, nil
}

// Grammar Rules
func (p *Parser) expression() expression.Expression {
	return p.assignment()
}

func (p *Parser) assignment() expression.Expression {
	expr := p.or()
	if p._error != nil {
		return nil
	}

	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()
		if p._error != nil {
			return nil
		}

		if variable, ok := expr.(expression.Variable); ok {
			return expression.Assign{
				Name:  variable.Name,
				Value: value,
			}
		}

		p.error(equals, "Invalid assignment target.")
		return nil
	}

	return expr
}

func (p *Parser) or() expression.Expression {
	return p.parseLogical(p.and, token.OR)
}

func (p *Parser) and() expression.Expression {
	return p.parseLogical(p.equality, token.AND)
}

func (p *Parser) equality() expression.Expression {
	return p.parseBinary(p.comparison, token.BANG_EQUAL, token.EQUAL_EQUAL)
}

func (p *Parser) comparison() expression.Expression {
	return p.parseBinary(p.term, token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL)
}

func (p *Parser) term() expression.Expression {
	return p.parseBinary(p.factor, token.MINUS, token.PLUS)
}

func (p *Parser) factor() expression.Expression {
	return p.parseBinary(p.unary, token.SLASH, token.STAR, token.PERCENT)
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

	if p.match(token.IDENTIFIER) {
		return expression.Variable{Name: p.previous()}
	}

	p.error(p.peek(), "Expected expression. Got "+p.peek().Lexeme+".")
	return nil
}

// Helper methods
func (p *Parser) parseBinary(nextLevel func() expression.Expression, types ...token.TokenType) expression.Expression {
	expr := nextLevel()
	if p._error != nil {
		return nil
	}

	for p.match(types...) {
		operator := p.previous()
		right := nextLevel()
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

func (p *Parser) parseLogical(nextLevel func() expression.Expression, types ...token.TokenType) expression.Expression {
	expr := nextLevel()
	if p._error != nil {
		return nil
	}

	for p.match(types...) {
		operator := p.previous()
		right := nextLevel()
		if p._error != nil {
			return nil
		}
		expr = expression.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
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

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) error(tok token.Token, message string) {
	p._error = fmt.Errorf("[line %d] Error: %s\n", tok.Line, message)
}
