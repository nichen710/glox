package compiler

const MaxLocals = 256

type local struct {
	name        string
	depth       int
	initialized bool
}

type CompilerContext struct {
	scopeDepth int // 0 es el scope global
	locals     []local
}

func NewCompilerContext() *CompilerContext {
	return &CompilerContext{
		scopeDepth: 0,
		locals:     make([]local, 0),
	}
}

func (ctx *CompilerContext) isGlobalScope() bool {
	return ctx.scopeDepth == 0
}

func (ctx *CompilerContext) beginScope() {
	ctx.scopeDepth++
}

// Saca las variables del scope que termina y devuelve cuantas eran
func (ctx *CompilerContext) endScope() int {
	ctx.scopeDepth--

	removed := 0
	for len(ctx.locals) > 0 && ctx.locals[len(ctx.locals)-1].depth > ctx.scopeDepth {
		ctx.locals = ctx.locals[:len(ctx.locals)-1]
		removed++
	}
	return removed
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
