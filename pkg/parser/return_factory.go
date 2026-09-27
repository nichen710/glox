package parser

import (
	"glox/pkg/expression"
	"glox/pkg/statement"
	"glox/pkg/token"
)

type ReturnFactory struct{}

func (f *ReturnFactory) CanParse(p *Parser) bool {
	return p.check(token.RETURN)
}

func (f *ReturnFactory) Parse(p *Parser) (statement.Statement, error) {
	p.advance()
	keyword := p.previous()

	var value expression.Expression
	if !p.check(token.SEMICOLON) {
		value = p.expression()
		if p._error != nil {
			return nil, p._error
		}
	}

	if !p.match(token.SEMICOLON) {
		p.error(p.peek(), "Expected ';' after return value. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	return statement.Return{
		Keyword: keyword,
		Value:   value,
	}, nil
}
