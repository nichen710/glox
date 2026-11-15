package compiler

import (
	"glox/pkg/chunk"
	"glox/pkg/token"
)

type WhileCompiler struct{}

func (wc *WhileCompiler) CanCompile(c *Compiler) bool {
	return c.check(token.WHILE)
}

func (wc *WhileCompiler) Compile(c *Compiler) error {
	c.advance()

	loopStart := len(c.chunk.Code)
	if err := c.condition("Expected '(' after 'while'.", "Expected ')' after while condition."); err != nil {
		return err
	}

	exitJump := c.emitJump(chunk.OP_JUMP_IF_FALSE)
	c.emitOp(chunk.OP_POP)
	if err := c.statement(); err != nil {
		return err
	}
	if err := c.emitLoop(loopStart); err != nil {
		return err
	}

	if err := c.patchJump(exitJump); err != nil {
		return err
	}
	c.emitOp(chunk.OP_POP)
	return nil
}
