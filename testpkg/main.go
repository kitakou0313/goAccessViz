package testpkg

import (
	"github.com/jmoiron/sqlx"
)

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

// User represents a user record
type User struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// Post represents a post record
type Post struct {
	ID     int    `db:"id"`
	Title  string `db:"title"`
	UserID int    `db:"user_id"`
}

// Order represents an order record
type Order struct {
	ID     int    `db:"id"`
	Status string `db:"status"`
}

// Functions using sqlx for testing SQL analysis
func GetUser(db *sqlx.DB, id int) (*User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE id = ?", id)
	return &user, err
}

func GetAllUsers(db *sqlx.DB) ([]User, error) {
	var users []User
	err := db.Select(&users, "SELECT * FROM users")
	return users, err
}

func CreatePost(db *sqlx.DB, title string, userID int) error {
	_, err := db.Exec("INSERT INTO posts (title, user_id, created_at) VALUES (?, ?, NOW())", title, userID)
	return err
}

func UpdateOrder(db *sqlx.DB, id int, status string) error {
	_, err := db.Exec("UPDATE orders SET status = ? WHERE id = ?", status, id)
	return err
}

func GetUserPosts(db *sqlx.DB, userID int) ([]Post, error) {
	var posts []Post
	err := db.Select(&posts, "SELECT p.id, p.title, p.user_id FROM users u JOIN posts p ON u.id = p.user_id WHERE u.id = ?", userID)
	return posts, err
}

func DeletePost(db *sqlx.DB, postID int) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", postID)
	return err
}

func GetUsersByStatus(db *sqlx.DB, status string) ([]User, error) {
	var users []User
	err := db.Select(&users, "SELECT u.* FROM users u JOIN orders o ON u.id = o.user_id WHERE o.status = ?", status)
	return users, err
}
