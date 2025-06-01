package application

import (
	"goAccessViz/cmd/goAccessViz/domain/node"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type dotNode struct {
	graph.Node
	label string
}

func (d *dotNode) DOTID() string {
	return d.label
}

func newDotNode(node node.Node, goNumDotNode graph.Node) *dotNode {
	return &dotNode{
		Node:  nil,
		label: node.GetLabel(),
	}
}

func NewDotGraph(rootNode node.Node) *simple.DirectedGraph {
	g := simple.NewDirectedGraph()

}
