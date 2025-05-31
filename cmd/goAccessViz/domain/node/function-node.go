package node

// 関数に相当するNode
type FunctionNode struct {
	funtionName string
	children    []Node
}

func NewFunctionNode(functionName string, children []Node) *FunctionNode {
	return &FunctionNode{
		funtionName: functionName,
		children:    children,
	}
}

func (fn *FunctionNode) GetChildren() []Node {
	return fn.children
}

func (fn *FunctionNode) GetLabel() string {
	return fn.funtionName
}
