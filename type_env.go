package main

type typEnv struct {
	names  map[string]typ
	parent *typEnv
}

func (e *typEnv) lookup(name string) typ {
	it, ok := e.names[name]
	if ok {
		return it
	}

	if e.parent != nil {
		return e.parent.lookup(name)
	}

	return nil
}

func (e *typEnv) bind(name string, r typ) {
	e.names[name] = r
}

func (e *typEnv) scope() *typEnv {
	return &typEnv{
		make(map[string]typ, 0),
		e,
	}
}
