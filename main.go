//go:generate goyacc -l -o parser.go parser.y
package main

import (
	"bufio"
	"fmt"
	"os"
)

func typecheckProgram(prg []definition) {
	mgr := newTypMgr()
	e := &typEnv{
		make(map[string]typ, 0),
		nil,
	}

	intTyp := &typBase{"Int"}
	binOpTyp := &typArr{
		intTyp,
		&typArr{
			intTyp,
			intTyp,
		},
	}

	e.bind("+", binOpTyp)
	e.bind("-", binOpTyp)
	e.bind("*", binOpTyp)
	e.bind("/", binOpTyp)

	for _, d := range prg {
		fmt.Printf("First typechecking: %v\n", d)
		d.typecheckFirst(mgr, e)
	}

	for _, d := range prg {
		fmt.Printf("Second typechecking: %v\n", d)
		d.typecheckSecond(mgr, e)
	}

	fmt.Println("Types:")
	for n, t := range e.names {
		fmt.Printf("\t%s: %v\n", n, t.typString(mgr))
	}
}

func main() {
	l := newLexer(bufio.NewReader(os.Stdin))
	yyParse(l)
	fmt.Printf("Parsed tree:\n%v\n", l.result)
	typecheckProgram(l.result)
}
