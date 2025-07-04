package application

import (
	"goAccessViz/cmd/goAccessViz/domain/node"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

type dotNode struct {
	graph.Node
	label string
}

func (d *dotNode) DOTID() string {
	return d.label
}

func (d *dotNode) Getlabel() string {
	return d.DOTID()
}

func newDotNode(node node.TrackedEntity, goNumDotNode graph.Node) *dotNode {
	return &dotNode{
		Node:  goNumDotNode,
		label: node.GetLabel(),
	}
}

// ToDo2度同じNodeに達した場合の対応を考える
func addDomainNodeChildrenToDotGraph(rootNode node.TrackedEntity, rootDotNode *dotNode, g *simple.DirectedGraph, dotIdToIDMap map[string]*dotNode) {
	for _, child := range rootNode.GetChildren() {
		label := child.GetLabel()

		childDotNode, exists := dotIdToIDMap[label]
		if !exists {
			childDotNode = newDotNode(child, g.NewNode())
			dotIdToIDMap[label] = childDotNode
		}

		g.SetEdge(g.NewEdge(rootDotNode, childDotNode))
		addDomainNodeChildrenToDotGraph(child, childDotNode, g, dotIdToIDMap)
	}
}

func NewDotGraph(rootNodes []node.TrackedEntity) *simple.DirectedGraph {
	g := simple.NewDirectedGraph()

	dotIdToDotNodeMap := make(map[string]*dotNode)
	for _, rootNode := range rootNodes {
		rootDotNode := newDotNode(rootNode, g.NewNode())
		g.AddNode(rootDotNode)
		dotIdToDotNodeMap[rootNode.GetLabel()] = rootDotNode
		addDomainNodeChildrenToDotGraph(rootNode, rootDotNode, g, dotIdToDotNodeMap)
	}

	return g
}

func ConvertDotGraphToString(dotGraph *simple.DirectedGraph) (string, error) {
	b, err := dot.Marshal(dotGraph, "Graph", "", " ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
