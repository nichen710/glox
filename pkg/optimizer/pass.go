package optimizer

import "glox/pkg/chunk"

type Pass interface {
	Run(c *chunk.Chunk)
}

func instructionSize(c *chunk.Chunk, offset int) int {
	switch chunk.OpCode(c.Code[offset]) {
	case chunk.OP_CONSTANT,
		chunk.OP_DEFINE_GLOBAL,
		chunk.OP_GET_GLOBAL,
		chunk.OP_SET_GLOBAL,
		chunk.OP_GET_LOCAL,
		chunk.OP_SET_LOCAL,
		chunk.OP_GET_UPVALUE,
		chunk.OP_SET_UPVALUE,
		chunk.OP_CALL,
		chunk.OP_TAIL_CALL,
		chunk.OP_CLASS,
		chunk.OP_GET_PROPERTY,
		chunk.OP_SET_PROPERTY,
		chunk.OP_METHOD,
		chunk.OP_GET_SUPER:
		return 2
	case chunk.OP_JUMP, chunk.OP_JUMP_IF_FALSE, chunk.OP_LOOP, chunk.OP_CONSTANT_LONG,
		chunk.OP_INVOKE, chunk.OP_SUPER_INVOKE:
		return 3
	case chunk.OP_CLOSURE:
		fnVal := c.Constants[c.Code[offset+1]]
		if fnVal.Fn != nil {
			return 2 + fnVal.Fn.UpvalueCount*2
		}
		return 2
	default:
		return 1
	}
}
