package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"glox/pkg/compiler"
	"glox/pkg/interpreter"
	"glox/pkg/parser"
	"glox/pkg/resolver"
	"glox/pkg/scanner"
	"glox/pkg/token"
	"glox/pkg/vm"
)

const (
	ModeRun      Mode = "run"
	ModeScanning Mode = "scanning"
	ModeParsing  Mode = "parsing"
	ModeResolve  Mode = "resolve"
	ModeCompile  Mode = "compile"
)

type (
	Glox struct {
		mode         Mode
		debug        bool
		inRepl       bool
		useVM        bool
		interpreter  *interpreter.Interpreter
		bindingTable *resolver.BindingTable
		vm           *vm.VM
	}

	Mode string
)

func NewGlox(mode Mode, useVM bool) *Glox {
	table := resolver.NewBindingTable()
	interp := interpreter.NewInterpreter()
	interp.SetBindings(table)

	return &Glox{
		mode:         mode,
		useVM:        useVM,
		interpreter:  interp,
		bindingTable: table,
		vm:           vm.NewVM(),
	}
}

func (g *Glox) Run(source string) {
	scnr := scanner.NewScanner(source)
	tokens, err := scnr.Scan()
	if err != nil {
		fmt.Printf("Scanning Error: %v\n", err)
		return
	}

	if g.mode == ModeScanning {
		for _, tok := range tokens {
			fmt.Println(tok)
		}
		return
	}

	if g.mode == ModeCompile || g.useVM {
		g.runBytecode(tokens)
		return
	}

	prsr := parser.NewParser(tokens)
	stmts, err := prsr.Parse()
	if err != nil {
		fmt.Printf("Parsing Error: %v\n", err)
		return
	}

	if g.mode == ModeParsing {
		for _, stmt := range stmts {
			fmt.Println(stmt)
		}
		return
	}

	rslvr := resolver.NewDefaultResolver(g.bindingTable)
	if err := rslvr.Resolve(stmts); err != nil {
		fmt.Printf("Resolve Error: %v\n", err)
		return
	}

	if g.mode == ModeResolve {
		for _, stmt := range stmts {
			fmt.Println(stmt)
		}
		fmt.Printf("Bindings: %s\n", g.bindingTable)
		return
	}

	// Execute statements
	result, err := g.interpreter.Interpret(stmts)
	if err != nil {
		fmt.Printf("Runtime Error: %v\n", err)
		return
	}

	if g.inRepl && result != nil {
		fmt.Println(interpreter.Stringify(result))
	}
}

func (g *Glox) runBytecode(tokens []token.Token) {
	cmplr := compiler.NewCompiler(tokens)
	chnk, err := cmplr.Compile()
	if err != nil {
		fmt.Printf("Compile Error: %v\n", err)
		return
	}

	if g.mode == ModeCompile {
		fmt.Println(chnk)
		return
	}

	if err := g.vm.Run(chnk); err != nil {
		fmt.Printf("Runtime Error: %v\n", err)
	}
}

func (g *Glox) RunFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file %s: %v\n", path, err)
	}
	g.Run(string(bytes))
}

func (g *Glox) RunPrompt() {
	g.inRepl = true
	fmt.Println("glox REPL - press Ctrl+C or Ctrl+D to quit")

	for {
		fmt.Print(">> ")
		scnr := bufio.NewScanner(os.Stdin)
		if !scnr.Scan() {
			break
		}

		if scnr.Text() == "" {
			continue
		}

		g.Run(scnr.Text())
	}
}

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: glox [flags] [file]")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
	scanningFlag := flag.Bool("scanning", false, "Run in scanning mode (prints tokens)")
	parsingFlag := flag.Bool("parsing", false, "Run in parsing mode (prints AST)")
	resolveFlag := flag.Bool("resolve", false, "Run in resolve mode (prints resolved scopes)")
	compileFlag := flag.Bool("compile", false, "Run in compile mode (prints bytecode)")
	vmFlag := flag.Bool("vm", false, "Execute compiling to bytecode instead of tree-walk")
	helpFlag := flag.Bool("help", false, "Print help message")

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	mode := ModeRun
	if *scanningFlag {
		mode = ModeScanning
	} else if *parsingFlag {
		mode = ModeParsing
	} else if *resolveFlag {
		mode = ModeResolve
	} else if *compileFlag {
		mode = ModeCompile
	}

	app := NewGlox(mode, *vmFlag)

	args := flag.Args()
	if len(args) > 1 {
		fmt.Println("Error: too many arguments")
		flag.Usage()
		os.Exit(64)
	} else if len(args) == 1 {
		app.RunFile(args[0])
	} else {
		app.RunPrompt()
	}
}
