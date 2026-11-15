package compiler

import (
	"glox/pkg/chunk"
	"glox/pkg/token"
)

type ExpressionStatementCompiler struct{}

func (ec *ExpressionStatementCompiler) CanCompile(_ *Compiler) bool {
	return true
}

func (ec *ExpressionStatementCompiler) Compile(c *Compiler) error {
	if err := c.expression(); err != nil {
		return err
	}
	if err := c.consume(token.SEMICOLON, "Expected ';' after expression."); err != nil {
		return err
	}
	c.emitOp(chunk.OP_POP)
	return nil
}

