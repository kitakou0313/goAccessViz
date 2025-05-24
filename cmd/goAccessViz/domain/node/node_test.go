package node

import (
	"testing"
)

func TestGetChildren(t *testing.T) {
	testChildren := []Node{
		NewFunctionNode("testFunction1", []Node{}),
		NewFunctionNode("testFunction2", []Node{}),
	}
	nodeInstance := NewFunctionNode("doTestFunction", testChildren)
	actual := nodeInstance.GetChildren()

	expected := testChildren

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Errorf("Expected %v, but got %v", expected[i], actual[i])
		}
	}

}
