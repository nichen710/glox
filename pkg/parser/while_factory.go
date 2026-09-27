package parser

import (
	"glox/pkg/statement"
	"glox/pkg/token"
)

type WhileFactory struct{}

func (f *WhileFactory) CanParse(p *Parser) bool {
	return p.check(token.WHILE)
}

func (f *WhileFactory) Parse(p *Parser) (statement.Statement, error) {
	p.advance()

	if !p.match(token.LEFT_PAREN) {
		p.error(p.peek(), "Expected '(' after 'while'. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	condition := p.expression()
	if p._error != nil {
		return nil, p._error
	}

	if !p.match(token.RIGHT_PAREN) {
		p.error(p.peek(), "Expected ')' after while condition. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return statement.While{
		Condition: condition,
		Body:      body,
	}, nil
}
