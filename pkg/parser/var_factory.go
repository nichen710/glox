package parser

import (
	"glox/pkg/expression"
	"glox/pkg/statement"
	"glox/pkg/token"
)

type VarFactory struct{}

func (f *VarFactory) CanParse(p *Parser) bool {
	return p.check(token.VAR)
}

func (f *VarFactory) Parse(p *Parser) (statement.Statement, error) {
	p.advance()

	if !p.match(token.IDENTIFIER) {
		p.error(p.peek(), "Expected variable name. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}
	name := p.previous()

	var initializer expression.Expression
	if p.match(token.EQUAL) {
		initializer = p.expression()
		if p._error != nil {
			return nil, p._error
		}
	}

	if !p.match(token.SEMICOLON) {
		p.error(p.peek(), "Expected ';' after variable declaration. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	return statement.Var{Name: name, Initializer: initializer}, nil
}
