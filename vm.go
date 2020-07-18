package main

import (
	"fmt"
)

type addrType = int
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
		fmt.Stringer
		getNodeTag() nodeTagType
	}

	nodeApp struct {
		left  addrType
		right addrType
	}

	nodeNum struct {
		value int
	}

	nodeGlobal struct {
		arity int
		code  []inst
	}

	nodeInd struct {
		next addrType
	}

	nodeData struct {
		tag   int
		array []addrType
	}

	stack struct {
		data []addrType
	}

	dumpEntry struct {
		insts []inst
		stack stack
	}

	gVM struct {
		insts     []inst
		stack     *stack
		dump      []dumpEntry
		heap      map[int]node
		globalMap map[string]int
		freeAddr  int
	}
)

func newGVM() *gVM {
	return &gVM{
		insts:     make([]inst, 0),
		stack:     newStack(),
		dump:      make([]dumpEntry, 0),
		heap:      make(map[int]node),
		globalMap: make(map[string]int),
		freeAddr:  0,
	}
}

func (g gVM) String() string {
	result := fmt.Sprintf("------ VM -----------\n")
	result += fmt.Sprintf("inst: %v\n", g.insts)
	result += fmt.Sprintf("stack: %v\n", g.stack)
	result += fmt.Sprintf("dump: %v\n", g.dump)
	result += fmt.Sprintf("heap: %v\n", g.heap)
	result += fmt.Sprintf("globalMap: %v\n", g.globalMap)
	result += fmt.Sprintf("======================\n")

	return result
}

func (g *gVM) pushInst(i inst) {
	g.insts = append(g.insts, i)
}

func (g *gVM) addGlobal(name string, arity int, inst []inst) {
	a := g.newFreeAddr()
	g.heap[a] = &nodeGlobal{arity, inst}
	g.globalMap[name] = a
}

func (g *gVM) run() {
	for {
		if len(g.insts) == 0 {
			break
		}
		head := g.peekInst()
		fmt.Println("------- START --------------")
		fmt.Println("BEFORE VM")
		fmt.Printf("%v", g)
		fmt.Printf("Execute: %v\n", head)
		head.execute(g)
		fmt.Println("AFTER VM")
		fmt.Printf("%v", g)
		fmt.Printf("------- END --------------\n\n")
	}
}

func (g *gVM) newFreeAddr() int {
	a := g.freeAddr
	g.freeAddr++

	return a
}

func (g *gVM) peekInst() inst {
	return g.insts[len(g.insts)-1]
}

func (g *gVM) popInst() {
	g.insts = g.insts[:len(g.insts)-1]
}

func (i instPushInt) execute(g *gVM) {
	g.popInst()
	a := g.newFreeAddr()
	g.heap[a] = &nodeNum{i.value}
	g.stack.push(a)
}

func (inst instPushGlobal) execute(g *gVM) {
	g.popInst()
	a, ok := g.globalMap[inst.name]
	if !ok {
		panic(fmt.Errorf("Undefined function: %s", inst.name))
	}
	g.stack.push(a)
}

func (inst instPush) execute(g *gVM) {
	g.popInst()
	an := g.stack.peek(inst.offset)
	g.stack.push(an)
}

func (inst instMkApp) execute(g *gVM) {
	g.popInst()
	a0 := g.stack.pop()
	a1 := g.stack.pop()
	a := g.newFreeAddr()
	g.heap[a] = &nodeApp{a0, a1}
	g.stack.push(a)
}

func (i instUnwind) execute(g *gVM) {
	a := g.stack.peek(0)
	heapNode := g.heap[a]
	if heapNode == nil {
		panic("Unexpected nil")
	}
	fmt.Printf("\tUnwinding: %v\n", heapNode)
	if appNode, ok := heapNode.(*nodeApp); ok {
		g.stack.push(appNode.left)
	} else if globalNode, ok := heapNode.(*nodeGlobal); ok {
		for i := 1; i <= int(globalNode.arity); i++ {
			ak := g.stack.data[len(g.stack.data)-1-i]
			appNode, ok := g.heap[ak].(*nodeApp)
			if !ok {
				panic("Stack is not appNode")
			}
			g.stack.data[len(g.stack.data)-i] = appNode.right
		}

		newInst := make([]inst, 0)
		for i := len(globalNode.code) - 1; i >= 0; i-- {
			newInst = append(newInst, globalNode.code[i])
		}
		g.insts = newInst
	} else if indNode, ok := heapNode.(*nodeInd); ok {
		g.stack.pop()
		g.stack.push(indNode.next)
	} else if len(g.dump) > 0 {
		a = g.stack.pop()
		de := g.dump[len(g.dump)-1]
		g.dump = g.dump[:len(g.dump)-1]
		g.insts = de.insts
		g.stack = &de.stack
		g.stack.push(a)
	}

	fmt.Println("\tAfter Unwinding:")
	fmt.Printf("\t%v\n", g)
}

func (i instUpdate) execute(g *gVM) {
	g.popInst()
	a := g.stack.pop()
	an := g.stack.peek(i.offset)
	g.heap[an] = &nodeInd{a}
}

func (ins instPack) execute(g *gVM) {
	g.popInst()
	a := g.newFreeAddr()
	arr := make([]addrType, ins.size)
	for i := 0; i < ins.size; i++ {
		arr[i] = g.stack.pop()
	}
	g.heap[a] = &nodeData{ins.tag, arr}
	g.stack.push(a)
}

func (ins instSplit) execute(g *gVM) {
	g.popInst()
	a := g.stack.pop()
	dataNode, ok := g.heap[a].(*nodeData)
	if !ok {
		panic("Expected data node")
	}
	for i := len(dataNode.array) - 1; i >= 0; i-- {
		g.stack.push(dataNode.array[i])
	}
}

func (ins instJump) execute(g *gVM) {
	g.popInst()
	a := g.stack.peek(0)
	dataNode, ok := g.heap[a].(*nodeData)
	if !ok {
		panic("Expected data node")
	}
	for i := len(ins.branches[dataNode.tag]) - 1; i >= 0; i-- {
		g.insts = append(g.insts, ins.branches[dataNode.tag][i])
	}
}

func (ins instSlide) execute(g *gVM) {
	g.popInst()
	a0 := g.stack.pop()
	for i := 1; i <= ins.offset; i++ {
		g.stack.pop()
	}
	g.stack.push(a0)
}

func (ins instBinOp) execute(g *gVM) {
	g.popInst()
	a0 := g.stack.pop()
	n, ok := g.heap[a0].(*nodeNum)
	if !ok {
		panic("Not a number")
	}
	a1 := g.stack.pop()
	m, ok := g.heap[a1].(*nodeNum)
	if !ok {
		panic("Not a number")
	}

	result := &nodeNum{0}

	switch ins.op {
	case binOpPlus:
		result.value = n.value + m.value
	case binOpMinus:
		result.value = n.value - m.value
	case binOpTimes:
		result.value = n.value * m.value
	case binOpDivide:
		result.value = n.value / m.value
	default:
		panic("Unsupported BinOp")
	}

	a := g.newFreeAddr()
	g.heap[a] = result
	g.stack.push(a)
}

func (ins instEval) execute(g *gVM) {
	g.popInst()
	newInst := make([]inst, 1)
	newInst[0] = &instUnwind{}
	a := g.stack.pop()
	newStack := newStack()
	newStack.push(a)

	de := dumpEntry{
		g.insts,
		*g.stack,
	}

	g.dump = append(g.dump, de)
	g.insts = newInst
	g.stack = newStack
}

func (ins instAlloc) execute(g *gVM) {
	g.popInst()
	for i := 0; i < ins.amount; i++ {
		ak := g.newFreeAddr()
		g.heap[ak] = &nodeInd{}
	}
}

func (ins instPop) execute(g *gVM) {
	g.popInst()
	for i := 0; i < ins.count; i++ {
		g.stack.pop()
	}
}
