package compiler

import (
	"glox/pkg/chunk"
	"glox/pkg/token"
)

type PrintCompiler struct{}

func (pc *PrintCompiler) CanCompile(c *Compiler) bool {
	return c.check(token.PRINT)
}

func (pc *PrintCompiler) Compile(c *Compiler) error {
	c.advance()

	if err := c.expression(); err != nil {
		return err
	}
	if err := c.consume(token.SEMICOLON, "Expected ';' after value."); err != nil {
		return err
	}
	c.emitOp(chunk.OP_PRINT)
	return nil
}
