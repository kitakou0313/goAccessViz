package application

import "testing"

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

// node DomainオブジェクトからDotNodeのGraphを生成するメソッドのテスト
