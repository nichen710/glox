package chunk

import (
	"fmt"
	"strings"
)

const MaxConstants = 256

const MaxJump = 1<<16 - 1

type Chunk struct {
	Code      []byte
	Constants []any
}

func NewChunk() *Chunk {
	return &Chunk{
		Code:      make([]byte, 0),
		Constants: make([]any, 0),
	}
}

func (c *Chunk) Write(b byte) {
	c.Code = append(c.Code, b)
}

// Agrega el valor al pool de constantes y devuelve su indice
func (c *Chunk) AddConstant(value any) int {
	c.Constants = append(c.Constants, value)
	return len(c.Constants) - 1
}

// Los operandos de los saltos son de 16 bits, big endian
func (c *Chunk) ReadShort(offset int) int {
	return int(c.Code[offset])<<8 | int(c.Code[offset+1])
}

// Desensamblado del chunk, una instruccion por linea
func (c *Chunk) String() string {
	var sb strings.Builder
	for offset := 0; offset < len(c.Code); {
		offset = c.disassembleInstruction(&sb, offset)
	}
	return strings.TrimSuffix(sb.String(), "\n")
}

func (c *Chunk) disassembleInstruction(sb *strings.Builder, offset int) int {
	op := OpCode(c.Code[offset])
	fmt.Fprintf(sb, "%04d| %s", offset, op)

	switch op {
	case OP_CONSTANT, OP_DEFINE_GLOBAL, OP_GET_GLOBAL, OP_SET_GLOBAL:
		index := c.Code[offset+1]
		fmt.Fprintf(sb, "<%s>\n", formatConstant(c.Constants[index]))
		return offset + 2
	case OP_GET_LOCAL, OP_SET_LOCAL:
		slot := c.Code[offset+1]
		fmt.Fprintf(sb, "<@%d>\n", slot)
		return offset + 2
	case OP_JUMP, OP_JUMP_IF_FALSE:
		fmt.Fprintf(sb, "<+%d>\n", c.ReadShort(offset+1)+3)
		return offset + 3
	case OP_LOOP:
		fmt.Fprintf(sb, "<-%d>\n", c.ReadShort(offset+1)-3)
		return offset + 3
	default:
		sb.WriteString("\n")
		return offset + 1
	}
}

func formatConstant(value any) string {
	if s, ok := value.(string); ok {
		return fmt.Sprintf("'%s'", s)
	}
	return fmt.Sprintf("%v", value)
}
