package optimizer

import (
	"glox/pkg/chunk"
	"glox/pkg/value"
)

type Optimizer struct {
	passes []Pass
}

func NewOptimizer(passes ...Pass) *Optimizer {
	return &Optimizer{passes: passes}
}

func NewDefaultOptimizer() *Optimizer {
	return NewOptimizer(
		&JumpThreadingPass{},
	)
}

func Optimize(c *chunk.Chunk) {
	NewDefaultOptimizer().Optimize(c)
}

// Aplica todas las pasadas sobre el chunk y recursivamente sobre sus funciones.
func (o *Optimizer) Optimize(c *chunk.Chunk) {
	o.optimizeChunk(c, make(map[*chunk.Chunk]bool))
}

func (o *Optimizer) optimizeChunk(c *chunk.Chunk, visited map[*chunk.Chunk]bool) {
	if c == nil || visited[c] {
		return
	}
	visited[c] = true

	for _, val := range c.Constants {
		if val.Kind == value.KindFunction && val.Fn != nil {
			if fnChunk, ok := val.Fn.Chunk.(*chunk.Chunk); ok {
				o.optimizeChunk(fnChunk, visited)
			}
		}
	}

	for _, pass := range o.passes {
		pass.Run(c)
	}
}
