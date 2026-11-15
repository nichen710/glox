package optimizer

import (
	"testing"

	"glox/pkg/chunk"
	"glox/pkg/compiler"
	"glox/pkg/scanner"
	"glox/pkg/vm"
)

func TestOptimizer_JumpThreading(t *testing.T) {
	c := chunk.NewChunk()
	// 0000: OP_JUMP_IF_FALSE -> 0006 (salta a otro OP_JUMP)
	c.Write(byte(chunk.OP_JUMP_IF_FALSE), 1)
	c.WriteShort(3, 1)
	// 0003: OP_JUMP -> 0009
	c.Write(byte(chunk.OP_JUMP), 1)
	c.WriteShort(3, 1)
	// 0006: OP_JUMP -> 0009 (destino intermedio)
	c.Write(byte(chunk.OP_JUMP), 1)
	c.WriteShort(0, 1)
	// 0009: OP_RETURN
	c.Write(byte(chunk.OP_RETURN), 1)

	Optimize(c)

	// 0000 debe apuntar directo a 0009 (offset +6 desde 0003)
	if got := c.ReadShort(1); got != 6 {
		t.Errorf("OP_JUMP_IF_FALSE offset: got %d, want 6", got)
	}
}

func TestOptimizer_PreservesSemantics(t *testing.T) {
	source := `
	var x = 0;
	if (false) {
		x = 1;
	} else if (false) {
		x = 2;
	} else {
		x = 3;
	}
	print x;
	print (false and true) or (true and 42);
	`
	tokens, err := scanner.NewScanner(source).Scan()
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}
	chnk, err := compiler.NewCompiler(tokens).Compile()
	if err != nil {
		t.Fatalf("compile error: %v", err)
	}
	Optimize(chnk)

	if err := vm.NewVM().Run(chnk); err != nil {
		t.Fatalf("runtime error: %v", err)
	}
}
