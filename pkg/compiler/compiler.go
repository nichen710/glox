package compiler

import (
	"glox/pkg/chunk"
	"glox/pkg/token"
)

type Compiler struct {
	tokens     []token.Token
	current    int
	chunk      *chunk.Chunk
	context    *CompilerContext
	statements []StatementCompiler
}

func NewCompiler(tokens []token.Token) *Compiler {
	return &Compiler{
		tokens:  tokens,
		current: 0,
		chunk:   chunk.NewChunk(),
		context: NewCompilerContext(),
		statements: []StatementCompiler{
			&PrintCompiler{},
			&IfCompiler{},
			&WhileCompiler{},
			&ForCompiler{},
			&BlockCompiler{},
			&ExpressionStatementCompiler{},
		},
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

func (c *Compiler) condition(openMessage, closeMessage string) error {
	if err := c.consume(token.LEFT_PAREN, openMessage); err != nil {
		return err
	}
	if err := c.expression(); err != nil {
		return err
	}
	return c.consume(token.RIGHT_PAREN, closeMessage)
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
