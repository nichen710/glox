package chunk

type OpCode byte

const (
	// Instrucciones sin operandos
	OP_NIL OpCode = iota
	OP_TRUE
	OP_FALSE
	OP_RETURN
	OP_PRINT
	OP_POP
	OP_NEGATE
	OP_NOT
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_MODULO
	OP_EQUAL
	OP_GREATER
	OP_LESS

	// Instrucciones con operandos
	OP_CONSTANT
	OP_DEFINE_GLOBAL
	OP_GET_GLOBAL
	OP_SET_GLOBAL
	OP_GET_LOCAL
	OP_SET_LOCAL

	OP_JUMP
	OP_JUMP_IF_FALSE
	OP_LOOP
)

var opCodeNames = map[OpCode]string{
	OP_NIL:           "OP_NIL",
	OP_TRUE:          "OP_TRUE",
	OP_FALSE:         "OP_FALSE",
	OP_RETURN:        "OP_RETURN",
	OP_PRINT:         "OP_PRINT",
	OP_POP:           "OP_POP",
	OP_NEGATE:        "OP_NEGATE",
	OP_NOT:           "OP_NOT",
	OP_ADD:           "OP_ADD",
	OP_SUBTRACT:      "OP_SUBTRACT",
	OP_MULTIPLY:      "OP_MULTIPLY",
	OP_DIVIDE:        "OP_DIVIDE",
	OP_MODULO:        "OP_MODULO",
	OP_EQUAL:         "OP_EQUAL",
	OP_GREATER:       "OP_GREATER",
	OP_LESS:          "OP_LESS",
	OP_CONSTANT:      "OP_CONSTANT",
	OP_DEFINE_GLOBAL: "OP_DEFINE_GLOBAL",
	OP_GET_GLOBAL:    "OP_GET_GLOBAL",
	OP_SET_GLOBAL:    "OP_SET_GLOBAL",
	OP_GET_LOCAL:     "OP_GET_LOCAL",
	OP_SET_LOCAL:     "OP_SET_LOCAL",
	OP_JUMP:          "OP_JUMP",
	OP_JUMP_IF_FALSE: "OP_JUMP_IF_FALSE",
	OP_LOOP:          "OP_LOOP",
}

func (op OpCode) String() string {
	if name, ok := opCodeNames[op]; ok {
		return name
	}
	return "OP_UNKNOWN"
}
