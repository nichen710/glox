package parser

import (
	"glox/pkg/expression"
	"glox/pkg/statement"
	"glox/pkg/token"
)

type ClassFactory struct{}

func (c *ClassFactory) CanParse(p *Parser) bool {
	return p.check(token.CLASS)
}

func (c *ClassFactory) Parse(p *Parser) (statement.Statement, error) {
	p.advance()

	if !p.match(token.IDENTIFIER) {
		p._error = NewParseError(p.peek(), "Expected class name. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}
	name := p.previous()

	var superclass *expression.Variable
	if p.match(token.LESS) {
		if !p.match(token.IDENTIFIER) {
			p._error = NewParseError(p.peek(), "Expected superclass name. Got "+p.peek().Lexeme+".")
			return nil, p._error
		}
		sc := expression.NewVariable(p.previous())
		superclass = &sc
	}

	if !p.match(token.LEFT_BRACE) {
		p._error = NewParseError(p.peek(), "Expected '{' before class body. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	var methods []statement.Function
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		method, err := parseFunction(p, "method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}

	if !p.match(token.RIGHT_BRACE) {
		p._error = NewParseError(p.peek(), "Expected '}' after class body. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	return statement.Class{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}, nil
}

