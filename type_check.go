package main

import (
	"fmt"
)

// Type Checks
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
	ltype, err := typeCheckCommon(a.left, mgr, e)
	if err != nil {
		return nil, err
	}
	rtype, err := typeCheckCommon(a.right, mgr, e)
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

func (a branch) typecheck(mgr *typMgr, e *typEnv) (typ, error) {
	panic("Unreachable code")
}

func (a astApp) typecheck(mgr *typMgr, e *typEnv) (typ, error) {
	ltype, err := typeCheckCommon(a.left, mgr, e)
	if err != nil {
		return nil, err
	}
	rtype, err := typeCheckCommon(a.right, mgr, e)
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
	caseType, err := typeCheckCommon(a.of, mgr, e)
	if err != nil {
		return nil, err
	}
	branchType := mgr.newTyp()

	for _, b := range a.branches {
		newEnv := e.scope()
		b.pat.match(caseType, mgr, newEnv)
		currBranchType, err := typeCheckCommon(b.expr, mgr, newEnv)
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

	e.bind(d.name, fullType)
}

func (d *definitionDefn) typecheckSecond(mgr *typMgr, e *typEnv) error {
	newEnv := e.scope()

	for i, p := range d.params {
		pt := d.paramTypes[len(d.paramTypes)-1-i]
		newEnv.bind(p, pt)
	}

	bodyType, err := typeCheckCommon(d.body, mgr, newEnv)
	if err != nil {
		return err
	}
	return mgr.unify(d.returnType, bodyType)
}

func (d *definitionData) typecheckFirst(mgr *typMgr, e *typEnv) {
	thisType := typData{
		typBase{
			d.name,
		},
		make(map[string]typDataConstr, 0),
	}
	returnType := &thisType
	nextTag := 0

	for _, c := range d.constructors {
		c.tag = nextTag
		thisType.constructors[c.name] = typDataConstr{nextTag + 1}
		nextTag++

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

func typeCheckCommon(a ast, m *typMgr, e *typEnv) (typ, error) {
	typ, err := a.typecheck(m, e)
	if err != nil {
		return nil, err
	}
	a.setNodeType(typ)

	return typ, nil
}
