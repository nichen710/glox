package parser

import (
	"glox/pkg/expression"
	"glox/pkg/statement"
	"glox/pkg/token"
)

type ForFactory struct{}

func (f *ForFactory) CanParse(p *Parser) bool {
	return p.check(token.FOR)
}

func (f *ForFactory) Parse(p *Parser) (statement.Statement, error) {
	p.advance()

	if !p.match(token.LEFT_PAREN) {
		p._error = NewParseError(p.peek(), "Expected '(' after 'for'. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	var initializer statement.Statement
	var err error
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.check(token.VAR) {
		initializer, err = (&VarFactory{}).Parse(p)
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = (&ExpressionStatementFactory{}).Parse(p)
		if err != nil {
			return nil, err
		}
	}

	var condition expression.Expression
	if !p.check(token.SEMICOLON) {
		condition = p.expression()
		if p._error != nil {
			return nil, p._error
		}
	}
	if !p.match(token.SEMICOLON) {
		p._error = NewParseError(p.peek(), "Expected ';' after loop condition. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	var increment expression.Expression
	if !p.check(token.RIGHT_PAREN) {
		increment = p.expression()
		if p._error != nil {
			return nil, p._error
		}
	}
	if !p.match(token.RIGHT_PAREN) {
		p._error = NewParseError(p.peek(), "Expected ')' after for clauses. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = statement.Block{
			Statements: []statement.Statement{
				body,
				statement.ExpressionStatement{Expression: increment},
			},
		}
	}

	if condition == nil {
		condition = expression.Literal{Value: true}
	}

	body = statement.While{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = statement.Block{
			Statements: []statement.Statement{
				initializer,
				body,
			},
		}
	}

	return body, nil
}
