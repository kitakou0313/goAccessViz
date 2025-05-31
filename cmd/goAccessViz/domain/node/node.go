package node

type Node interface {
	GetChildren() []Node
	GetLabel() string
}
