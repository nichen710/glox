package compiler

import (
	"glox/pkg/chunk"
	"glox/pkg/token"
)

type Compiler struct {
	tokens  []token.Token
	current int
	chunk   *chunk.Chunk
}

func NewCompiler(tokens []token.Token) *Compiler {
	return &Compiler{
		tokens:  tokens,
		current: 0,
		chunk:   chunk.NewChunk(),
	}
}

func (c *Compiler) Compile() (*chunk.Chunk, error) {
	for !c.isAtEnd() {
		if err := c.statement(); err != nil {
			return nil, err
		}
	}
	c.emitOp(chunk.OP_RETURN)
	return c.chunk, nil
}

// Statements: stack effect 0
func (c *Compiler) statement() error {
	if c.match(token.VAR) {
		return c.varDeclaration()
	}
	if c.match(token.PRINT) {
		return c.printStatement()
	}
	return c.expressionStatement()
}

func (c *Compiler) varDeclaration() error {
	if err := c.consume(token.IDENTIFIER, "Expected variable name."); err != nil {
		return err
	}
	name := c.previous()
	index, err := c.identifierConstant(name)
	if err != nil {
		return err
	}

	if c.match(token.EQUAL) {
		if err := c.expression(); err != nil {
			return err
		}
	} else {
		c.emitOp(chunk.OP_NIL)
	}
	if err := c.consume(token.SEMICOLON, "Expected ';' after variable declaration."); err != nil {
		return err
	}

	c.emitOp(chunk.OP_DEFINE_GLOBAL)
	c.chunk.Write(index)
	return nil
}

func (c *Compiler) printStatement() error {
	if err := c.expression(); err != nil {
		return err
	}
	if err := c.consume(token.SEMICOLON, "Expected ';' after value."); err != nil {
		return err
	}
	c.emitOp(chunk.OP_PRINT)
	return nil
}

func (c *Compiler) expressionStatement() error {
	if err := c.expression(); err != nil {
		return err
	}
	if err := c.consume(token.SEMICOLON, "Expected ';' after expression."); err != nil {
		return err
	}
	c.emitOp(chunk.OP_POP)
	return nil
}

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
		c.emitOps(chunk.OP_EQUAL, chunk.OP_NOT)
	case token.GREATER:
		c.emitOp(chunk.OP_GREATER)
	case token.GREATER_EQUAL:
		c.emitOps(chunk.OP_LESS, chunk.OP_NOT)
	case token.LESS:
		c.emitOp(chunk.OP_LESS)
	case token.LESS_EQUAL:
		c.emitOps(chunk.OP_GREATER, chunk.OP_NOT)
	}
	return nil
}

// Si viene un `=` es una asignacion, si no es una lectura de la variable
func (c *Compiler) variable(canAssign bool) error {
	index, err := c.identifierConstant(c.previous())
	if err != nil {
		return err
	}

	if canAssign && c.match(token.EQUAL) {
		if err := c.expression(); err != nil {
			return err
		}
		c.emitOp(chunk.OP_SET_GLOBAL)
	} else {
		c.emitOp(chunk.OP_GET_GLOBAL)
	}
	c.chunk.Write(index)
	return nil
}

func (c *Compiler) emitOp(op chunk.OpCode) {
	c.chunk.Write(byte(op))
}

func (c *Compiler) emitOps(ops ...chunk.OpCode) {
	for _, op := range ops {
		c.emitOp(op)
	}
}

func (c *Compiler) emitConstant(tok token.Token) error {
	index, err := c.makeConstant(tok, tok.Literal)
	if err != nil {
		return err
	}
	c.emitOp(chunk.OP_CONSTANT)
	c.chunk.Write(index)
	return nil
}

// Los nombres de las variables globales se guardan como strings en el pool de constantes
func (c *Compiler) identifierConstant(name token.Token) (byte, error) {
	return c.makeConstant(name, name.Lexeme)
}

func (c *Compiler) makeConstant(tok token.Token, value any) (byte, error) {
	index := c.chunk.AddConstant(value)
	if index >= chunk.MaxConstants {
		return 0, NewCompileError(tok, "Too many constants in one chunk.")
	}
	return byte(index), nil
}

func (c *Compiler) consume(tokenType token.TokenType, message string) error {
	if c.check(tokenType) {
		c.advance()
		return nil
	}
	return NewCompileError(c.peek(), message)
}

func (c *Compiler) match(tokenType token.TokenType) bool {
	if !c.check(tokenType) {
		return false
	}
	c.advance()
	return true
}

func (c *Compiler) check(tokenType token.TokenType) bool {
	return c.peek().Type == tokenType
}

func (c *Compiler) advance() token.Token {
	tok := c.peek()
	if !c.isAtEnd() {
		c.current++
	}
	return tok
}

func (c *Compiler) isAtEnd() bool {
	return c.peek().Type == token.EOF
}

func (c *Compiler) peek() token.Token {
	return c.tokens[c.current]
}

func (c *Compiler) previous() token.Token {
	return c.tokens[c.current-1]
}
