//go:generate goyacc -l -o parser.go parser.y
package main

import (
	"fmt"
	"log"

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

func f_add(s *stack) error {
	stackTop, err := s.stackPeek(0)
	if err != nil {
		return err
	}
	aLeft, err := eval(stackTop)
	if err != nil {
		return err
	}
	leftNum, ok := aLeft.(*nodeNum)
	if !ok {
		return fmt.Errorf("Left Node is not a number")
	}

	stackSecond, err := s.stackPeek(1)
	if err != nil {
		return err
	}
	aRight, err := eval(stackSecond)
	if err != nil {
		return err
	}
	rightNum, ok := aRight.(*nodeNum)
	if !ok {
		return fmt.Errorf("Right Node is not a number")
	}

	s.stackPush(newNodeNum(leftNum.value + rightNum.value))

	return nil
}

func f_main(s *stack) error {
	s.stackPush(newNodeNum(320))
	s.stackPush(newNodeNum(6))
	s.stackPush(newNodeGlobal(f_add, 2))

	left, err := s.stackPop()
	if err != nil {
		return err
	}
	right, err := s.stackPop()
	if err != nil {
		return err
	}
	s.stackPush(newNodeApp(left, right))

	left, err = s.stackPop()
	if err != nil {
		return err
	}
	right, err = s.stackPop()
	if err != nil {
		return err
	}
	s.stackPush(newNodeApp(left, right))

	return nil
}

func main() {
	firstNode := newNodeGlobal(f_main, 0)
	result, err := eval(firstNode)
	if err != nil {
		log.Fatalln(err)
	}

	resultNum, ok := result.(*nodeNum)
	if !ok {
		log.Fatalln("Result is not a number: ", result)
	}
	fmt.Println(resultNum.value)
	/*
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
	*/

}
