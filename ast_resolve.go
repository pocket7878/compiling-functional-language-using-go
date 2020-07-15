package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func resolveCommon(a ast, mgr *typMgr) error {
	var v *typVar
	resolvedType := mgr.resolve(a.getNodeType(), &v)
	if v != nil {
		return fmt.Errorf("Ambiguously typed program in resolveCommon: %v", a)
	}

	err := a.resolve(mgr)
	if err != nil {
		return err
	}
	a.setNodeType(resolvedType)

	return nil
}

func (a *astBinOp) resolve(mgr *typMgr) error {
	err := resolveCommon(a.left, mgr)
	if err != nil {
		return errors.Wrap(err, "resolve astBinOp left")
	}

	err = resolveCommon(a.right, mgr)
	if err != nil {
		return errors.Wrap(err, "resolve astBinOp right")
	}

	return nil
}

func (a *astApp) resolve(mgr *typMgr) error {
	err := resolveCommon(a.left, mgr)
	if err != nil {
		return errors.Wrap(err, "resolve astApp left")
	}

	err = resolveCommon(a.right, mgr)
	if err != nil {
		return errors.Wrap(err, "resolve astApp righ")
	}

	return nil
}

func (a *astInt) resolve(mgr *typMgr) error {
	return nil
}

func (a *astLID) resolve(mgr *typMgr) error {
	return nil
}

func (a *astUID) resolve(mgr *typMgr) error {
	return nil
}

func (a *astCase) resolve(mgr *typMgr) error {
	err := resolveCommon(a.of, mgr)
	if err != nil {
		return errors.Wrap(err, "resolve astCase")
	}

	for i := 0; i < len(a.branches); i++ {
		err = resolveCommon(&a.branches[i], mgr)
		if err != nil {
			return errors.Wrap(err, "resolve astCase branch")
		}
	}

	return nil
}

func (b *branch) resolve(mgr *typMgr) error {
	err := resolveCommon(b.expr, mgr)
	if err != nil {
		return errors.Wrap(err, "resolve branch")
	}

	return nil
}

func (b *constructor) resolve(mgr *typMgr) error {
	return nil
}

func (a *definitionDefn) resolve(mgr *typMgr) error {
	var v *typVar
	err := resolveCommon(a.body, mgr)
	if err != nil {
		return errors.Wrap(err, "resolve definitionDefn")
	}

	a.returnType = mgr.resolve(a.returnType, &v)

	if v != nil {
		return fmt.Errorf("Ambiguously typed program in resolve")
	}

	for i := 0; i < len(a.paramTypes); i++ {
		a.paramTypes[i] = mgr.resolve(a.paramTypes[i], &v)
		if v != nil {
			return fmt.Errorf("Ambiguously typed program in resolve")
		}
	}

	return nil
}

func (a *definitionData) resolve(mgr *typMgr) error {
	return nil
}
