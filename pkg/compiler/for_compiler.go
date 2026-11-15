package compiler

import (
	"glox/pkg/chunk"
	"glox/pkg/token"
)

type ForCompiler struct{}

func (fc *ForCompiler) CanCompile(c *Compiler) bool {
	return c.check(token.FOR)
}

// for (<inicializador>; <condicion>; <incremento>) <cuerpo>
func (fc *ForCompiler) Compile(c *Compiler) error {
	c.advance()
	c.context.beginScope()

	if err := c.consume(token.LEFT_PAREN, "Expected '(' after 'for'."); err != nil {
		return err
	}
	if err := fc.initializer(c); err != nil {
		return err
	}

	loopStart := len(c.chunk.Code)
	exitJump := -1
	if !c.match(token.SEMICOLON) {
		if err := c.expression(); err != nil {
			return err
		}
		if err := c.consume(token.SEMICOLON, "Expected ';' after loop condition."); err != nil {
			return err
		}
		exitJump = c.emitJump(chunk.OP_JUMP_IF_FALSE)
		c.emitOp(chunk.OP_POP)
	}

	if !c.match(token.RIGHT_PAREN) {
		bodyJump := c.emitJump(chunk.OP_JUMP)
		incrementStart := len(c.chunk.Code)
		if err := c.expression(); err != nil {
			return err
		}
		c.emitOp(chunk.OP_POP)
		if err := c.consume(token.RIGHT_PAREN, "Expected ')' after for clauses."); err != nil {
			return err
		}
		if err := c.emitLoop(loopStart); err != nil {
			return err
		}
		loopStart = incrementStart
		if err := c.patchJump(bodyJump); err != nil {
			return err
		}
	}

	if err := c.statement(); err != nil {
		return err
	}
	if err := c.emitLoop(loopStart); err != nil {
		return err
	}

	if exitJump != -1 {
		if err := c.patchJump(exitJump); err != nil {
			return err
		}
		c.emitOp(chunk.OP_POP)
	}
	c.endScope()
	return nil
}

func (fc *ForCompiler) initializer(c *Compiler) error {
	if c.match(token.SEMICOLON) {
		return nil
	}
	if c.check(token.VAR) {
		return (&VarCompiler{}).Compile(c)
	}
	return (&ExpressionStatementCompiler{}).Compile(c)
}
