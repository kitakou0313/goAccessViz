package node

// 関数に相当するNode
type FunctionTrackedEntity struct {
	funtionName string
	children    []TrackedEntity
}

func NewFunctionTrackedEntity(functionName string, children []TrackedEntity) *FunctionTrackedEntity {
	return &FunctionTrackedEntity{
		funtionName: functionName,
		children:    children,
	}
}

func (fn *FunctionTrackedEntity) GetChildren() []TrackedEntity {
	return fn.children
}

func (fn *FunctionTrackedEntity) GetLabel() string {
	return fn.funtionName
}
