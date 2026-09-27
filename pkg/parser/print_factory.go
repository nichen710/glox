package parser

import (
	"glox/pkg/statement"
	"glox/pkg/token"
)

type PrintFactory struct{}

func (f *PrintFactory) CanParse(p *Parser) bool {
	return p.check(token.PRINT)
}

func (f *PrintFactory) Parse(p *Parser) (statement.Statement, error) {
	p.advance()

	expr := p.expression()
	if p._error != nil {
		return nil, p._error
	}

	if !p.match(token.SEMICOLON) {
		p.error(p.peek(), "Expected ';' after value to print. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	return statement.Print{Expression: expr}, nil
}
