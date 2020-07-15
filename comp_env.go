package main

import "fmt"

type (
	compEnv interface {
		getOffset(name string) (int, error)
		hasVariable(name string) bool
	}

	compEnvVar struct {
		name   string
		parent compEnv
	}

	compEnvOffset struct {
		offset int
		parent compEnv
	}
)

func (e compEnvVar) getOffset(name string) (int, error) {
	if e.name == name {
		return 0, nil
	}

	if e.parent != nil {
		pOffset, err := e.parent.getOffset(name)
		if err != nil {
			return 0, err
		}
		return pOffset + 1, nil
	}

	return 0, fmt.Errorf("Failed to getOffset")
}

func (e compEnvVar) hasVariable(name string) bool {
	if name == e.name {
		return true
	}

	if e.parent != nil {
		return e.parent.hasVariable(name)
	}

	return false
}

func (e compEnvOffset) getOffset(name string) (int, error) {
	if e.parent != nil {
		pOffset, err := e.parent.getOffset(name)
		if err != nil {
			return 0, err
		}

		return pOffset + e.offset, nil
	}

	return 0, fmt.Errorf("Failed to getOffset")
}
func (e compEnvOffset) hasVariable(name string) bool {
	if e.parent != nil {
		return e.parent.hasVariable(name)
	}
	return false
}
