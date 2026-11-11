package chunk

import (
	"fmt"
	"strings"
)

const MaxConstants = 256

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
	case OP_CONSTANT:
		index := c.Code[offset+1]
		fmt.Fprintf(sb, "<%s>\n", formatConstant(c.Constants[index]))
		return offset + 2
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
