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
	if c.match(token.PRINT) {
		return c.printStatement()
	}
	return c.expressionStatement()
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
	if err := rule.prefix(c); err != nil {
		return err
	}

	for precedence <= getRule(c.peek().Type).infixPrecedence {
		c.advance()
		if err := getRule(c.previous().Type).infix(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) value() error {
	return c.emitConstant(c.previous())
}

func (c *Compiler) literal() error {
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

func (c *Compiler) grouping() error {
	if err := c.expression(); err != nil {
		return err
	}
	return c.consume(token.RIGHT_PAREN, "Expected ')' after expression.")
}

func (c *Compiler) unary() error {
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

func (c *Compiler) binary() error {
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

func (c *Compiler) emitOp(op chunk.OpCode) {
	c.chunk.Write(byte(op))
}

func (c *Compiler) emitOps(ops ...chunk.OpCode) {
	for _, op := range ops {
		c.emitOp(op)
	}
}

func (c *Compiler) emitConstant(tok token.Token) error {
	index := c.chunk.AddConstant(tok.Literal)
	if index >= chunk.MaxConstants {
		return NewCompileError(tok, "Too many constants in one chunk.")
	}
	c.emitOp(chunk.OP_CONSTANT)
	c.chunk.Write(byte(index))
	return nil
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
