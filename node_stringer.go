package main

import "fmt"

func (a nodeApp) String() string {
	return fmt.Sprintf("NApp %d %d", a.left, a.right)
}

func (a nodeNum) String() string {
	return fmt.Sprintf("NInt %d", a.value)
}

func (a nodeGlobal) String() string {
	return fmt.Sprintf("NGlobal %d c", a.arity)
}

func (a nodeInd) String() string {
	return fmt.Sprintf("NInd %d", a.next)
}

func (a nodeData) String() string {
	return fmt.Sprintf("NData %d %v", a.tag, a.array)
}
