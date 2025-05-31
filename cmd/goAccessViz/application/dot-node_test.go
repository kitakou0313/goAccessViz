package application

import (
	"goAccessViz/cmd/goAccessViz/domain/node"
	"testing"
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
	dotNode := NewDotNode(testNode)

	if dotNode.DOTID() != testNodeName {
		t.Errorf("Expected DOTID to be '%s', but got '%s'", testNodeName, dotNode.DOTID())
	}
}

// node DomainオブジェクトからDotNodeのGraphを生成するメソッドのテスト
