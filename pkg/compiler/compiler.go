package compiler

import (
	"glox/pkg/chunk"
	"glox/pkg/token"
)


const (
	FnScript FunctionType = iota
	FnFunction
)

type (
	FunctionType uint8

	Compiler struct {
		tokens     []token.Token
		current    *int
		chunk      *chunk.Chunk
		context    *CompilerContext
		fnType     FunctionType
		statements []StatementCompiler
	}
)

func NewCompiler(tokens []token.Token) *Compiler {
	pos := 0
	return &Compiler{
		tokens:  tokens,
		current: &pos,
		chunk:   chunk.NewChunk(),
		context: NewCompilerContext(),
		fnType:  FnScript,
		statements: []StatementCompiler{
			&PrintCompiler{},
			&IfCompiler{},
			&WhileCompiler{},
			&ForCompiler{},
			&BlockCompiler{},
			&ReturnCompiler{},
			&ExpressionStatementCompiler{},
		},
	}
}

func (c *Compiler) newChild(fnType FunctionType) *Compiler {
	return &Compiler{
		tokens:     c.tokens,
		current:    c.current,
		chunk:      chunk.NewChunk(),
		context:    NewCompilerContext(),
		fnType:     fnType,
		statements: c.statements,
	}
}

func (c *Compiler) Compile() (*chunk.Chunk, error) {
	for !c.isAtEnd() {
		if err := c.declaration(); err != nil {
			return nil, err
		}
	}
	c.emitOp(chunk.OP_RETURN)
	return c.chunk, nil
}

// Statements: stack effect 0
func (c *Compiler) declaration() error {
	if c.check(token.FUN) {
		return (&FunctionCompiler{}).Compile(c)
	}
	if c.check(token.VAR) {
		return (&VarCompiler{}).Compile(c)
	}
	return c.statement()
}

func (c *Compiler) statement() error {
	for _, stmt := range c.statements {
		if stmt.CanCompile(c) {
			return stmt.Compile(c)
		}
	}
	return nil
}

func (c *Compiler) block() error {
	for !c.check(token.RIGHT_BRACE) && !c.isAtEnd() {
		if err := c.declaration(); err != nil {
			return err
		}
	}
	return c.consume(token.RIGHT_BRACE, "Expected '}' after block.")
}

func (c *Compiler) condition(openMessage, closeMessage string) error {
	if err := c.consume(token.LEFT_PAREN, openMessage); err != nil {
		return err
	}
	if err := c.expression(); err != nil {
		return err
	}
	return c.consume(token.RIGHT_PAREN, closeMessage)
}

func (c *Compiler) declareLocal(name token.Token) error {
	if c.context.isDeclaredInCurrentScope(name.Lexeme) {
		return NewCompileError(name, "Already a variable with this name in this scope.")
	}
	if c.context.isFull() {
		return NewCompileError(name, "Too many local variables in scope.")
	}
	c.context.declare(name.Lexeme)
	return nil
}

func (c *Compiler) endScope() {
	for range c.context.endScope() {
		c.emitOp(chunk.OP_POP)
	}
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
		*c.current++
	}
	return tok
}

func (c *Compiler) isAtEnd() bool {
	return c.peek().Type == token.EOF
}

func (c *Compiler) peek() token.Token {
	return c.tokens[*c.current]
}

func (c *Compiler) previous() token.Token {
	return c.tokens[*c.current-1]
}
