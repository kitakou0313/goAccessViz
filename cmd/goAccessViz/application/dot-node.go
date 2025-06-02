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
		Node:  goNumDotNode,
		label: node.GetLabel(),
	}
}

// ToDo2度同じNodeに達した場合の対応を考える
func addDomainNodeChildrenToDotGraph(rootNode node.Node, rootDotNode *dotNode, g *simple.DirectedGraph) {
	for _, child := range rootNode.GetChildren() {
		childDotNode := newDotNode(child, g.NewNode())
		g.SetEdge(g.NewEdge(rootDotNode, childDotNode))

		addDomainNodeChildrenToDotGraph(child, childDotNode, g)
	}
}

func NewDotGraph(rootNode node.Node) *simple.DirectedGraph {
	g := simple.NewDirectedGraph()

	rootDotNode := newDotNode(rootNode, g.NewNode())
	g.AddNode(rootDotNode)
	addDomainNodeChildrenToDotGraph(rootNode, rootDotNode, g)

	return g
}
