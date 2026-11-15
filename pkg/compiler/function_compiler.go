package compiler

import (
	"glox/pkg/chunk"
	"glox/pkg/token"
	"glox/pkg/value"
)

const MaxArgs = 255

type FunctionCompiler struct{}

func (fc *FunctionCompiler) CanCompile(c *Compiler) bool {
	return c.check(token.FUN)
}

func (fc *FunctionCompiler) Compile(c *Compiler) error {
	c.advance()

	if err := c.consume(token.IDENTIFIER, "Expected function name."); err != nil {
		return err
	}
	name := c.previous()

	var globalIndex byte
	if c.context.isGlobalScope() {
		var err error
		globalIndex, err = c.identifierConstant(name)
		if err != nil {
			return err
		}
	} else {
		if err := c.declareLocal(name); err != nil {
			return err
		}
		// Se marca inicializada antes del cuerpo para permitir recursion
		c.context.markInitialized()
	}

	fn, err := fc.compileFunction(c, name.Lexeme)
	if err != nil {
		return err
	}

	fnIndex, err := c.makeConstant(name, value.FromFunction(fn))
	if err != nil {
		return err
	}
	c.emitOpByte(chunk.OP_CONSTANT, fnIndex)

	if c.context.isGlobalScope() {
		c.emitOpByte(chunk.OP_DEFINE_GLOBAL, globalIndex)
	}
	return nil
}

func (fc *FunctionCompiler) compileFunction(c *Compiler, name string) (*value.Function, error) {
	child := c.newChild(FnFunction)
	child.context.beginScope()

	if err := child.consume(token.LEFT_PAREN, "Expected '(' after function name."); err != nil {
		return nil, err
	}

	arity := 0
	if !child.check(token.RIGHT_PAREN) {
		for {
			if arity >= MaxArgs {
				return nil, NewCompileError(child.peek(), "Can't have more than 255 parameters.")
			}
			if err := child.consume(token.IDENTIFIER, "Expected parameter name."); err != nil {
				return nil, err
			}
			if err := child.declareLocal(child.previous()); err != nil {
				return nil, err
			}
			child.context.markInitialized()
			arity++

			if !child.match(token.COMMA) {
				break
			}
		}
	}

	if err := child.consume(token.RIGHT_PAREN, "Expected ')' after parameters."); err != nil {
		return nil, err
	}
	if err := child.consume(token.LEFT_BRACE, "Expected '{' before function body."); err != nil {
		return nil, err
	}
	if err := child.block(); err != nil {
		return nil, err
	}

	child.emitOps(chunk.OP_NIL, chunk.OP_RETURN)
	return &value.Function{
		Name:  name,
		Arity: arity,
		Chunk: child.chunk,
	}, nil
}

