package chunk

import (
	"testing"

	"glox/pkg/value"
)

func TestChunk_WriteAndConstants(t *testing.T) {
	c := NewChunk()
	index := c.AddConstant(value.Number(1.5))
	c.Write(byte(OP_CONSTANT), 1)
	c.Write(byte(index), 1)
	c.Write(byte(OP_RETURN), 1)

	expectedCode := []byte{byte(OP_CONSTANT), 0, byte(OP_RETURN)}
	if len(c.Code) != len(expectedCode) {
		t.Fatalf("got %d bytes, want %d", len(c.Code), len(expectedCode))
	}
	for i, b := range expectedCode {
		if c.Code[i] != b {
			t.Errorf("byte %d: got %d, want %d", i, c.Code[i], b)
		}
	}

	if !c.Constants[index].Equals(value.Number(1.5)) {
		t.Errorf("constant %d: got %v, want 1.5", index, c.Constants[index])
	}
}

func TestChunk_Disassemble(t *testing.T) {
	c := NewChunk()
	c.Write(byte(OP_CONSTANT), 1)
	c.Write(byte(c.AddConstant(value.Number(6.0))), 1)
	c.Write(byte(OP_CONSTANT), 2)
	c.Write(byte(c.AddConstant(value.String("abc"))), 2)
	c.Write(byte(OP_TRUE), 3)
	c.Write(byte(OP_PRINT), 3)
	c.Write(byte(OP_RETURN), 4)

	expected := "0000| OP_CONSTANT<6>\n" +
		"0002| OP_CONSTANT<'abc'>\n" +
		"0004| OP_TRUE\n" +
		"0005| OP_PRINT\n" +
		"0006| OP_RETURN"

	if got := c.String(); got != expected {
		t.Errorf("got:\n%s\nwant:\n%s", got, expected)
	}
}

func TestChunk_DisassembleJumps(t *testing.T) {
	c := NewChunk()
	c.Write(byte(OP_JUMP_IF_FALSE), 1)
	c.Write(0, 1)
	c.Write(1, 1)
	c.Write(byte(OP_POP), 2)
	c.Write(byte(OP_LOOP), 3)
	c.Write(1, 3)
	c.Write(0, 3)

	expected := "0000| OP_JUMP_IF_FALSE<+4>\n" +
		"0003| OP_POP\n" +
		"0004| OP_LOOP<-253>"

	if got := c.String(); got != expected {
		t.Errorf("got:\n%s\nwant:\n%s", got, expected)
	}
}

func TestOpCode_String(t *testing.T) {
	tests := []struct {
		op       OpCode
		expected string
	}{
		{OP_NIL, "OP_NIL"},
		{OP_TRUE, "OP_TRUE"},
		{OP_FALSE, "OP_FALSE"},
		{OP_RETURN, "OP_RETURN"},
		{OP_PRINT, "OP_PRINT"},
		{OP_POP, "OP_POP"},
		{OP_CONSTANT, "OP_CONSTANT"},
		{OpCode(255), "OP_UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.op.String(); got != tt.expected {
				t.Errorf("got %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestChunk_DeduplicatesAndLines(t *testing.T) {
	c := NewChunk()
	i1 := c.AddConstant(value.Number(42.0))
	i2 := c.AddConstant(value.Number(42.0))
	i3 := c.AddConstant(value.String("x"))
	i4 := c.AddConstant(value.String("x"))

	if i1 != i2 || i3 != i4 {
		t.Errorf("expected deduplicated indices, got (%d, %d) and (%d, %d)", i1, i2, i3, i4)
	}
	if len(c.Constants) != 2 {
		t.Errorf("got %d constants, want 2", len(c.Constants))
	}

	c.Write(byte(OP_CONSTANT), 3)
	c.Write(byte(i1), 3)
	c.WriteShort(10, 5)

	if got := c.LineAt(0); got != 3 {
		t.Errorf("offset 0 line: got %d, want 3", got)
	}
	if got := c.LineAt(2); got != 5 {
		t.Errorf("offset 2 line: got %d, want 5", got)
	}
	if got := c.LineAt(99); got != 0 {
		t.Errorf("out-of-bounds offset line: got %d, want 0", got)
	}
}
