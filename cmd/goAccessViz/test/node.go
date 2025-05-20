package test

import "testing"

func TestGetChildren(t *testing.T) {
	node := &Node{}
	actual := node.GetChildren()

	expected := []Node{}

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Errorf("Expected %v, but got %v", expected[i], actual[i])
		}
	}

}
