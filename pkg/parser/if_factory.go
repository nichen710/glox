package parser

import (
	"glox/pkg/statement"
	"glox/pkg/token"
)

type IfFactory struct{}

func (f *IfFactory) CanParse(p *Parser) bool {
	return p.check(token.IF)
}

func (f *IfFactory) Parse(p *Parser) (statement.Statement, error) {
	p.advance()

	if !p.match(token.LEFT_PAREN) {
		p._error = NewParseError(p.peek(), "Expected '(' after 'if'. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	condition := p.expression()
	if p._error != nil {
		return nil, p._error
	}

	if !p.match(token.RIGHT_PAREN) {
		p._error = NewParseError(p.peek(), "Expected ')' after if condition. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch statement.Statement
	if p.match(token.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return statement.If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}
