package node

import (
	"testing"
)

var testChildren = []Node{
	NewFunctionNode("testFunction1", []Node{}),
	NewFunctionNode("testFunction2", []Node{}),
}

func TestFunctionNodeGetChildren(t *testing.T) {
	nodeInstance := NewFunctionNode("doTestFunction", testChildren)
	actual := nodeInstance.GetChildren()

	expected := testChildren

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Errorf("Expected %v, but got %v", expected[i], actual[i])
		}
	}

}

func TestNetFunctionNode(t *testing.T) {
	expectedName := "testFunction"
	functionNode := NewFunctionNode(expectedName, testChildren)
	if functionNode.funtionName != "testFunction" {
		t.Errorf("Expected function name '%s', but got '%s'", expectedName, functionNode.funtionName)
	}

}
