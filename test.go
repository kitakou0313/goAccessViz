package main

import (
	"os"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

type CustomNode struct {
	graph.Node
	Label string
}

func (n *CustomNode) DOTID() string {
	return n.Label
}

func main() {
	g := simple.NewDirectedGraph()

	n1 := &CustomNode{
		Node:  g.NewNode(),
		Label: "input",
	}
	g.AddNode(n1)
	n2 := &CustomNode{
		Node:  g.NewNode(),
		Label: "Process",
	}
	g.AddNode(n2)
	n3 := &CustomNode{
		Node:  g.NewNode(),
		Label: "OutPut",
	}
	g.AddNode(n3)

	g.SetEdge(g.NewEdge(n1, n2))
	g.SetEdge(g.NewEdge(n2, n3))

	b, _ := dot.Marshal(g, "Graph", "", " ")
	os.WriteFile("graph.dot", b, 0644)
}
