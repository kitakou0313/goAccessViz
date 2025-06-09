package application

import (
	"goAccessViz/cmd/goAccessViz/domain/node"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

func TestNew(t *testing.T) {}

// DotNodeのDOTIDメソッドのテスト
func TestDOTID(t *testing.T) {
	expected := "test.dot.node"

	dotNode := &dotNode{
		label: expected,
		Node:  nil,
	}

	actual := dotNode.DOTID()

	if actual != expected {
		t.Errorf("Expected DOTID to be '%s', but got '%s'", expected, actual)
	}
}

func TestGetLabel(t *testing.T) {
	expected := "test.dot.node"

	dotNode := &dotNode{
		label: expected,
		Node:  nil,
	}

	actual := dotNode.Getlabel()

	if actual != expected {
		t.Errorf("Expected Getlabel res to be '%s', but got '%s'", expected, actual)
	}
}

// DotNodeの生成メソッドに対してのテスト
func TestNewDotNode(t *testing.T) {
	testNodeName := "testFunction"
	testNode := node.NewFunctionNode(testNodeName, nil)

	tmpGraph := simple.NewDirectedGraph()
	tmpGoNumDotNode := tmpGraph.NewNode()

	dotNode := newDotNode(testNode, tmpGoNumDotNode)

	if dotNode.DOTID() != testNodeName {
		t.Errorf("Expected DOTID to be '%s', but got '%s'", testNodeName, dotNode.DOTID())
	}
}

// node DomainオブジェクトからDotNodeのGraphを生成するメソッドのテスト
// TODO:　網羅性を考慮したテストを書く
func TestNewDotGraphWithTreeGraph(t *testing.T) {
	testNodeName := "testFunction"
	testChildrenNodes := []node.Node{
		node.NewFunctionNode("childFunctionNode", nil),
		node.NewDBTableNode("childDBNode", nil),
	}
	testRootNode := node.NewFunctionNode(testNodeName, testChildrenNodes)

	dotGrapth := NewDotGraph([]node.Node{testRootNode})

	// トポロジカルソートして確認
	sortedDotNodes, err := topo.Sort(dotGrapth)
	if err != nil {
		t.Errorf("Failed to sort graph: %v", err)
	}

	// ルートノードが最初に来ることを確認
	if testRootNode.GetLabel() != sortedDotNodes[0].(*dotNode).DOTID() {
		t.Errorf("Expected root node '%s' to be first, but got '%s'", testRootNode.GetLabel(), sortedDotNodes[0].(*dotNode).DOTID())

	}

	// 子ノードが正しく追加されていることを確認
	if len(sortedDotNodes)-1 != len(testChildrenNodes) {
		t.Errorf("Expected %d child nodes, but got %d", len(testChildrenNodes), len(sortedDotNodes)-1)
	}

	childDomainNodesExistingInDotGraph := make(map[string]bool)
	for _, childDomainNode := range testChildrenNodes {
		childDomainNodesExistingInDotGraph[childDomainNode.GetLabel()] = false
		for _, childDotNode := range sortedDotNodes[1:] {
			childDomainNodesExistingInDotGraph[childDotNode.(*dotNode).DOTID()] = true
		}
	}

	for _, childDomainNode := range testChildrenNodes {
		if _, ok := childDomainNodesExistingInDotGraph[childDomainNode.GetLabel()]; !ok {
			t.Errorf("Child node '%s' not found in dot graph", childDomainNode.GetLabel())
		}
	}
}

// 一つのNodeが二つ入力辺を持つ場合のテスト 例:A -> C, B -> C
func TestNewDotGraphWithSome2IncomingEdges(t *testing.T) {
	nodeHaving2IncomingEdges := node.NewFunctionNode("nodeHaving2IncomingEdges", nil)
	testChildrenNodes := []node.Node{
		nodeHaving2IncomingEdges,
		node.NewFunctionNode("childFunctionNode", []node.Node{nodeHaving2IncomingEdges}),
		node.NewDBTableNode("childDBNode", nil),
	}
	testNodeName := "testFunction"
	testRootNode := node.NewFunctionNode(testNodeName, testChildrenNodes)

	dotGrapth := NewDotGraph([]node.Node{testRootNode})

	// トポロジカルソートして確認
	sortedDotNodes, err := topo.Sort(dotGrapth)
	if err != nil {
		t.Errorf("Failed to sort graph: %v", err)
	}

	// ルートノードが最初に来ることを確認
	if testRootNode.GetLabel() != sortedDotNodes[0].(*dotNode).DOTID() {
		t.Errorf("Expected root node '%s' to be first, but got '%s'", testRootNode.GetLabel(), sortedDotNodes[0].(*dotNode).DOTID())

	}

	// 子ノードが正しく追加されていることを確認
	if len(sortedDotNodes)-1 != 3 {
		t.Errorf("Expected %d child nodes, but got %d", 3, len(sortedDotNodes)-1)
	}

	childDomainNodesExistingInDotGraph := make(map[string]bool)
	for _, childDomainNode := range testChildrenNodes {
		childDomainNodesExistingInDotGraph[childDomainNode.GetLabel()] = false
		for _, childDotNode := range sortedDotNodes[1:] {
			childDomainNodesExistingInDotGraph[childDotNode.(*dotNode).DOTID()] = true
		}
	}

	for _, childDomainNode := range testChildrenNodes {
		if _, ok := childDomainNodesExistingInDotGraph[childDomainNode.GetLabel()]; !ok {
			t.Errorf("Child node '%s' not found in dot graph", childDomainNode.GetLabel())
		}
	}
}

func validIfDotNodeIsInSline(domainNode node.Node, gonumNodesList []graph.Node) bool {
	for _, goNumNode := range gonumNodesList {
		if domainNode.GetLabel() == goNumNode.(*dotNode).DOTID() {
			return true
		}
	}
	return false

}

func TestNewDotGraphWithSomeRootNodes(t *testing.T) {
	testNodeName1 := "testFunction1"
	testChildrenNodes := []node.Node{
		node.NewFunctionNode("childFunctionNode", nil),
		node.NewDBTableNode("childDBNode", nil),
	}
	testRootNode1 := node.NewFunctionNode(testNodeName1, testChildrenNodes)
	testNodeName2 := "testFunction2"
	testRootNode2 := node.NewFunctionNode(testNodeName2, testChildrenNodes)

	dotGrapth := NewDotGraph([]node.Node{testRootNode1, testRootNode2})

	// トポロジカルソートして確認
	sortedDotNodes, err := topo.Sort(dotGrapth)
	if err != nil {
		t.Errorf("Failed to sort graph: %v", err)
	}

	// ルートノードが最初に来ることを確認
	if validIfDotNodeIsInSline(testRootNode1, sortedDotNodes[:2]) == false {
		t.Errorf("Expected root node '%s' to be first, but got '%s'", testRootNode1.GetLabel(), sortedDotNodes[:2])

	}
	if validIfDotNodeIsInSline(testRootNode2, sortedDotNodes[:2]) == false {
		t.Errorf("Expected root node '%s' to be first, but got '%s'", testRootNode2.GetLabel(), sortedDotNodes[:2])

	}

	// 子ノードが正しく追加されていることを確認
	if len(sortedDotNodes)-2 != len(testChildrenNodes) {
		t.Errorf("Expected %d child nodes, but got %d", len(testChildrenNodes), len(sortedDotNodes)-2)
	}

	childDomainNodesExistingInDotGraph := make(map[string]bool)
	for _, childDomainNode := range testChildrenNodes {
		childDomainNodesExistingInDotGraph[childDomainNode.GetLabel()] = false
		for _, childDotNode := range sortedDotNodes[2:] {
			childDomainNodesExistingInDotGraph[childDotNode.(*dotNode).DOTID()] = true
		}
	}

	for _, childDomainNode := range testChildrenNodes {
		if _, ok := childDomainNodesExistingInDotGraph[childDomainNode.GetLabel()]; !ok {
			t.Errorf("Child node '%s' not found in dot graph", childDomainNode.GetLabel())
		}
	}
}
