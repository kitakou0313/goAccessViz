package application

import (
	"goAccessViz/cmd/goAccessViz/domain/node"
	"testing"

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

// DotNodeの生成メソッドに対してのテスト
func TestNewDotNode(t *testing.T) {
	testNodeName := "testFunction"
	testNode := node.NewFunctionNode(testNodeName, nil)
	dotNode := newDotNode(testNode)

	if dotNode.DOTID() != testNodeName {
		t.Errorf("Expected DOTID to be '%s', but got '%s'", testNodeName, dotNode.DOTID())
	}
}

// node DomainオブジェクトからDotNodeのGraphを生成するメソッドのテスト
// TODO:　網羅性を考慮したテストを書く
func TestNewDotGraph(t *testing.T) {
	testNodeName := "testFunction"
	testChildrenNodes := []node.Node{
		node.NewFunctionNode("childFunctionNode", nil),
		node.NewDBTableNode("childDBNode", nil),
	}
	testRootNode := node.NewFunctionNode(testNodeName, testChildrenNodes)

	dotGrapth := NewDotGraph(testRootNode)

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
