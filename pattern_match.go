package main

import (
	"fmt"
)

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

	_, ok := constrTyp.(*typBase)

	if !ok {
		return fmt.Errorf("Failed to unify constructor type: %v", constrTyp)
	}

	return nil
}
