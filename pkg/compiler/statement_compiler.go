package compiler

type StatementCompiler interface {
	CanCompile(c *Compiler) bool
	Compile(c *Compiler) error
}

