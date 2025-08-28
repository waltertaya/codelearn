package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Challenge struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Difficulty  string    `json:"difficulty"`
	Language    string    `json:"language"`
	TestCases   string    `json:"test_cases"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Submission struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ChallengeID int       `json:"challenge_id"`
	Code        string    `json:"code"`
	Language    string    `json:"language"`
	Status      string    `json:"status"` // pending, passed, failed
	Score       int       `json:"score"`
	Output      string    `json:"output"`
	CreatedAt   time.Time `json:"created_at"`
}

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./codelearn.db")
	if err != nil {
		return err
	}

	// Test connection
	if err = db.Ping(); err != nil {
		return err
	}

	// Create tables
	if err = createTables(); err != nil {
		return err
	}

	// Insert sample data
	if err = insertSampleData(); err != nil {
		return err
	}

	return nil
}

func createTables() error {
	// Users table
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Challenges table
	challengeTable := `
	CREATE TABLE IF NOT EXISTS challenges (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		difficulty TEXT NOT NULL,
		language TEXT NOT NULL,
		test_cases TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Submissions table
	submissionTable := `
	CREATE TABLE IF NOT EXISTS submissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		challenge_id INTEGER NOT NULL,
		code TEXT NOT NULL,
		language TEXT NOT NULL,
		status TEXT DEFAULT 'pending',
		score INTEGER DEFAULT 0,
		output TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id),
		FOREIGN KEY (challenge_id) REFERENCES challenges (id)
	);`

	tables := []string{userTable, challengeTable, submissionTable}
	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return err
		}
	}

	return nil
}

func insertSampleData() error {
	// Check if challenges already exist
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM challenges").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Sample data already exists
	}

	// Sample challenges
	challenges := []struct {
		title       string
		description string
		difficulty  string
		language    string
		testCases   string
	}{
		{
			title:       "Two Sum",
			description: "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.",
			difficulty:  "Easy",
			language:    "python",
			testCases:   `[{"input": "[2,7,11,15], 9", "expected": "[0,1]"}, {"input": "[3,2,4], 6", "expected": "[1,2]"}]`,
		},
		{
			title:       "Reverse String",
			description: "Write a function that reverses a string. The input string is given as an array of characters s.",
			difficulty:  "Easy",
			language:    "javascript",
			testCases:   `[{"input": "['h','e','l','l','o']", "expected": "['o','l','l','e','h']"}, {"input": "['H','a','n','n','a','h']", "expected": "['h','a','n','n','a','H']"}]`,
		},
		{
			title:       "Binary Search",
			description: "Given an array of integers nums which is sorted in ascending order, and an integer target, write a function to search target in nums.",
			difficulty:  "Medium",
			language:    "go",
			testCases:   `[{"input": "[-1,0,3,5,9,12], 9", "expected": "4"}, {"input": "[-1,0,3,5,9,12], 2", "expected": "-1"}]`,
		},
		{
			title:       "Valid Parentheses",
			description: "Given a string s containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid.",
			difficulty:  "Easy",
			language:    "python",
			testCases:   `[{"input": "()", "expected": "true"}, {"input": "()[]{}", "expected": "true"}, {"input": "(]", "expected": "false"}]`,
		},
		{
			title:       "Fibonacci Sequence",
			description: "Write a function to generate the nth Fibonacci number.",
			difficulty:  "Easy",
			language:    "javascript",
			testCases:   `[{"input": "0", "expected": "0"}, {"input": "1", "expected": "1"}, {"input": "10", "expected": "55"}]`,
		},
	}

	for _, challenge := range challenges {
		_, err := db.Exec(`
			INSERT INTO challenges (title, description, difficulty, language, test_cases)
			VALUES (?, ?, ?, ?, ?)
		`, challenge.title, challenge.description, challenge.difficulty, challenge.language, challenge.testCases)
		if err != nil {
			return err
		}
	}

	return nil
}
