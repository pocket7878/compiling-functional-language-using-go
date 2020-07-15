package main

import (
	"fmt"
)

type (
	ast interface {
		fmt.Stringer
		setNodeType(t typ)
		getNodeType() typ
		typecheck(mgr *typMgr, env *typEnv) (typ, error)
		compile(env compEnv, into *[]inst) error
		resolve(mgr *typMgr) error
	}

	pattern interface {
		fmt.Stringer
		match(t typ, mgr *typMgr, e *typEnv) error
	}

	branch struct {
		pat     pattern
		expr    ast
		nodeTyp typ
	}

	constructor struct {
		name    string
		types   []string
		tag     int
		nodeTyp typ
	}

	definition interface {
		typecheckFirst(mgr *typMgr, e *typEnv)
		typecheckSecond(mgr *typMgr, e *typEnv) error
		resolve(mgr *typMgr) error
		compile() error
	}

	astInt struct {
		value   int
		nodeTyp typ
	}

	astLID struct {
		ID      string
		nodeTyp typ
	}

	astUID struct {
		ID      string
		nodeTyp typ
	}

	astBinOp struct {
		op      binOpType
		left    ast
		right   ast
		nodeTyp typ
	}

	astApp struct {
		left    ast
		right   ast
		nodeTyp typ
	}

	astCase struct {
		of       ast
		branches []branch
		nodeTyp  typ
	}

	patternVar struct {
		variable string
		nodeTyp  typ
	}

	patternConstr struct {
		constr  string
		params  []string
		nodeTyp typ
	}

	definitionDefn struct {
		name         string
		params       []string
		body         ast
		paramTypes   []typ
		returnType   typ
		nodeTyp      typ
		instructions []inst
	}

	definitionData struct {
		name         string
		constructors []constructor
		nodeTyp      typ
	}
)

func newDefinitionDefn(name string, params []string, body ast) *definitionDefn {
	return &definitionDefn{
		name,
		params,
		body,
		make([]typ, 0),
		nil,
		nil,
		make([]inst, 0),
	}
}
