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

	dbNodesMap := createDBTableNodesMap(sqlStrings)

	expectedTables := map[string]bool{"users": true, "posts": true, "comments": true}
	if len(dbNodesMap) != len(expectedTables) {
		t.Errorf("Expected %d DB table nodes, got %d", len(expectedTables), len(dbNodesMap))
	}

	foundTables := make(map[string]bool)
	for tableName := range dbNodesMap {
		foundTables[tableName] = true
	}

	for expectedTable := range expectedTables {
		if !foundTables[expectedTable] {
			t.Errorf("Expected to find DB table node for %s", expectedTable)
		}
	}
}

func TestReadGraphWithSQLAnalysis(t *testing.T) {
	// This test will verify that ReadGraph includes SQL table analysis
	// In the new implementation, DB tables are children of functions, not top-level nodes
	nodes, err := ReadGraph("goAccessViz/testpkg")
	if err != nil {
		t.Fatalf("Failed to read graph with SQL analysis: %v", err)
	}

	// Count function nodes and find DB table nodes as children
	functionNodes := 0
	var dbTableNames []string

	for _, n := range nodes {
		if fnNode, ok := n.(*node.FunctionNode); ok {
			functionNodes++
			
			// Check children for DB table nodes
			for _, child := range fnNode.GetChildren() {
				if dbNode, ok := child.(*node.DBTableNode); ok {
					dbTableNames = append(dbTableNames, dbNode.GetLabel())
				}
			}
		}
	}

	// Should have function nodes (existing functionality)
	if functionNodes == 0 {
		t.Error("Expected to find function nodes")
	}

	// Should also have DB table nodes as children (since we added SQL strings to testpkg)
	if len(dbTableNames) == 0 {
		t.Error("Expected to find DB table nodes as children of functions")
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

	t.Logf("Found %d function nodes and %d DB table relationships", functionNodes, len(dbTableNames))
	t.Logf("DB tables found: %v", dbTableNames)
}

func TestFunctionToTableRelationships(t *testing.T) {
	// Test that functions that query SQL tables have those tables as children
	nodes, err := ReadGraph("goAccessViz/testpkg")
	if err != nil {
		t.Fatalf("Failed to read graph: %v", err)
	}


	// Find functions that should have SQL table children
	var getUserFunc, createPostFunc, updateOrderFunc *node.FunctionNode
	
	for _, n := range nodes {
		if fnNode, ok := n.(*node.FunctionNode); ok {
			label := fnNode.GetLabel()
			if contains(label, "GetUser") {
				getUserFunc = fnNode
			} else if contains(label, "CreatePost") {
				createPostFunc = fnNode
			} else if contains(label, "UpdateOrder") {
				updateOrderFunc = fnNode
			}
		}
	}

	// Test GetUser function should have 'users' table as child
	if getUserFunc != nil {
		hasUsersTable := false
		for _, child := range getUserFunc.GetChildren() {
			if dbNode, ok := child.(*node.DBTableNode); ok && dbNode.GetLabel() == "users" {
				hasUsersTable = true
				break
			}
		}
		if !hasUsersTable {
			t.Error("GetUser function should have 'users' table as child node")
		}
	} else {
		t.Error("Could not find GetUser function")
	}

	// Test CreatePost function should have 'posts' table as child
	if createPostFunc != nil {
		hasPostsTable := false
		for _, child := range createPostFunc.GetChildren() {
			if dbNode, ok := child.(*node.DBTableNode); ok && dbNode.GetLabel() == "posts" {
				hasPostsTable = true
				break
			}
		}
		if !hasPostsTable {
			t.Error("CreatePost function should have 'posts' table as child node")
		}
	} else {
		t.Error("Could not find CreatePost function")
	}

	// Test UpdateOrder function should have 'orders' table as child
	if updateOrderFunc != nil {
		hasOrdersTable := false
		for _, child := range updateOrderFunc.GetChildren() {
			if dbNode, ok := child.(*node.DBTableNode); ok && dbNode.GetLabel() == "orders" {
				hasOrdersTable = true
				break
			}
		}
		if !hasOrdersTable {
			t.Error("UpdateOrder function should have 'orders' table as child node")
		}
	} else {
		t.Error("Could not find UpdateOrder function")
	}

	// Test that GetUserPosts function has both 'users' and 'posts' tables as children
	var getUserPostsFunc *node.FunctionNode
	for _, n := range nodes {
		if fnNode, ok := n.(*node.FunctionNode); ok {
			if contains(fnNode.GetLabel(), "GetUserPosts") {
				getUserPostsFunc = fnNode
				break
			}
		}
	}

	if getUserPostsFunc != nil {
		foundTables := make(map[string]bool)
		for _, child := range getUserPostsFunc.GetChildren() {
			if dbNode, ok := child.(*node.DBTableNode); ok {
				foundTables[dbNode.GetLabel()] = true
			}
		}
		
		if !foundTables["users"] {
			t.Error("GetUserPosts function should have 'users' table as child node")
		}
		if !foundTables["posts"] {
			t.Error("GetUserPosts function should have 'posts' table as child node")
		}
	} else {
		t.Error("Could not find GetUserPosts function")
	}
}
