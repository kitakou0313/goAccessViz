package node

type GraphNode interface {
	GetChildren() []GraphNode
	GetLabel() string
}
