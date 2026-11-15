package compiler

import "glox/pkg/token"

type BlockCompiler struct{}

func (bc *BlockCompiler) CanCompile(c *Compiler) bool {
	return c.check(token.LEFT_BRACE)
}

// Al cerrar el bloque sus variables locales se sacan del stack
func (bc *BlockCompiler) Compile(c *Compiler) error {
	c.advance()
	c.context.beginScope()
	if err := c.block(); err != nil {
		return err
	}
	c.endScope()
	return nil
}
