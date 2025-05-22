package node

import (
	"testing"
)

func TestGetChildren(t *testing.T) {
	testChildren := []nod.FunctionNode{}
	testChildren := FunctionNode.Ne
	nodeInstance := domain.NewFunctionNode("doTestFunction")
	actual := nodeInstance.GetChildren()

	expected := []Node{}

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Errorf("Expected %v, but got %v", expected[i], actual[i])
		}
	}

}
