package compiler

import (
	"glox/pkg/chunk"
	"glox/pkg/token"
)

type VarCompiler struct{}

func (vc *VarCompiler) CanCompile(c *Compiler) bool {
	return c.check(token.VAR)
}

func (vc *VarCompiler) Compile(c *Compiler) error {
	c.advance()

	if err := c.consume(token.IDENTIFIER, "Expected variable name."); err != nil {
		return err
	}
	name := c.previous()

	if c.context.isGlobalScope() {
		index, err := c.identifierConstant(name)
		if err != nil {
			return err
		}
		if err := vc.initializer(c); err != nil {
			return err
		}
		c.emitOpByte(chunk.OP_DEFINE_GLOBAL, index)
		return nil
	}

	if err := c.declareLocal(name); err != nil {
		return err
	}
	if err := vc.initializer(c); err != nil {
		return err
	}
	c.context.markInitialized()
	return nil
}

// Sin `=` la variable empieza en nil
func (vc *VarCompiler) initializer(c *Compiler) error {
	if c.match(token.EQUAL) {
		if err := c.expression(); err != nil {
			return err
		}
	} else {
		c.emitOp(chunk.OP_NIL)
	}
	return c.consume(token.SEMICOLON, "Expected ';' after variable declaration.")
}
