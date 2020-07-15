package main

import "fmt"

// Implement String
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
		panic(fmt.Errorf("Unsupported Operator: %d", a.op))
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

func (c constructor) String() string {
	return fmt.Sprintf("%s", c.name)
}

func (a astCase) String() string {
	return fmt.Sprintf("case %v of { %v }", a.of, a.branches)
}
