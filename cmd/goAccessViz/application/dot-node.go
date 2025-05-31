package application

import (
	"goAccessViz/cmd/goAccessViz/domain/node"

	"gonum.org/v1/gonum/graph"
)

type dotNode struct {
	graph.Node
	label string
}

func (d *dotNode) DOTID() string {
	return d.label
}

func NewDotNode(node node.Node) *dotNode {
	return &dotNode{
		Node:  nil,
		label: node.GetLabel(),
	}
}
