package compiler

import (
	"glox/pkg/chunk"
	"glox/pkg/token"
)

type IfCompiler struct{}

func (ic *IfCompiler) CanCompile(c *Compiler) bool {
	return c.check(token.IF)
}

func (ic *IfCompiler) Compile(c *Compiler) error {
	c.advance()

	if err := c.condition("Expected '(' after 'if'.", "Expected ')' after if condition."); err != nil {
		return err
	}

	thenJump := c.emitJump(chunk.OP_JUMP_IF_FALSE)
	c.emitOp(chunk.OP_POP)
	if err := c.statement(); err != nil {
		return err
	}
	elseJump := c.emitJump(chunk.OP_JUMP)

	if err := c.patchJump(thenJump); err != nil {
		return err
	}
	c.emitOp(chunk.OP_POP)
	if c.match(token.ELSE) {
		if err := c.statement(); err != nil {
			return err
		}
	}
	return c.patchJump(elseJump)
}

