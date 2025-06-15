package repository

import (
	"testing"
)

func TestReadGraphWithTestPackage(t *testing.T) {
	// Test reading our custom test package
	nodes, err := ReadGraph("goAccessViz/testpkg")
	if err != nil {
		t.Fatalf("Failed to read graph: %v", err)
	}

	// Check that we have nodes
	if len(nodes) == 0 {
		t.Error("Expected graph to have nodes, but it is empty")
	}

	// Check that at least one node has children (function calls)
	hasChildren := false
	for _, node := range nodes {
		if len(node.GetChildren()) > 0 {
			hasChildren = true
			break
		}
	}
	
	if !hasChildren {
		t.Error("Expected at least one node to have children (edges)")
	}

	// Verify we have the expected functions
	functionNames := make(map[string]bool)
	for _, node := range nodes {
		functionNames[node.GetLabel()] = true
	}

	// We should have our test functions
	expectedFunctions := []string{"FunctionA", "FunctionB", "FunctionC", "FunctionD"}
	foundCount := 0
	for _, expected := range expectedFunctions {
		for label := range functionNames {
			if contains(label, expected) {
				foundCount++
				break
			}
		}
	}

	if foundCount == 0 {
		t.Error("Expected to find test functions in the graph")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    (len(s) > len(substr) && 
		     s[len(s)-len(substr):] == substr))
}