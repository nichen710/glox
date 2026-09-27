package parser

import "glox/pkg/statement"

type StatementFactory interface {
	CanParse(p *Parser) bool
	Parse(p *Parser) (statement.Statement, error)
}
