package domain

type Node interface {
	GetChildren() []Node
}
