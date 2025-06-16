package testpkg

// Simple test functions for ReadGraph testing
func FunctionA() {
	FunctionB()
	FunctionC()
}

func FunctionB() {
	FunctionC()
}

func FunctionC() {
	// Base function with no calls
}

func FunctionD() {
	FunctionA()
}

// Functions with SQL strings for testing SQL analysis
func GetUser(id int) {
	query := "SELECT * FROM users WHERE id = ?"
	// Simulate database query
	_ = query
}

func CreatePost(title string) {
	sql := "INSERT INTO posts (title, created_at) VALUES (?, NOW())"
	_ = sql
}

func UpdateOrder(id int, status string) {
	updateSQL := "UPDATE orders SET status = ? WHERE id = ?"
	_ = updateSQL
}

func GetUserPosts(userID int) {
	joinQuery := "SELECT u.name, p.title FROM users u JOIN posts p ON u.id = p.user_id WHERE u.id = ?"
	_ = joinQuery
}
