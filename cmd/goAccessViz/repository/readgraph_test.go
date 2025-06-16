package repository

import (
	"testing"

	"goAccessViz/cmd/goAccessViz/domain/node"
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

// Tests for SQL analysis functionality
func TestExtractTablesFromSQL(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		expected []string
	}{
		{
			name:     "Simple SELECT",
			sql:      "SELECT * FROM users",
			expected: []string{"users"},
		},
		{
			name:     "SELECT with JOIN",
			sql:      "SELECT u.name, p.title FROM users u JOIN posts p ON u.id = p.user_id",
			expected: []string{"users", "posts"},
		},
		{
			name:     "INSERT statement",
			sql:      "INSERT INTO products (name, price) VALUES (?, ?)",
			expected: []string{"products"},
		},
		{
			name:     "UPDATE statement",
			sql:      "UPDATE orders SET status = ? WHERE id = ?",
			expected: []string{"orders"},
		},
		{
			name:     "DELETE statement",
			sql:      "DELETE FROM comments WHERE post_id = ?",
			expected: []string{"comments"},
		},
		{
			name:     "Multiple tables with subquery",
			sql:      "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders)",
			expected: []string{"users", "orders"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tables := extractTablesFromSQL(tt.sql)
			if len(tables) != len(tt.expected) {
				t.Errorf("Expected %d tables, got %d", len(tt.expected), len(tables))
				return
			}
			for i, expected := range tt.expected {
				if tables[i] != expected {
					t.Errorf("Expected table %s, got %s", expected, tables[i])
				}
			}
		})
	}
}

func TestDetectSQLStrings(t *testing.T) {
	testCode := `
package main

func GetUser(id int) (*User, error) {
	query := "SELECT * FROM users WHERE id = ?"
	return db.Query(query, id)
}

func CreatePost(title string) error {
	sql := "INSERT INTO posts (title) VALUES (?)"
	_, err := db.Exec(sql, title)
	return err
}
`

	sqlStrings := detectSQLStrings(testCode)
	expectedCount := 2
	if len(sqlStrings) != expectedCount {
		t.Errorf("Expected %d SQL strings, got %d", expectedCount, len(sqlStrings))
	}

	expectedTables := map[string]bool{"users": true, "posts": true}
	allTables := make(map[string]bool)
	for _, sql := range sqlStrings {
		tables := extractTablesFromSQL(sql)
		for _, table := range tables {
			allTables[table] = true
		}
	}

	for expectedTable := range expectedTables {
		if !allTables[expectedTable] {
			t.Errorf("Expected to find table %s", expectedTable)
		}
	}
}

func TestCreateDBTableNodes(t *testing.T) {
	sqlStrings := []string{
		"SELECT * FROM users",
		"SELECT * FROM posts WHERE user_id = ?",
		"INSERT INTO comments (post_id, content) VALUES (?, ?)",
	}

	dbNodes := createDBTableNodes(sqlStrings)

	expectedTables := map[string]bool{"users": true, "posts": true, "comments": true}
	if len(dbNodes) != len(expectedTables) {
		t.Errorf("Expected %d DB table nodes, got %d", len(expectedTables), len(dbNodes))
	}

	foundTables := make(map[string]bool)
	for _, dbNode := range dbNodes {
		foundTables[dbNode.GetLabel()] = true
	}

	for expectedTable := range expectedTables {
		if !foundTables[expectedTable] {
			t.Errorf("Expected to find DB table node for %s", expectedTable)
		}
	}
}

func TestReadGraphWithSQLAnalysis(t *testing.T) {
	// This test will verify that ReadGraph includes SQL table analysis
	// We'll need a test package with SQL strings
	nodes, err := ReadGraph("goAccessViz/testpkg")
	if err != nil {
		t.Fatalf("Failed to read graph with SQL analysis: %v", err)
	}

	// Count different node types
	functionNodes := 0
	dbTableNodes := 0
	var dbTableNames []string

	for _, n := range nodes {
		switch dbNode := n.(type) {
		case *node.FunctionNode:
			functionNodes++
		case *node.DBTableNode:
			dbTableNodes++
			dbTableNames = append(dbTableNames, dbNode.GetLabel())
		}
	}

	// Should have function nodes (existing functionality)
	if functionNodes == 0 {
		t.Error("Expected to find function nodes")
	}

	// Should also have DB table nodes (since we added SQL strings to testpkg)
	if dbTableNodes == 0 {
		t.Error("Expected to find DB table nodes from SQL analysis")
	}

	// Check for specific tables we expect
	expectedTables := map[string]bool{"users": true, "posts": true, "orders": true}
	foundTables := make(map[string]bool)
	for _, tableName := range dbTableNames {
		foundTables[tableName] = true
	}

	for expectedTable := range expectedTables {
		if !foundTables[expectedTable] {
			t.Errorf("Expected to find table %s but it was not detected", expectedTable)
		}
	}

	t.Logf("Found %d function nodes and %d DB table nodes", functionNodes, dbTableNodes)
	t.Logf("DB tables found: %v", dbTableNames)
}
