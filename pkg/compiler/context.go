package compiler

const MaxLocals = 256

type (
	local struct {
		name        string
		depth       int
		initialized bool
		captured    bool
	}

	upvalue struct {
		index   byte
		isLocal bool
	}

	CompilerContext struct {
		scopeDepth int // 0 es el scope global
		locals     []local
		upvalues   []upvalue
	}
)

func NewCompilerContext() *CompilerContext {
	return &CompilerContext{
		scopeDepth: 0,
		locals:     make([]local, 0),
		upvalues:   make([]upvalue, 0),
	}
}

func (ctx *CompilerContext) isGlobalScope() bool {
	return ctx.scopeDepth == 0
}

func (ctx *CompilerContext) beginScope() {
	ctx.scopeDepth++
}

// Saca las variables del scope que termina y devuelve si cada una fue capturada
func (ctx *CompilerContext) endScope() []bool {
	ctx.scopeDepth--

	captured := make([]bool, 0)
	for len(ctx.locals) > 0 && ctx.locals[len(ctx.locals)-1].depth > ctx.scopeDepth {
		last := ctx.locals[len(ctx.locals)-1]
		ctx.locals = ctx.locals[:len(ctx.locals)-1]
		captured = append(captured, last.captured)
	}
	return captured
}

func (ctx *CompilerContext) declare(name string) {
	ctx.locals = append(ctx.locals, local{name: name, depth: ctx.scopeDepth})
}

func (ctx *CompilerContext) markInitialized() {
	ctx.locals[len(ctx.locals)-1].initialized = true
}

func (ctx *CompilerContext) isFull() bool {
	return len(ctx.locals) >= MaxLocals
}

// Solo mira las variables del scope actual
func (ctx *CompilerContext) isDeclaredInCurrentScope(name string) bool {
	for i := len(ctx.locals) - 1; i >= 0; i-- {
		if ctx.locals[i].depth < ctx.scopeDepth {
			return false
		}
		if ctx.locals[i].name == name {
			return true
		}
	}
	return false
}

func (ctx *CompilerContext) resolve(name string) (int, *local) {
	for i := len(ctx.locals) - 1; i >= 0; i-- {
		if ctx.locals[i].name == name {
			return i, &ctx.locals[i]
		}
	}
	return -1, nil
}

func (ctx *CompilerContext) addUpvalue(index byte, isLocal bool) (int, bool) {
	for i, uv := range ctx.upvalues {
		if uv.index == index && uv.isLocal == isLocal {
			return i, true
		}
	}
	if len(ctx.upvalues) >= MaxLocals {
		return 0, false
	}
	ctx.upvalues = append(ctx.upvalues, upvalue{index: index, isLocal: isLocal})
	return len(ctx.upvalues) - 1, true
}
