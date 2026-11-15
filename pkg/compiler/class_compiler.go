package compiler

import (
	"glox/pkg/chunk"
	"glox/pkg/token"
)

type ClassCompiler struct{}

func (cc *ClassCompiler) CanCompile(c *Compiler) bool {
	return c.check(token.CLASS)
}

func (cc *ClassCompiler) Compile(c *Compiler) error {
	c.advance()

	if err := c.consume(token.IDENTIFIER, "Expected class name."); err != nil {
		return err
	}
	className := c.previous()
	nameConst, err := c.identifierConstant(className)
	if err != nil {
		return err
	}

	if !c.context.isGlobalScope() {
		if err := c.declareLocal(className); err != nil {
			return err
		}
		c.context.markInitialized()
	}

	c.emitOpByte(chunk.OP_CLASS, nameConst)
	if c.context.isGlobalScope() {
		c.emitOpByte(chunk.OP_DEFINE_GLOBAL, nameConst)
	}

	classCtx := &ClassContext{enclosing: c.currentClass}
	c.currentClass = classCtx
	defer func() {
		c.currentClass = classCtx.enclosing
	}()

	if c.match(token.LESS) {
		if err := c.consume(token.IDENTIFIER, "Expected superclass name."); err != nil {
			return err
		}
		superName := c.previous()
		if err := c.namedVariable(superName, false); err != nil {
			return err
		}
		if className.Lexeme == superName.Lexeme {
			return NewCompileError(superName, "A class can't inherit from itself.")
		}

		c.context.beginScope()
		c.context.declare("super")
		c.context.markInitialized()

		if err := c.namedVariable(className, false); err != nil {
			return err
		}
		c.emitOp(chunk.OP_INHERIT)
		classCtx.hasSuperclass = true
	}

	if err := c.namedVariable(className, false); err != nil {
		return err
	}
	if err := c.consume(token.LEFT_BRACE, "Expected '{' before class body."); err != nil {
		return err
	}

	for !c.check(token.RIGHT_BRACE) && !c.isAtEnd() {
		if err := cc.compileMethod(c); err != nil {
			return err
		}
	}

	if err := c.consume(token.RIGHT_BRACE, "Expected '}' after class body."); err != nil {
		return err
	}
	c.emitOp(chunk.OP_POP)

	if classCtx.hasSuperclass {
		c.endScope()
	}
	return nil
}

func (cc *ClassCompiler) compileMethod(c *Compiler) error {
	if err := c.consume(token.IDENTIFIER, "Expected method name."); err != nil {
		return err
	}
	methodName := c.previous()
	nameConst, err := c.identifierConstant(methodName)
	if err != nil {
		return err
	}

	fnType := FnMethod
	if methodName.Lexeme == "init" {
		fnType = FnInitializer
	}

	fn, upvalues, err := (&FunctionCompiler{}).compileFunction(c, methodName.Lexeme, fnType)
	if err != nil {
		return err
	}

	if err := c.emitClosure(methodName, fn, upvalues); err != nil {
		return err
	}

	c.emitOpByte(chunk.OP_METHOD, nameConst)
	return nil
}
