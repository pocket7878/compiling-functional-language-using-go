package main

import "fmt"

type binOpType = int

const (
	binOpPlus binOpType = iota
	binOpMinus
	binOpTimes
	binOpDivide
)

type (
	ast interface {
		fmt.Stringer
	}

	pattern interface {
		ast
	}

	branch struct {
		pat  pattern
		expr ast
	}

	constructor struct {
		name  string
		types []string
	}

	definition interface{}

	astInt struct {
		value int
	}

	astLID struct {
		ID string
	}

	astUID struct {
		ID string
	}

	astBinOp struct {
		op    binOpType
		left  ast
		right ast
	}

	astApp struct {
		left  ast
		right ast
	}

	astCase struct {
		of       ast
		branches []branch
	}

	patternVar struct {
		variable string
	}

	patternConstr struct {
		constr string
		params []string
	}

	definitionDefn struct {
		name   string
		params []string
		body   ast
	}

	definitionData struct {
		name         string
		constructors []constructor
	}
)

func (d definitionDefn) String() string {
	return fmt.Sprintf("defn %s %v { %v }", d.name, d.params, d.body.String())
}

func (d definitionData) String() string {
	return fmt.Sprintf("data %s [%v]", d.name, d.constructors)
}

func (a astBinOp) String() string {
	switch a.op {
	case binOpPlus:
		return fmt.Sprintf("%v + %v", a.left, a.right)
	case binOpMinus:
		return fmt.Sprintf("%v - %v", a.left, a.right)
	case binOpTimes:
		return fmt.Sprintf("%v * %v", a.left, a.right)
	case binOpDivide:
		return fmt.Sprintf("%v / %v", a.left, a.right)
	default:
		panic("Usupported operator")
	}
}

func (a astInt) String() string {
	return fmt.Sprintf("%d", a.value)
}

func (a astApp) String() string {
	return fmt.Sprintf("%s(%v)", a.left, a.right)
}

func (a astLID) String() string {
	return a.ID
}

func (a astUID) String() string {
	return a.ID
}

func (pv patternVar) String() string {
	return pv.variable
}

func (b branch) String() string {
	return fmt.Sprintf("%v -> %v", b.pat, b.expr)
}

func (pc patternConstr) String() string {
	return fmt.Sprintf("%s(%v)", pc.constr, pc.params)
}

func (a astCase) String() string {
	return fmt.Sprintf("case %v of { %v }", a.of, a.branches)
}
