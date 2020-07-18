package main

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
