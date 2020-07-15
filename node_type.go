package main

// Set Node Type
func (a *astInt) setNodeType(t typ) {
	a.nodeTyp = t
}

func (a *astBinOp) setNodeType(t typ) {
	a.nodeTyp = t
}

func (a *astApp) setNodeType(t typ) {
	a.nodeTyp = t
}

func (a *astLID) setNodeType(t typ) {
	a.nodeTyp = t
}

func (a *astUID) setNodeType(t typ) {
	a.nodeTyp = t
}

func (pv *patternVar) setNodeType(t typ) {
	pv.nodeTyp = t
}

func (b *branch) setNodeType(t typ) {
	b.nodeTyp = t
}

func (c *constructor) setNodeType(t typ) {
	c.nodeTyp = t
}

func (pc *patternConstr) setNodeType(t typ) {
	pc.nodeTyp = t
}

func (a *astCase) setNodeType(t typ) {
	a.nodeTyp = t
}

// Get NodeType
func (a astInt) getNodeType() typ {
	return a.nodeTyp
}

func (a astBinOp) getNodeType() typ {
	return a.nodeTyp
}

func (a astApp) getNodeType() typ {
	return a.nodeTyp
}

func (a astLID) getNodeType() typ {
	return a.nodeTyp
}

func (a astUID) getNodeType() typ {
	return a.nodeTyp
}

func (pv patternVar) getNodeType() typ {
	return pv.nodeTyp
}

func (b branch) getNodeType() typ {
	return b.nodeTyp
}

func (c constructor) getNodeType() typ {
	return c.nodeTyp
}

func (pc patternConstr) getNodeType() typ {
	return pc.nodeTyp
}

func (a astCase) getNodeType() typ {
	return a.nodeTyp
}
