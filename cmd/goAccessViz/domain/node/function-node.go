package node

// 関数に相当するNode
type FunctionGraphNode struct {
	funtionName string
	children    []GraphNode
}

func NewFunctionGraphNode(functionName string, children []GraphNode) *FunctionGraphNode {
	return &FunctionGraphNode{
		funtionName: functionName,
		children:    children,
	}
}

func (fn *FunctionGraphNode) GetChildren() []GraphNode {
	return fn.children
}

func (fn *FunctionGraphNode) GetLabel() string {
	return fn.funtionName
}
