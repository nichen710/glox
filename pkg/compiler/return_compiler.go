package compiler

import (
	"glox/pkg/chunk"
	"glox/pkg/token"
)

type ReturnCompiler struct{}

func (rc *ReturnCompiler) CanCompile(c *Compiler) bool {
	return c.check(token.RETURN)
}

func (rc *ReturnCompiler) Compile(c *Compiler) error {
	tok := c.advance()
	if c.fnType == FnScript {
		return NewCompileError(tok, "Can't return from top-level code.")
	}

	if c.match(token.SEMICOLON) {
		c.emitOps(chunk.OP_NIL, chunk.OP_RETURN)
		return nil
	}

	if err := c.expression(); err != nil {
		return err
	}
	if err := c.consume(token.SEMICOLON, "Expected ';' after return value."); err != nil {
		return err
	}
	c.emitOp(chunk.OP_RETURN)
	return nil
}

