package main

import (
	"goAccessViz/cmd/goAccessViz/application"
	"goAccessViz/cmd/goAccessViz/domain/node"
	"os"

	"gonum.org/v1/gonum/graph/encoding/dot"
)

func main() {
	testNodeName1 := "testFunction1"
	testChildrenNodes := []node.Node{
		node.NewFunctionNode("childFunctionNode", nil),
		node.NewDBTableNode("childDBNode", nil),
	}
	testRootNode1 := node.NewFunctionNode(testNodeName1, testChildrenNodes)
	testNodeName2 := "testFunction2"
	testRootNode2 := node.NewFunctionNode(testNodeName2, testChildrenNodes)

	dotGrapth := application.NewDotGraph([]node.Node{testRootNode1, testRootNode2})

	b, _ := dot.Marshal(dotGrapth, "Graph", "", " ")
	os.WriteFile("graph.dot", b, 0644)

}
