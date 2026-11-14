package chunk

import "testing"

func TestChunk_WriteAndConstants(t *testing.T) {
	c := NewChunk()
	index := c.AddConstant(1.5)
	c.Write(byte(OP_CONSTANT))
	c.Write(byte(index))
	c.Write(byte(OP_RETURN))

	expectedCode := []byte{byte(OP_CONSTANT), 0, byte(OP_RETURN)}
	if len(c.Code) != len(expectedCode) {
		t.Fatalf("got %d bytes, want %d", len(c.Code), len(expectedCode))
	}
	for i, b := range expectedCode {
		if c.Code[i] != b {
			t.Errorf("byte %d: got %d, want %d", i, c.Code[i], b)
		}
	}

	if c.Constants[index] != 1.5 {
		t.Errorf("constant %d: got %v, want 1.5", index, c.Constants[index])
	}
}

func TestChunk_Disassemble(t *testing.T) {
	c := NewChunk()
	c.Write(byte(OP_CONSTANT))
	c.Write(byte(c.AddConstant(6.0)))
	c.Write(byte(OP_CONSTANT))
	c.Write(byte(c.AddConstant("abc")))
	c.Write(byte(OP_TRUE))
	c.Write(byte(OP_PRINT))
	c.Write(byte(OP_RETURN))

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
	c.Write(byte(OP_JUMP_IF_FALSE))
	c.Write(0)
	c.Write(1)
	c.Write(byte(OP_POP))
	c.Write(byte(OP_LOOP))
	c.Write(1)
	c.Write(0)

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
