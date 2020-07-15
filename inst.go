package main

import (
	"fmt"
)

type (
	inst interface {
		fmt.Stringer
	}

	instPushInt struct {
		value int
	}

	instPushGlobal struct {
		name string
	}

	instPush struct {
		offset int
	}

	instPop struct {
		count int
	}

	instMkApp struct{}

	instUpdate struct {
		offset int
	}

	instPack struct {
		tag  int
		size int
	}

	instSplit struct{}

	instJump struct {
		branches    [][]inst
		tagMappings map[int]int
	}

	instSlide struct {
		offset int
	}

	instBinOp struct {
		op binOpType
	}

	instEval struct{}

	instAlloc struct {
		amount int
	}

	instUnwind struct{}
)

func (i instPushInt) String() string {
	return fmt.Sprintf("PushInt(%d)", i.value)
}

func (i instPushGlobal) String() string {
	return fmt.Sprintf("PushGlobal(%s)", i.name)
}

func (i instPush) String() string {
	return fmt.Sprintf("Push(%d)", i.offset)
}

func (i instPop) String() string {
	return fmt.Sprintf("Pop(%d)", i.count)
}

func (i instMkApp) String() string {
	return fmt.Sprintf("MkApp()")
}

func (i instUpdate) String() string {
	return fmt.Sprintf("Update(%d)", i.offset)
}

func (i instPack) String() string {
	return fmt.Sprintf("Pack(%d, %d)", i.tag, i.size)
}

func (i instSplit) String() string {
	return fmt.Sprintf("Split()")
}

func (i instJump) String() string {
	//TODO: Update Code to display Jump
	result := "Jump(\n"
	for _, inss := range i.branches {
		for _, ins := range inss {
			result += fmt.Sprintf("\t%v\n", ins)
		}
		result += "\n"
	}
	result += ")"

	return result
}

func (i instSlide) String() string {
	return fmt.Sprintf("Slide(%d)", i.offset)
}

func (i instBinOp) String() string {
	o, err := opName(i.op)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("BinOp(%s)", o)
}

func (i instEval) String() string {
	return fmt.Sprintf("Eval()")
}

func (i instAlloc) String() string {
	return fmt.Sprintf("Alloc(%d)", i.amount)
}

func (i instUnwind) String() string {
	return fmt.Sprintf("Unwind()")
}
