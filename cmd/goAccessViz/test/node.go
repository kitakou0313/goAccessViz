package test

import "testing"

func TestGetChildren(t *testing.T) {
	testChildren := []FunctionNode{}
	node := NewFunctionNode("doTestFunction")
	actual := node.GetChildren()

	expected := []Node{}

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Errorf("Expected %v, but got %v", expected[i], actual[i])
		}
	}

}
