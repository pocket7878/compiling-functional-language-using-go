package main

import (
	"fmt"
	"log"
	"reflect"
)

func (a astInt) compile(e compEnv, into *[]inst) error {
	*into = append(*into, instPushInt{a.value})

	return nil
}

func (a astLID) compile(e compEnv, into *[]inst) error {
	if e.hasVariable(a.ID) {
		idOffset, err := e.getOffset(a.ID)
		if err != nil {
			return err
		}
		*into = append(*into, instPush{idOffset})
	} else {
		*into = append(*into, instPushGlobal{a.ID})
	}

	return nil
}

func (a astUID) compile(e compEnv, into *[]inst) error {
	*into = append(*into, instPushGlobal{a.ID})

	return nil
}

func (a astBinOp) compile(e compEnv, into *[]inst) error {
	err := a.right.compile(e, into)
	if err != nil {
		return err
	}
	err = a.left.compile(compEnvOffset{1, e}, into)
	if err != nil {
		return err
	}

	opA, err := opAction(a.op)
	if err != nil {
		return err
	}

	*into = append(*into, instPushGlobal{opA})
	*into = append(*into, instMkApp{})
	*into = append(*into, instMkApp{})

	return nil
}

func (a astApp) compile(e compEnv, into *[]inst) error {
	err := a.right.compile(e, into)
	if err != nil {
		return err
	}

	err = a.left.compile(compEnvOffset{1, e}, into)
	if err != nil {
		return err
	}

	*into = append(*into, instMkApp{})

	return nil
}

func (b branch) compile(e compEnv, into *[]inst) error {
	panic("Unreachable Code")
}

func (b constructor) compile(e compEnv, into *[]inst) error {
	panic("Unreachable Code")
}

func (a astCase) compile(e compEnv, into *[]inst) error {
	ty, ok := a.of.getNodeType().(*typData)

	if !ok {
		return fmt.Errorf("Unexpected case of type %s", reflect.TypeOf(a.of))
	}

	err := a.of.compile(e, into)
	if err != nil {
		return err
	}

	*into = append(*into, instEval{})

	jmpInst := instJump{}
	jmpInst.tagMappings = make(map[int]int, 0)

	for _, branch := range a.branches {
		branchInst := make([]inst, 0)
		log.Printf("Branch type: %v\n", reflect.TypeOf(branch.pat))
		if _, ok := branch.pat.(patternVar); ok {
			branch.expr.compile(compEnvOffset{1, e}, &branchInst)

			for _, constPair := range ty.constructors {
				if _, ok := jmpInst.tagMappings[constPair.tag]; ok {
					break
				}

				jmpInst.tagMappings[constPair.tag] = len(jmpInst.branches)
			}

			jmpInst.branches = append(jmpInst.branches, branchInst)
		} else if cpat, ok := branch.pat.(patternConstr); ok {
			newEnv := e
			for i := len(cpat.params) - 1; i >= 0; i-- {
				newEnv = compEnvVar{cpat.params[i], newEnv}
			}

			branchInst = append(branchInst, instSplit{})
			branch.expr.compile(newEnv, &branchInst)
			branchInst = append(branchInst, instSlide{len(cpat.params)})

			newTag := ty.constructors[cpat.constr].tag

			if _, ok := jmpInst.tagMappings[newTag]; ok {
				return fmt.Errorf("Technically not a type error: duplicate pattern")
			}

			jmpInst.tagMappings[newTag] = len(jmpInst.branches)
			jmpInst.branches = append(jmpInst.branches, branchInst)
		}
	}

	for _, constPair := range ty.constructors {
		if _, ok := jmpInst.tagMappings[constPair.tag]; !ok {
			return fmt.Errorf("Non total pattern")
		}
	}

	*into = append(*into, jmpInst)

	return nil
}

func (a *definitionData) compile() error {
	return nil
}

func (a *definitionDefn) compile() error {
	var newEnv compEnv = compEnvOffset{0, nil}
	for i := len(a.params) - 1; i >= 0; i-- {
		newEnv = compEnvVar{a.params[i], newEnv}
	}
	err := a.body.compile(newEnv, &a.instructions)
	if err != nil {
		return err
	}
	a.instructions = append(a.instructions, instUpdate{len(a.params)})
	a.instructions = append(a.instructions, instPop{len(a.params)})

	return nil
}
