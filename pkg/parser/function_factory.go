package parser

import (
	"glox/pkg/statement"
	"glox/pkg/token"
)

type FunctionFactory struct{}

func (f *FunctionFactory) CanParse(p *Parser) bool {
	return p.check(token.FUN)
}

func (f *FunctionFactory) Parse(p *Parser) (statement.Statement, error) {
	p.advance()

	if !p.match(token.IDENTIFIER) {
		p._error = NewParseError(p.peek(), "Expected function name. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}
	name := p.previous()

	if !p.match(token.LEFT_PAREN) {
		p._error = NewParseError(p.peek(), "Expected '(' after function name. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	var params []token.Token
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(params) >= 255 {
				p._error = NewParseError(p.peek(), "Can't have more than 255 parameters.")
				return nil, p._error
			}
			if !p.match(token.IDENTIFIER) {
				p._error = NewParseError(p.peek(), "Expected parameter name. Got "+p.peek().Lexeme+".")
				return nil, p._error
			}
			params = append(params, p.previous())
			if !p.match(token.COMMA) {
				break
			}
		}
	}

	if !p.match(token.RIGHT_PAREN) {
		p._error = NewParseError(p.peek(), "Expected ')' after parameters. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	if !p.check(token.LEFT_BRACE) {
		p._error = NewParseError(p.peek(), "Expected '{' before function body. Got "+p.peek().Lexeme+".")
		return nil, p._error
	}

	bodyStmt, err := (&BlockFactory{}).Parse(p)
	if err != nil {
		return nil, err
	}
	bodyBlock := bodyStmt.(statement.Block)

	return statement.Function{
		Name:   name,
		Params: params,
		Body:   bodyBlock.Statements,
	}, nil
}
