package main

import "fmt"

type nodeTagType = int

const (
	nodeAppTag nodeTagType = iota
	nodeNumTag
	nodeGlobalTag
	nodeIndTag
	nodeDataTag
)

type (
	node interface {
		getNodeTag() nodeTagType
	}

	stack struct {
		data []node
	}

	nodeApp struct {
		left  node
		right node
	}

	nodeNum struct {
		value int32
	}

	globalFuncPointer func(*stack) error

	nodeGlobal struct {
		arity int32
		fn    globalFuncPointer
	}

	nodeInd struct {
		next node
	}

	nodeData struct {
		tag   int8
		array []node
	}
)

func (a *nodeApp) getNodeTag() nodeTagType {
	return nodeAppTag
}

func (a *nodeNum) getNodeTag() nodeTagType {
	return nodeNumTag
}

func (a *nodeGlobal) getNodeTag() nodeTagType {
	return nodeGlobalTag
}

func (a *nodeInd) getNodeTag() nodeTagType {
	return nodeIndTag
}

func (a *nodeData) getNodeTag() nodeTagType {
	return nodeDataTag
}

func newNodeApp(l node, r node) *nodeApp {
	return &nodeApp{
		left:  l,
		right: r,
	}
}

func newNodeNum(n int32) *nodeNum {
	return &nodeNum{
		n,
	}
}

func newNodeGlobal(f globalFuncPointer, a int32) *nodeGlobal {
	return &nodeGlobal{
		arity: a,
		fn:    f,
	}
}

func newNodeInd(n node) *nodeInd {
	return &nodeInd{
		next: n,
	}
}

func newStack() *stack {
	return &stack{
		data: make([]node, 0),
	}
}

func (s *stack) stackPush(n node) {
	s.data = append(s.data, n)
}

func (s *stack) stackPop() (node, error) {
	if len(s.data) <= 0 {
		return nil, fmt.Errorf("Stack underflow")
	}

	v := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]

	return v, nil
}

func (s *stack) stackPeek(offset int) (node, error) {
	if len(s.data) <= offset {
		return nil, fmt.Errorf("Stack underflow")
	}

	v := s.data[len(s.data)-offset-1]

	return v, nil
}

func (s *stack) stackPopN(n int) error {
	if len(s.data) <= n {
		return fmt.Errorf("Stack underflow")
	}

	s.data = s.data[:len(s.data)-n-1]

	return nil
}

func (s *stack) stackSlide(n int) error {
	if len(s.data) <= n {
		return fmt.Errorf("Stack underflow")
	}

	s.data[len(s.data)-n-1] = s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-n-1]

	return nil
}

func (s *stack) stackUpdate(o int) error {
	if len(s.data) <= o+1 {
		return fmt.Errorf("Stack underflow")
	}

	ind, ok := s.data[len(s.data)-o-2].(*nodeInd)
	if !ok {
		return fmt.Errorf("Stack underflow")
	}

	ind.next = s.data[len(s.data)-1]

	return nil
}

func (s *stack) stackAlloc(o int) {
	for o != 0 {
		s.stackPush(newNodeInd(nil))
	}
}

func (s *stack) stackPack(n int, t int8) {
	data := make([]node, n)
	for i := 0; i < n; i++ {
		data[i] = s.data[len(s.data)-1-i]
	}

	dataNode := &nodeData{
		tag:   t,
		array: data,
	}

	s.stackPopN(n)
	s.stackPush(dataNode)
}

func (s *stack) stackSplit(n int) error {
	aNode, err := s.stackPop()
	if err != nil {
		return err
	}

	datNode, ok := aNode.(*nodeData)
	if !ok {
		return fmt.Errorf("This is not data node")
	}

	for i := 0; i < n; i++ {
		s.stackPush(datNode.array[i])
	}

	return nil
}

func unwind(s *stack) error {
	for {
		peek, err := s.stackPeek(0)
		if err != nil {
			return err
		}
		if appNode, ok := peek.(*nodeApp); ok {
			s.stackPush(appNode.left)
		} else if globalNode, ok := peek.(*nodeGlobal); ok {
			if len(s.data) <= int(globalNode.arity) {
				return fmt.Errorf("Arity mismatch with stack side")
			}

			for i := 1; i <= int(globalNode.arity); i++ {
				appNode, ok := s.data[len(s.data)-1-i].(*nodeApp)
				if !ok {
					return fmt.Errorf("Stack is not appNode")
				}
				s.data[len(s.data)-i] = appNode.right
			}

			globalNode.fn(s)
		} else if indNode, ok := peek.(*nodeInd); ok {
			s.stackPop()
			s.stackPush(indNode.next)
		} else {
			break
		}
	}

	return nil
}

func eval(n node) (node, error) {
	programStack := newStack()
	programStack.stackPush(n)
	err := unwind(programStack)
	if err != nil {
		return nil, err
	}
	result, err := programStack.stackPop()
	if err != nil {
		return nil, err
	}

	return result, nil
}
