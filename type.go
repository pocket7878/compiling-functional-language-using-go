package main

import (
	"fmt"
	"log"
	"reflect"
)

type (
	typStringer interface {
		fmt.Stringer
		typString(m *typMgr) string
	}

	typ interface {
		typStringer
	}

	typVar struct {
		name string
	}

	typBase struct {
		name string
	}

	typArr struct {
		left  typ
		right typ
	}

	typMgr struct {
		lastId int
		types  map[string]typ
	}
)

func newTypMgr() *typMgr {
	return &typMgr{
		0,
		make(map[string]typ, 0),
	}
}

func (m *typMgr) newTypName() string {
	tmp := m.lastId
	m.lastId++
	str := ""
	for tmp != -1 {
		str += string(rune('a' + (tmp % 26)))
		tmp = tmp/26 - 1
	}

	log.Println("New type Name: ", str)

	return str
}

func (m *typMgr) newTyp() *typVar {
	return &typVar{
		m.newTypName(),
	}
}

func (m *typMgr) newArrTyp() *typArr {
	return &typArr{
		m.newTyp(),
		m.newTyp(),
	}
}

func (m *typMgr) resolve(t typ, v **typVar) typ {
	log.Printf("\t\tResolve: %v typ: %v\n", t, reflect.TypeOf(t))
	for {
		cast, ok := t.(*typVar)
		log.Printf("\t\tt is typVar? = %v\n", ok)
		if !ok {
			break
		}

		it, ok := m.types[cast.name]
		if !ok {
			log.Printf("\t\tfound base variable: %v\n", cast)
			*v = cast
			break
		}
		t = it
	}

	return t
}

func (m *typMgr) bind(s string, t typ) {
	log.Printf("\tbind %s to %v (%v)\n", s, t, reflect.TypeOf(t))
	other, ok := t.(*typVar)

	if ok && other.name == s {
		log.Printf("\t\tSkip binding because already same name")
		return
	}

	m.types[s] = t
}

func (m *typMgr) unify(l typ, r typ) {
	log.Printf("Unify %v with %v\n", l, r)
	var lvar *typVar
	var rvar *typVar

	l = m.resolve(l, &lvar)
	r = m.resolve(r, &rvar)

	log.Printf("\tresolved lvar: %v\n", lvar)
	log.Printf("\tresolved rvar: %v\n", rvar)

	if lvar != nil {
		m.bind(lvar.name, r)
		return
	} else if rvar != nil {
		m.bind(rvar.name, l)
		return
	}

	larr, larrOk := l.(*typArr)
	rarr, rarrOk := r.(*typArr)

	if larrOk && rarrOk {
		m.unify(larr.left, rarr.left)
		m.unify(larr.right, rarr.right)
		return
	}

	lid, lidOk := l.(*typBase)
	rid, ridOk := l.(*typBase)

	if lidOk && ridOk {
		if lid.name == rid.name {
			return
		}
	}

	panic(fmt.Sprintf("Failed to unify: %v with %v", l, r))
}

// Print type
func (v typVar) String() string {
	return fmt.Sprintf("TypVar(%s)", v.name)
}

func (b typBase) String() string {
	return b.name
}

func (a typArr) String() string {
	switch a.right.(type) {
	case *typArr:
		return fmt.Sprintf("%v -> (%v)", a.left, a.right)
	default:
		return fmt.Sprintf("%v -> %v", a.left, a.right)
	}
}

func (v typVar) typString(m *typMgr) string {
	it, ok := m.types[v.name]
	if ok {
		return it.typString(m)
	}
	return fmt.Sprintf("TypVar(%s)", v.name)
}

func (b typBase) typString(m *typMgr) string {
	return b.name
}

func (a typArr) typString(m *typMgr) string {
	switch a.right.(type) {
	case *typArr:
		return fmt.Sprintf("%v -> (%v)", a.left.typString(m), a.right.typString(m))
	default:
		return fmt.Sprintf("%v -> %v", a.left.typString(m), a.right.typString(m))
	}
}
