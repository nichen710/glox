package chunk

import (
	"fmt"
	"strings"

	"glox/pkg/value"
)

const (
	MaxConstants = 1 << 16
	MaxJump      = 1<<16 - 1
)

type Chunk struct {
	Code       []byte
	Lines      []int
	Constants  []value.Value
	constIndex map[value.Value]int
}

func NewChunk() *Chunk {
	return &Chunk{
		Code:       make([]byte, 0),
		Lines:      make([]int, 0),
		Constants:  make([]value.Value, 0),
		constIndex: make(map[value.Value]int),
	}
}

func (c *Chunk) Write(b byte, line int) {
	c.Code = append(c.Code, b)
	c.Lines = append(c.Lines, line)
}

// Agrega el valor al pool de constantes (reutilizando si ya existe) y devuelve su indice
func (c *Chunk) AddConstant(val value.Value) int {
	if c.constIndex == nil {
		c.constIndex = make(map[value.Value]int)
	}
	if idx, ok := c.constIndex[val]; ok {
		return idx
	}
	idx := len(c.Constants)
	c.Constants = append(c.Constants, val)
	c.constIndex[val] = idx
	return idx
}

// Los operandos de los saltos son de 16 bits, big endian
func (c *Chunk) ReadShort(offset int) int {
	return int(c.Code[offset])<<8 | int(c.Code[offset+1])
}

func (c *Chunk) WriteShort(value int, line int) {
	c.Code = append(c.Code, byte(value>>8), byte(value))
	c.Lines = append(c.Lines, line, line)
}

func (c *Chunk) PatchShort(offset, value int) {
	c.Code[offset] = byte(value >> 8)
	c.Code[offset+1] = byte(value)
}

func (c *Chunk) LineAt(offset int) int {
	if offset >= 0 && offset < len(c.Lines) {
		return c.Lines[offset]
	}
	return 0
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
	case OP_CONSTANT, OP_DEFINE_GLOBAL, OP_GET_GLOBAL, OP_SET_GLOBAL,
		OP_CLASS, OP_GET_PROPERTY, OP_SET_PROPERTY, OP_METHOD, OP_GET_SUPER:
		index := c.Code[offset+1]
		fmt.Fprintf(sb, "<%s>\n", formatConstant(c.Constants[index]))
		return offset + 2
	case OP_INVOKE, OP_SUPER_INVOKE:
		index := c.Code[offset+1]
		argCount := c.Code[offset+2]
		fmt.Fprintf(sb, "<%s, %d>\n", formatConstant(c.Constants[index]), argCount)
		return offset + 3
	case OP_CONSTANT_LONG:
		index := c.ReadShort(offset + 1)
		fmt.Fprintf(sb, "<%s>\n", formatConstant(c.Constants[index]))
		return offset + 3
	case OP_ARRAY:
		count := c.ReadShort(offset + 1)
		fmt.Fprintf(sb, "<%d>\n", count)
		return offset + 3
	case OP_GET_LOCAL, OP_SET_LOCAL:
		slot := c.Code[offset+1]
		fmt.Fprintf(sb, "<@%d>\n", slot)
		return offset + 2
	case OP_CALL, OP_TAIL_CALL:
		argCount := c.Code[offset+1]
		fmt.Fprintf(sb, "<%d>\n", argCount)
		return offset + 2
	case OP_GET_UPVALUE, OP_SET_UPVALUE:
		slot := c.Code[offset+1]
		fmt.Fprintf(sb, "<^%d>\n", slot)
		return offset + 2
	case OP_CLOSURE:
		index := c.Code[offset+1]
		fnVal := c.Constants[index]
		fmt.Fprintf(sb, "<%s>\n", formatConstant(fnVal))
		offset += 2
		if fnVal.Fn != nil {
			for i := 0; i < fnVal.Fn.UpvalueCount; i++ {
				kind := "upvalue"
				if c.Code[offset] == 1 {
					kind = "local"
				}
				uvIndex := c.Code[offset+1]
				fmt.Fprintf(sb, "%04d|      | %s %d\n", offset, kind, uvIndex)
				offset += 2
			}
		}
		return offset
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

func formatConstant(val value.Value) string {
	if s, ok := val.AsString(); ok {
		return fmt.Sprintf("'%s'", s)
	}
	return val.String()
}
