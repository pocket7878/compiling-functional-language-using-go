package main

import (
	"fmt"
	"log"
	"reflect"
)

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
		typecheck(mgr *typMgr, env *typEnv) (typ, error)
	}

	pattern interface {
		fmt.Stringer
		match(t typ, mgr *typMgr, e *typEnv) error
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
		typecheckSecond(mgr *typMgr, e *typEnv) error
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

	unknownOpError struct {
		op binOpType
	}
)

func (e unknownOpError) Error() string {
	return fmt.Sprintf("Unknown operator: %d", e.op)
}

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
		panic(unknownOpError{a.op})
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
func (pv patternVar) match(t typ, mgr *typMgr, e *typEnv) error {
	e.bind(pv.variable, t)

	return nil
}

func (pc patternConstr) match(t typ, mgr *typMgr, e *typEnv) error {
	constrTyp := e.lookup(pc.constr)
	if constrTyp == nil {
		return fmt.Errorf("Failed to lookp constructor type: %s", pc.constr)
	}

	for _, param := range pc.params {
		arr, ok := constrTyp.(*typArr)
		if !ok {
			return fmt.Errorf("Unexpected constructor type: %v", constrTyp)
		}

		e.bind(param, arr.left)
		constrTyp = arr.right
	}

	err := mgr.unify(t, constrTyp)
	if err != nil {
		return err
	}

	log.Printf("Constructor type: %v, %v\n", constrTyp, reflect.TypeOf(constrTyp))
	_, ok := constrTyp.(*typBase)

	if !ok {
		return fmt.Errorf("Failed to unify constructor type: %v", constrTyp)
	}

	return nil
}

// Type Checks
func opName(op binOpType) (string, error) {
	switch op {
	case binOpPlus:
		return "+", nil
	case binOpMinus:
		return "-", nil
	case binOpTimes:
		return "*", nil
	case binOpDivide:
		return "/", nil
	default:
		return "", fmt.Errorf("Unsupported operator: %d", op)
	}
}

func (a astInt) typecheck(mgr *typMgr, e *typEnv) (typ, error) {
	return &typBase{"Int"}, nil
}

func (a astLID) typecheck(mgr *typMgr, e *typEnv) (typ, error) {
	return e.lookup(a.ID), nil
}

func (a astUID) typecheck(mgr *typMgr, e *typEnv) (typ, error) {
	return e.lookup(a.ID), nil
}

func (a astBinOp) typecheck(mgr *typMgr, e *typEnv) (typ, error) {
	ltype, err := a.left.typecheck(mgr, e)
	if err != nil {
		return nil, err
	}
	rtype, err := a.right.typecheck(mgr, e)
	if err != nil {
		return nil, err
	}
	o, error := opName(a.op)
	if error != nil {
		return nil, error
	}

	ftype := e.lookup(o)
	if ftype == nil {
		return nil, fmt.Errorf("Failed to typecheck bin op")
	}

	returnType := mgr.newTyp()
	arrowOne := &typArr{rtype, returnType}
	arrowTwo := &typArr{ltype, arrowOne}

	err = mgr.unify(arrowTwo, ftype)
	if err != nil {
		return nil, err
	}

	return returnType, nil
}

func (a astApp) typecheck(mgr *typMgr, e *typEnv) (typ, error) {
	ltype, err := a.left.typecheck(mgr, e)
	if err != nil {
		return nil, err
	}
	rtype, err := a.right.typecheck(mgr, e)
	if err != nil {
		return nil, err
	}

	returnType := mgr.newTyp()
	arrow := &typArr{rtype, returnType}

	err = mgr.unify(arrow, ltype)
	if err != nil {
		return nil, err
	}

	return returnType, nil
}

func (a astCase) typecheck(mgr *typMgr, e *typEnv) (typ, error) {
	caseType, err := a.of.typecheck(mgr, e)
	if err != nil {
		return nil, err
	}
	branchType := mgr.newTyp()

	for _, b := range a.branches {
		newEnv := e.scope()
		b.pat.match(caseType, mgr, newEnv)
		currBranchType, err := b.expr.typecheck(mgr, newEnv)
		if err != nil {
			return nil, err
		}
		err = mgr.unify(branchType, currBranchType)
		if err != nil {
			return nil, err
		}
	}

	return branchType, nil
}

func (d *definitionDefn) typecheckFirst(mgr *typMgr, e *typEnv) {
	d.returnType = mgr.newTyp()
	var fullType typ = d.returnType

	for i := len(d.params) - 1; i >= 0; i-- {
		paramType := mgr.newTyp()
		fullType = &typArr{paramType, fullType}
		d.paramTypes = append(d.paramTypes, paramType)
	}

	log.Println("Full type: ", fullType)

	e.bind(d.name, fullType)
}

func (d *definitionDefn) typecheckSecond(mgr *typMgr, e *typEnv) error {
	newEnv := e.scope()

	for i, p := range d.params {
		pt := d.paramTypes[len(d.paramTypes)-1-i]
		log.Printf("Nested env bind: %s to %v\n", p, pt)
		newEnv.bind(p, d.paramTypes[len(d.paramTypes)-1-i])
	}

	bodyType, err := d.body.typecheck(mgr, newEnv)
	if err != nil {
		return err
	}
	return mgr.unify(d.returnType, bodyType)
}

func (d *definitionData) typecheckFirst(mgr *typMgr, e *typEnv) {
	returnType := &typBase{d.name}

	for _, c := range d.constructors {
		var fullType typ = returnType

		for i := len(c.types) - 1; i >= 0; i-- {
			ty := &typBase{c.types[i]}
			fullType = &typArr{ty, fullType}
		}

		e.bind(c.name, fullType)
	}
}

func (d *definitionData) typecheckSecond(mgr *typMgr, e *typEnv) error {
	return nil
}
