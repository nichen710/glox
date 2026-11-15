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
		if c.fnType == FnInitializer {
			c.emitOpByte(chunk.OP_GET_LOCAL, 0)
			c.emitOp(chunk.OP_RETURN)
			return nil
		}
		c.emitOps(chunk.OP_NIL, chunk.OP_RETURN)
		return nil
	}

	if c.fnType == FnInitializer {
		return NewCompileError(tok, "Can't return a value from an initializer.")
	}

	if err := c.expression(); err != nil {
		return err
	}
	if err := c.consume(token.SEMICOLON, "Expected ';' after return value."); err != nil {
		return err
	}
	if c.lastOpOffset >= 0 && chunk.OpCode(c.chunk.Code[c.lastOpOffset]) == chunk.OP_CALL {
		c.chunk.Code[c.lastOpOffset] = byte(chunk.OP_TAIL_CALL)
	}
	c.emitOp(chunk.OP_RETURN)
	return nil
}
