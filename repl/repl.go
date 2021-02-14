package repl

import (
	"bufio"
	"fmt"
	"interpreter-go/evaluator"
	"interpreter-go/lexer"
	"interpreter-go/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		parser := parser.New(l)
		program := parser.ParseProgram()

		if len(parser.Errors()) != 0 {
			printParseErrors(out, parser.Errors())
			continue
		}

		obj := evaluator.Eval(program)
		if obj != nil {
			io.WriteString(out, obj.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! we ran into some monkey business here!\n")
	io.WriteString(out, "parse error:\n")
	for _, e := range errors {
		io.WriteString(out, "\t"+e+"\n")
	}
}
