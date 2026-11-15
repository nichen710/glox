package compiler

import (
	"glox/pkg/chunk"
	"glox/pkg/token"
	"glox/pkg/value"
)

func (c *Compiler) currentLine() int {
	if *c.current > 0 {
		return c.previous().Line
	}
	return c.peek().Line
}

func (c *Compiler) emitOp(op chunk.OpCode) {
	c.lastOpOffset = len(c.chunk.Code)
	c.chunk.Write(byte(op), c.currentLine())
}

func (c *Compiler) emitOpByte(op chunk.OpCode, operand byte) {
	c.lastOpOffset = len(c.chunk.Code)
	line := c.currentLine()
	c.chunk.Write(byte(op), line)
	c.chunk.Write(operand, line)
}

func (c *Compiler) emitOpShort(op chunk.OpCode, operand int) {
	c.lastOpOffset = len(c.chunk.Code)
	line := c.currentLine()
	c.chunk.Write(byte(op), line)
	c.chunk.WriteShort(operand, line)
}

func (c *Compiler) emitOps(ops ...chunk.OpCode) {
	for _, op := range ops {
		c.emitOp(op)
	}
}

// Todavia no se sabe a donde salta: deja 0xFF 0xFF y devuelve donde
// esta el operando para emparcharlo despues
func (c *Compiler) emitJump(op chunk.OpCode) int {
	c.lastOpOffset = len(c.chunk.Code)
	line := c.currentLine()
	c.chunk.Write(byte(op), line)
	c.chunk.WriteShort(0xFFFF, line)
	return len(c.chunk.Code) - 2
}

// El salto es relativo al final del operando
func (c *Compiler) patchJump(operand int) error {
	jump := len(c.chunk.Code) - operand - 2
	if jump > chunk.MaxJump {
		return NewCompileError(c.previous(), "Too much code to jump over.")
	}
	c.chunk.PatchShort(operand, jump)
	return nil
}

// Salta para atras, contando los dos bytes del propio operando
func (c *Compiler) emitLoop(loopStart int) error {
	c.lastOpOffset = len(c.chunk.Code)
	line := c.currentLine()
	c.chunk.Write(byte(chunk.OP_LOOP), line)
	jump := len(c.chunk.Code) - loopStart + 2
	if jump > chunk.MaxJump {
		return NewCompileError(c.previous(), "Loop body too large.")
	}
	c.chunk.WriteShort(jump, line)
	return nil
}

func (c *Compiler) emitConstant(tok token.Token) error {
	index, err := c.makeConstant(tok, value.FromAny(tok.Literal))
	if err != nil {
		return err
	}
	if index < 256 {
		c.emitOpByte(chunk.OP_CONSTANT, byte(index))
	} else {
		c.emitOpShort(chunk.OP_CONSTANT_LONG, index)
	}
	return nil
}

// Los nombres de las variables globales se guardan como strings en el pool de constantes
func (c *Compiler) identifierConstant(name token.Token) (byte, error) {
	return c.makeByteConstant(name, value.String(name.Lexeme))
}

func (c *Compiler) makeByteConstant(tok token.Token, val value.Value) (byte, error) {
	index := c.chunk.AddConstant(val)
	if index >= 256 {
		return 0, NewCompileError(tok, "Too many constants in one chunk.")
	}
	return byte(index), nil
}

func (c *Compiler) makeConstant(tok token.Token, val value.Value) (int, error) {
	index := c.chunk.AddConstant(val)
	if index >= chunk.MaxConstants {
		return 0, NewCompileError(tok, "Too many constants in one chunk.")
	}
	return index, nil
}
