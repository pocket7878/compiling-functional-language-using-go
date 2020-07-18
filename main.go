//go:generate goyacc -l -o parser.go parser.y
package main

import (
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
		d.typecheckFirst(mgr, e)
	}

	for _, d := range prg {
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

func runProgram(prog []definition) {
	// Boot G-Machine VM
	vm := newGVM()
	//Store every function to heap
	for _, d := range prog {
		switch def := d.(type) {
		case *definitionDefn:
			vm.addGlobal(def.name, len(def.params), def.instructions)
		case *definitionData:
			for _, c := range def.constructors {
				packInsts := make([]inst, 0)
				packInsts = append(packInsts, instPack{c.tag, len(c.types)})
				packInsts = append(packInsts, instUpdate{})
				packInsts = append(packInsts, instUnwind{})
				vm.addGlobal(c.name, len(c.types), packInsts)
			}
		}
	}

	vm.pushInst(&instEval{})
	vm.pushInst(&instPushGlobal{"main"})

	vm.run()

	resultAddr := vm.stack.pop()
	resultNode, ok := vm.heap[resultAddr]
	if !ok {
		log.Fatal("Failed to retrieve resutl")
	}

	log.Println("The Result is : ", resultNode)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage %s <file>\n", os.Args[0])
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	l := newLexer(file)
	yyParse(l)
	err = typecheckProgram(l.result)
	if err != nil {
		log.Fatalln("Typecheck Error: ", err)
	}
	err = compileProgram(l.result)
	if err != nil {
		log.Fatalln("Compile Error: ", err)
	}

	runProgram(l.result)
}
