package optimizer

import "glox/pkg/chunk"

type JumpThreadingPass struct{}

// Si un salto apunta a otro OP_JUMP incondicional, redirige directo al destino final.
func (p *JumpThreadingPass) Run(c *chunk.Chunk) {
	for offset := 0; offset < len(c.Code); offset += instructionSize(c, offset) {
		op := chunk.OpCode(c.Code[offset])
		if op != chunk.OP_JUMP && op != chunk.OP_JUMP_IF_FALSE {
			continue
		}

		firstTarget := offset + 3 + c.ReadShort(offset+1)
		finalTarget := p.followJumps(c, firstTarget)
		if finalTarget == firstTarget {
			continue
		}

		newJump := finalTarget - (offset + 3)
		if newJump >= 0 && newJump <= chunk.MaxJump {
			c.PatchShort(offset+1, newJump)
		}
	}
}

func (p *JumpThreadingPass) followJumps(c *chunk.Chunk, target int) int {
	for steps := 0; steps < len(c.Code) && target >= 0 && target+2 < len(c.Code); steps++ {
		if chunk.OpCode(c.Code[target]) != chunk.OP_JUMP {
			break
		}
		next := target + 3 + c.ReadShort(target+1)
		if next == target {
			break
		}
		target = next
	}
	return target
}
