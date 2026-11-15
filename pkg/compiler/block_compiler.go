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

	for !c.check(token.RIGHT_BRACE) && !c.isAtEnd() {
		if err := c.declaration(); err != nil {
			return err
		}
	}
	if err := c.consume(token.RIGHT_BRACE, "Expected '}' after block."); err != nil {
		return err
	}

	c.endScope()
	return nil
}
