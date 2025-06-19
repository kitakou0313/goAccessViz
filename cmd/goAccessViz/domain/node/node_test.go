package node

import (
	"testing"
)

var testChildren = []TrackedEntity{
	NewFunctionTrackedEntity("testFunction1", []TrackedEntity{}),
	NewFunctionTrackedEntity("testFunction2", []TrackedEntity{}),
}

func TestFunctionNodeGetChildren(t *testing.T) {
	nodeInstance := NewFunctionTrackedEntity("doTestFunction", testChildren)
	actual := nodeInstance.GetChildren()

	expected := testChildren

	if len(actual) != len(expected) {
		t.Errorf("Expected %d children, but got %d", len(expected), len(actual))
	}

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Errorf("Expected %v, but got %v", expected[i], actual[i])
		}
	}

}

func TestNetFunctionNode(t *testing.T) {
	expectedName := "testFunction"
	functionNode := NewFunctionTrackedEntity(expectedName, testChildren)
	if functionNode.GetLabel() != "testFunction" {
		t.Errorf("Expected function name '%s', but got '%s'", expectedName, functionNode.funtionName)
	}

}

func TestDBTableNode(t *testing.T) {
	expectedName := "test-table"

	dbtableNode := NewDatabaseTableTrackedEntity(expectedName, testChildren)

	if dbtableNode.GetLabel() != expectedName {
		t.Errorf("Expected table name '%s', but got '%s'", expectedName, dbtableNode.tableName)
	}

	for i := 0; i < len(testChildren); i++ {
		if dbtableNode.children[i] != testChildren[i] {
			t.Errorf("Expected child %d to be %v, but got %v", i, testChildren[i], dbtableNode.children[i])
		}
	}
}

func TestDBTableNodeGetChildren(t *testing.T) {
	dbTableNode := NewDatabaseTableTrackedEntity("test-table", testChildren)
	actual := dbTableNode.GetChildren()

	expected := testChildren

	if len(actual) != len(expected) {
		t.Errorf("Expected %d children, but got %d", len(expected), len(actual))
	}

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Errorf("Expected %v, but got %v", expected[i], actual[i])
		}
	}

}
