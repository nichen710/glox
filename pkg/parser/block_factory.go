package parser

import (
	"glox/pkg/statement"
	"glox/pkg/token"
)

type BlockFactory struct{}

func (f *BlockFactory) CanParse(p *Parser) bool {
	return p.check(token.LEFT_BRACE)
}

func (f *BlockFactory) Parse(p *Parser) (statement.Statement, error) {
	p.advance()

	var stmts []statement.Statement
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	if !p.match(token.RIGHT_BRACE) {
		p._error = NewParseError(p.peek(), "Expected '}' after block. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	return statement.Block{Statements: stmts}, nil
}
