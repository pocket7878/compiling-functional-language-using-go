//go:generate goyacc -l -o parser.go parser.y
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
)

func typecheckProgram(prg []definition) error {
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
		err := d.typecheckSecond(mgr, e)
		if err != nil {
			return err
		}
	}

	for _, d := range prg {
		err := d.resolve(mgr)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("resolving: %v", d))
		}
	}

	fmt.Println("Types:")
	for n, t := range e.names {
		fmt.Printf("\t%s: %v\n", n, t.typString(mgr))
	}

	return nil
}

func compileProgram(prog []definition) error {
	for _, d := range prog {
		err := d.compile()
		if err != nil {
			return err
		}

		defn, ok := d.(*definitionDefn)
		if !ok {
			continue
		}

		for _, i := range defn.instructions {
			fmt.Printf("%v\n", i)
		}
		fmt.Println()
	}

	return nil
}

func main() {
	l := newLexer(bufio.NewReader(os.Stdin))
	yyParse(l)
	fmt.Printf("Parsed tree:\n%v\n", l.result)
	err := typecheckProgram(l.result)
	if err != nil {
		log.Fatalln("Typecheck Error: ", err)
	}
	err = compileProgram(l.result)
	if err != nil {
		log.Fatalln("Compile Error: ", err)
	}
}
