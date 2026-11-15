package compiler

import (
	"glox/pkg/chunk"
	"glox/pkg/token"
)

// Expresiones: stack effect +1
func (c *Compiler) expression() error {
	return c.parsePrecedence(PREC_ASSIGNMENT)
}

func (c *Compiler) parsePrecedence(precedence Precedence) error {
	tok := c.advance()
	rule := getRule(tok.Type)
	if rule.prefix == nil {
		return newCompileErrorf(tok, "Expected expression, got `%s`.", tok.Lexeme)
	}
	// Solo una expresion de la precedencia mas baja puede ser target: `a + b = 1` no es valido
	canAssign := precedence <= PREC_ASSIGNMENT
	if err := rule.prefix(c, canAssign); err != nil {
		return err
	}

	for precedence <= getRule(c.peek().Type).infixPrecedence {
		c.advance()
		if err := getRule(c.previous().Type).infix(c, canAssign); err != nil {
			return err
		}
	}

	// Si nadie consumio el `=`, lo que esta a la izquierda no era un target valido
	if canAssign && c.match(token.EQUAL) {
		return NewCompileError(c.previous(), "Invalid assignment target.")
	}
	return nil
}

func (c *Compiler) value(_ bool) error {
	return c.emitConstant(c.previous())
}

func (c *Compiler) literal(_ bool) error {
	switch c.previous().Type {
	case token.NIL:
		c.emitOp(chunk.OP_NIL)
	case token.TRUE:
		c.emitOp(chunk.OP_TRUE)
	case token.FALSE:
		c.emitOp(chunk.OP_FALSE)
	}
	return nil
}

func (c *Compiler) grouping(_ bool) error {
	if err := c.expression(); err != nil {
		return err
	}
	return c.consume(token.RIGHT_PAREN, "Expected ')' after expression.")
}

func (c *Compiler) unary(_ bool) error {
	operator := c.previous()
	if err := c.parsePrecedence(getRule(operator.Type).prefixPrecedence); err != nil {
		return err
	}

	switch operator.Type {
	case token.MINUS:
		c.emitOp(chunk.OP_NEGATE)
	case token.BANG:
		c.emitOp(chunk.OP_NOT)
	}
	return nil
}

func (c *Compiler) binary(_ bool) error {
	operator := c.previous()
	if err := c.parsePrecedence(getRule(operator.Type).infixPrecedence + 1); err != nil {
		return err
	}

	switch operator.Type {
	case token.PLUS:
		c.emitOp(chunk.OP_ADD)
	case token.MINUS:
		c.emitOp(chunk.OP_SUBTRACT)
	case token.STAR:
		c.emitOp(chunk.OP_MULTIPLY)
	case token.SLASH:
		c.emitOp(chunk.OP_DIVIDE)
	case token.PERCENT:
		c.emitOp(chunk.OP_MODULO)
	case token.EQUAL_EQUAL:
		c.emitOp(chunk.OP_EQUAL)
	case token.BANG_EQUAL:
		c.emitOp(chunk.OP_NOT_EQUAL)
	case token.GREATER:
		c.emitOp(chunk.OP_GREATER)
	case token.GREATER_EQUAL:
		c.emitOp(chunk.OP_GREATER_EQUAL)
	case token.LESS:
		c.emitOp(chunk.OP_LESS)
	case token.LESS_EQUAL:
		c.emitOp(chunk.OP_LESS_EQUAL)
	}
	return nil
}

// El operando izquierdo ya esta en el stack: si es falso, es el resultado
// y se saltea el derecho
func (c *Compiler) logicAnd(_ bool) error {
	endJump := c.emitJump(chunk.OP_JUMP_IF_FALSE)
	c.emitOp(chunk.OP_POP)
	if err := c.parsePrecedence(PREC_AND); err != nil {
		return err
	}
	return c.patchJump(endJump)
}

// No hay OP_JUMP_IF_TRUE: si es falso salta a evaluar el derecho,
// si no un salto incondicional lo saltea
func (c *Compiler) logicOr(_ bool) error {
	elseJump := c.emitJump(chunk.OP_JUMP_IF_FALSE)
	endJump := c.emitJump(chunk.OP_JUMP)

	if err := c.patchJump(elseJump); err != nil {
		return err
	}
	c.emitOp(chunk.OP_POP)
	if err := c.parsePrecedence(PREC_OR); err != nil {
		return err
	}
	return c.patchJump(endJump)
}

// Si viene un `=` es una asignacion, si no es una lectura de la variable
func (c *Compiler) variable(canAssign bool) error {
	return c.namedVariable(c.previous(), canAssign)
}

func (c *Compiler) namedVariable(name token.Token, canAssign bool) error {
	getOp, setOp, operand, err := c.resolveVariable(name)
	if err != nil {
		return err
	}

	if canAssign && c.match(token.EQUAL) {
		if err := c.expression(); err != nil {
			return err
		}
		c.emitOpByte(setOp, operand)
	} else {
		c.emitOpByte(getOp, operand)
	}
	return nil
}

func (c *Compiler) resolveVariable(name token.Token) (chunk.OpCode, chunk.OpCode, byte, error) {
	if slot, local := c.context.resolve(name.Lexeme); local != nil {
		if !local.initialized {
			return 0, 0, 0, NewCompileError(name, "Can't read local variable in its own initializer.")
		}
		return chunk.OP_GET_LOCAL, chunk.OP_SET_LOCAL, byte(slot), nil
	}

	if upIdx, found, err := c.resolveUpvalue(name); err != nil {
		return 0, 0, 0, err
	} else if found {
		return chunk.OP_GET_UPVALUE, chunk.OP_SET_UPVALUE, byte(upIdx), nil
	}

	index, err := c.identifierConstant(name)
	if err != nil {
		return 0, 0, 0, err
	}
	return chunk.OP_GET_GLOBAL, chunk.OP_SET_GLOBAL, index, nil
}

func (c *Compiler) call(_ bool) error {
	argCount, err := c.argumentList()
	if err != nil {
		return err
	}
	c.emitOpByte(chunk.OP_CALL, argCount)
	return nil
}

func (c *Compiler) argumentList() (byte, error) {
	argCount := 0
	if !c.check(token.RIGHT_PAREN) {
		for {
			if err := c.expression(); err != nil {
				return 0, err
			}
			if argCount >= MaxArgs {
				return 0, NewCompileError(c.previous(), "Can't have more than 255 arguments.")
			}
			argCount++
			if !c.match(token.COMMA) {
				break
			}
		}
	}
	if err := c.consume(token.RIGHT_PAREN, "Expected ')' after arguments."); err != nil {
		return 0, err
	}
	return byte(argCount), nil
}

func (c *Compiler) dot(canAssign bool) error {
	if err := c.consume(token.IDENTIFIER, "Expected property name after '.'."); err != nil {
		return err
	}
	nameIdx, err := c.identifierConstant(c.previous())
	if err != nil {
		return err
	}

	if canAssign && c.match(token.EQUAL) {
		if err := c.expression(); err != nil {
			return err
		}
		c.emitOpByte(chunk.OP_SET_PROPERTY, nameIdx)
	} else if c.match(token.LEFT_PAREN) {
		argCount, err := c.argumentList()
		if err != nil {
			return err
		}
		c.emitOp(chunk.OP_INVOKE)
		c.chunk.Write(nameIdx, c.currentLine())
		c.chunk.Write(argCount, c.currentLine())
	} else {
		c.emitOpByte(chunk.OP_GET_PROPERTY, nameIdx)
	}
	return nil
}

func (c *Compiler) this_(_ bool) error {
	if c.currentClass == nil {
		return NewCompileError(c.previous(), "Can't use 'this' outside of a class.")
	}
	return c.variable(false)
}

func (c *Compiler) super_(_ bool) error {
	kw := c.previous()
	if c.currentClass == nil {
		return NewCompileError(kw, "Can't use 'super' outside of a class.")
	}
	if !c.currentClass.hasSuperclass {
		return NewCompileError(kw, "Can't use 'super' in a class with no superclass.")
	}

	if err := c.consume(token.DOT, "Expected '.' after 'super'."); err != nil {
		return err
	}
	if err := c.consume(token.IDENTIFIER, "Expected superclass method name."); err != nil {
		return err
	}
	nameIdx, err := c.identifierConstant(c.previous())
	if err != nil {
		return err
	}

	if err := c.namedVariable(syntheticToken("this", kw.Line), false); err != nil {
		return err
	}
	if c.match(token.LEFT_PAREN) {
		argCount, err := c.argumentList()
		if err != nil {
			return err
		}
		if err := c.namedVariable(syntheticToken("super", kw.Line), false); err != nil {
			return err
		}
		c.emitOp(chunk.OP_SUPER_INVOKE)
		c.chunk.Write(nameIdx, c.currentLine())
		c.chunk.Write(argCount, c.currentLine())
	} else {
		if err := c.namedVariable(syntheticToken("super", kw.Line), false); err != nil {
			return err
		}
		c.emitOpByte(chunk.OP_GET_SUPER, nameIdx)
	}
	return nil
}

func syntheticToken(lexeme string, line int) token.Token {
	return token.NewToken(token.IDENTIFIER, lexeme, nil, line)
}
