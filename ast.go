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
		typecheck(mgr *typMgr, env *typEnv) typ
	}

	pattern interface {
		fmt.Stringer
		match(t typ, mgr *typMgr, e *typEnv)
	}

	branch struct {
		pat  pattern
		expr ast
	}

	constructor struct {
		name  string
		types []string
	}

	definition interface {
		typecheckFirst(mgr *typMgr, e *typEnv)
		typecheckSecond(mgr *typMgr, e *typEnv)
	}

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
		name       string
		params     []string
		body       ast
		paramTypes []typ
		returnType typ
	}

	definitionData struct {
		name         string
		constructors []constructor
	}
)

func newDefinitionDefn(name string, params []string, body ast) *definitionDefn {
	return &definitionDefn{
		name,
		params,
		body,
		make([]typ, 0),
		nil,
	}
}

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

//Pattern Match
func (pv patternVar) match(t typ, mgr *typMgr, e *typEnv) {
	e.bind(pv.variable, t)
}

func (pc patternConstr) match(t typ, mgr *typMgr, e *typEnv) {
	constrTyp := e.lookup(pc.constr)
	if constrTyp == nil {
		panic("Failed to lookp constructor Type")
	}

	for _, param := range pc.params {
		arr, ok := constrTyp.(typArr)
		if !ok {
			panic("Unexpected constructor type")
		}

		e.bind(param, arr.left)
		constrTyp = arr.right
	}

	mgr.unify(t, constrTyp)
	_, ok := constrTyp.(typBase)

	if !ok {
		panic("Failed to unify constructor type")
	}
}

// Type Checks
func opName(op binOpType) string {
	switch op {
	case binOpPlus:
		return "+"
	case binOpMinus:
		return "-"
	case binOpTimes:
		return "*"
	case binOpDivide:
		return "/"
	default:
		panic(fmt.Sprintf("Unsupported operator: %d", op))
		return ""
	}
}

func (a astInt) typecheck(mgr *typMgr, e *typEnv) typ {
	return &typBase{"Int"}
}

func (a astLID) typecheck(mgr *typMgr, e *typEnv) typ {
	return e.lookup(a.ID)
}

func (a astUID) typecheck(mgr *typMgr, e *typEnv) typ {
	return e.lookup(a.ID)
}

func (a astBinOp) typecheck(mgr *typMgr, e *typEnv) typ {
	ltype := a.left.typecheck(mgr, e)
	rtype := a.left.typecheck(mgr, e)
	ftype := e.lookup(opName(a.op))
	if ftype == nil {
		panic("Failed to typecheck bin op")
	}

	returnType := mgr.newTyp()
	arrowOne := &typArr{rtype, returnType}
	arrowTwo := &typArr{ltype, arrowOne}

	mgr.unify(arrowTwo, ftype)

	return returnType
}

func (a astApp) typecheck(mgr *typMgr, e *typEnv) typ {
	ltype := a.left.typecheck(mgr, e)
	rtype := a.right.typecheck(mgr, e)

	returnType := mgr.newTyp()
	arrow := &typArr{rtype, returnType}

	mgr.unify(arrow, ltype)

	return returnType
}

func (a astCase) typecheck(mgr *typMgr, e *typEnv) typ {
	caseType := a.of.typecheck(mgr, e)
	branchType := mgr.newTyp()

	for _, b := range a.branches {
		newEnv := e.scope()
		b.pat.match(caseType, mgr, newEnv)
		currBranchType := b.expr.typecheck(mgr, newEnv)
		mgr.unify(branchType, currBranchType)
	}

	return branchType
}

func (d *definitionDefn) typecheckFirst(mgr *typMgr, e *typEnv) {
	d.returnType = mgr.newTyp()
	var fullType typ = d.returnType

	for i := len(d.params) - 1; i >= 0; i-- {
		paramType := mgr.newTyp()
		fullType = &typArr{paramType, fullType}
		d.paramTypes = append(d.paramTypes, paramType)
	}

	e.bind(d.name, fullType)
}

func (d *definitionDefn) typecheckSecond(mgr *typMgr, e *typEnv) {
	newEnv := e.scope()

	if len(d.params) > len(d.paramTypes) {
		for i, pt := range d.paramTypes {
			newEnv.bind(d.params[i], pt)
		}
	} else {
		for i, p := range d.params {
			newEnv.bind(p, d.paramTypes[i])
		}
	}

	bodyType := d.body.typecheck(mgr, newEnv)
	mgr.unify(d.returnType, bodyType)
}

func (d *definitionData) typecheckFirst(mgr *typMgr, e *typEnv) {
	returnType := &typBase{d.name}

	for _, c := range d.constructors {
		var fullType typ = returnType

		for _, tn := range c.types {
			ty := &typBase{tn}
			fullType = &typArr{ty, fullType}
		}

		e.bind(c.name, fullType)
	}
}

func (d *definitionData) typecheckSecond(mgr *typMgr, e *typEnv) {}
