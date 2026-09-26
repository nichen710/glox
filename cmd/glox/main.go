package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"glox/pkg/scanner"
	"glox/pkg/parser"
)

const (
	ModeRun      Mode = "run"
	ModeScanning Mode = "scanning"
	ModeParsing  Mode = "parsing"
	ModeResolve  Mode = "resolve"
)

type (
	Glox struct {
		mode   Mode
		debug  bool
		inRepl bool
	}

	Mode string
)

func NewGlox(mode Mode) *Glox {
	return &Glox{
		mode: mode,
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

	prsr := parser.NewParser(tokens)
	expr, err := prsr.Parse()
	if err != nil {
		fmt.Printf("Parsing Error: %v\n", err)
		return
	}

	if g.mode == ModeParsing {
		fmt.Println(expr)
		return
	}

	// TODO: Implement resolver
	if g.mode == ModeResolve {
		fmt.Println("Resolver not implemented yet.")
		return
	}

	// TODO: Implement interpreter
	fmt.Println("Interpreter not implemented yet.")
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
	}

	app := NewGlox(mode)

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
