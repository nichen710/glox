package parser

import (
	"glox/pkg/statement"
	"glox/pkg/token"
)

type ExpressionStatementFactory struct{}

func (f *ExpressionStatementFactory) CanParse(p *Parser) bool {
	return true
}

func (f *ExpressionStatementFactory) Parse(p *Parser) (statement.Statement, error) {
	expr := p.expression()
	if p._error != nil {
		return nil, p._error
	}

	if !p.match(token.SEMICOLON) {
		p._error = NewParseError(p.peek(), "Expected ';' after expression. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	return statement.ExpressionStatement{Expression: expr}, nil
}
