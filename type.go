package main

import (
	"fmt"
	"log"
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

	typDataConstr struct {
		tag int
	}

	typData struct {
		typBase
		constructors map[string]typDataConstr
	}

	typArr struct {
		left  typ
		right typ
	}

	typMgr struct {
		lastID int
		types  map[string]typ
	}

	unificationError struct {
		left  typ
		right typ
	}
)

// Error
func (e unificationError) Error() string {
	return fmt.Sprintf("Failed to unify type: %s with %s", e.left, e.right)
}

func newTypMgr() *typMgr {
	return &typMgr{
		0,
		make(map[string]typ, 0),
	}
}

func (m *typMgr) newTypName() string {
	tmp := m.lastID
	m.lastID++
	str := ""
	for tmp != -1 {
		str += string(rune('a' + (tmp % 26)))
		tmp = tmp/26 - 1
	}

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

// Type Check
func (m *typMgr) resolve(t typ, v **typVar) typ {
	for {
		cast, ok := t.(*typVar)
		if !ok {
			break
		}

		it, ok := m.types[cast.name]
		if !ok {
			*v = cast
			break
		}
		t = it
	}

	return t
}

func (m *typMgr) bind(s string, t typ) {
	other, ok := t.(*typVar)

	if ok && other.name == s {
		return
	}

	log.Printf("Bind %s to %v\n", s, t)

	m.types[s] = t
}

func (m *typMgr) unify(l typ, r typ) error {
	var lvar *typVar
	var rvar *typVar

	l = m.resolve(l, &lvar)
	r = m.resolve(r, &rvar)

	if lvar != nil {
		m.bind(lvar.name, r)
		return nil
	} else if rvar != nil {
		m.bind(rvar.name, l)
		return nil
	}

	larr, larrOk := l.(*typArr)
	rarr, rarrOk := r.(*typArr)

	if larrOk && rarrOk {
		err := m.unify(larr.left, rarr.left)
		if err != nil {
			return err
		}
		return m.unify(larr.right, rarr.right)
	}

	lbase, lbaseOk := l.(*typBase)
	ldata, ldataOk := l.(*typData)
	rbase, rbaseOk := r.(*typBase)
	rdata, rdataOk := r.(*typData)

	if lbaseOk && rbaseOk && lbase.name == rbase.name {
		return nil
	}

	if lbaseOk && rdataOk && lbase.name == rdata.name {
		return nil
	}

	if ldataOk && rbaseOk && ldata.name == rbase.name {
		return nil
	}

	if ldataOk && rdataOk && ldata.name == rdata.name {
		return nil
	}

	return unificationError{l, r}
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

func (a typData) String() string {
	return fmt.Sprintf("TypData(%v)", a.name)
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
