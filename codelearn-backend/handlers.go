package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SubmitSolutionRequest struct {
	Code     string `json:"code" binding:"required"`
	Language string `json:"language" binding:"required"`
}

type LeaderboardEntry struct {
	Username     string `json:"username"`
	TotalScore   int    `json:"total_score"`
	Submissions  int    `json:"submissions"`
	LastActivity string `json:"last_activity"`
}

func GetChallengesHandler(c *gin.Context) {
	// Get query parameters for filtering
	difficulty := c.Query("difficulty")
	language := c.Query("language")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	// Build query
	query := "SELECT id, title, description, difficulty, language, created_at, updated_at FROM challenges WHERE 1=1"
	args := []interface{}{}

	if difficulty != "" {
		query += " AND difficulty = ?"
		args = append(args, difficulty)
	}

	if language != "" {
		query += " AND language = ?"
		args = append(args, language)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch challenges"})
		return
	}
	defer rows.Close()

	var challenges []Challenge
	for rows.Next() {
		var challenge Challenge
		err := rows.Scan(&challenge.ID, &challenge.Title, &challenge.Description,
			&challenge.Difficulty, &challenge.Language, &challenge.CreatedAt, &challenge.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan challenge"})
			return
		}
		challenges = append(challenges, challenge)
	}

	c.JSON(http.StatusOK, gin.H{
		"challenges": challenges,
		"total":      len(challenges),
	})
}

func GetChallengeHandler(c *gin.Context) {
	challengeID := c.Param("id")
	id, err := strconv.Atoi(challengeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid challenge ID"})
		return
	}

	var challenge Challenge
	err = db.QueryRow(`
		SELECT id, title, description, difficulty, language, test_cases, created_at, updated_at
		FROM challenges WHERE id = ?
	`, id).Scan(&challenge.ID, &challenge.Title, &challenge.Description,
		&challenge.Difficulty, &challenge.Language, &challenge.TestCases,
		&challenge.CreatedAt, &challenge.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch challenge"})
		return
	}

	c.JSON(http.StatusOK, challenge)
}

func SubmitSolutionHandler(c *gin.Context) {
	challengeID := c.Param("id")
	id, err := strconv.Atoi(challengeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid challenge ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var req SubmitSolutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify challenge exists
	var challengeExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM challenges WHERE id = ?)", id).Scan(&challengeExists)
	if err != nil || !challengeExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
		return
	}

	// Process the submission (simplified - in real implementation, this would run in a sandbox)
	status, score, output := processSubmission(req.Code, req.Language, id)

	// Insert submission
	result, err := db.Exec(`
		INSERT INTO submissions (user_id, challenge_id, code, language, status, score, output)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, userID, id, req.Code, req.Language, status, score, output)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save submission"})
		return
	}

	submissionID, _ := result.LastInsertId()

	submission := Submission{
		ID:          int(submissionID),
		UserID:      userID.(int),
		ChallengeID: id,
		Code:        req.Code,
		Language:    req.Language,
		Status:      status,
		Score:       score,
		Output:      output,
		CreatedAt:   time.Now(),
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Solution submitted successfully",
		"submission": submission,
	})
}

func GetSubmissionsHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	rows, err := db.Query(`
		SELECT s.id, s.user_id, s.challenge_id, s.code, s.language, s.status, s.score, s.output, s.created_at,
		       c.title as challenge_title
		FROM submissions s
		JOIN challenges c ON s.challenge_id = c.id
		WHERE s.user_id = ?
		ORDER BY s.created_at DESC
		LIMIT ? OFFSET ?
	`, userID, limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submissions"})
		return
	}
	defer rows.Close()

	var submissions []map[string]interface{}
	for rows.Next() {
		var submission Submission
		var challengeTitle string
		err := rows.Scan(&submission.ID, &submission.UserID, &submission.ChallengeID,
			&submission.Code, &submission.Language, &submission.Status, &submission.Score,
			&submission.Output, &submission.CreatedAt, &challengeTitle)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan submission"})
			return
		}

		submissionData := map[string]interface{}{
			"id":              submission.ID,
			"user_id":         submission.UserID,
			"challenge_id":    submission.ChallengeID,
			"challenge_title": challengeTitle,
			"code":            submission.Code,
			"language":        submission.Language,
			"status":          submission.Status,
			"score":           submission.Score,
			"output":          submission.Output,
			"created_at":      submission.CreatedAt,
		}
		submissions = append(submissions, submissionData)
	}

	c.JSON(http.StatusOK, gin.H{
		"submissions": submissions,
		"total":       len(submissions),
	})
}

func GetSubmissionHandler(c *gin.Context) {
	submissionID := c.Param("id")
	id, err := strconv.Atoi(submissionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var submission Submission
	err = db.QueryRow(`
		SELECT id, user_id, challenge_id, code, language, status, score, output, created_at
		FROM submissions WHERE id = ? AND user_id = ?
	`, id, userID).Scan(&submission.ID, &submission.UserID, &submission.ChallengeID,
		&submission.Code, &submission.Language, &submission.Status, &submission.Score,
		&submission.Output, &submission.CreatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submission"})
		return
	}

	c.JSON(http.StatusOK, submission)
}

func GetLeaderboardHandler(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")

	rows, err := db.Query(`
		SELECT u.username, 
		       COALESCE(SUM(s.score), 0) as total_score,
		       COUNT(s.id) as submissions,
		       COALESCE(MAX(s.created_at), u.created_at) as last_activity
		FROM users u
		LEFT JOIN submissions s ON u.id = s.user_id
		GROUP BY u.id, u.username, u.created_at
		ORDER BY total_score DESC, submissions DESC
		LIMIT ?
	`, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}
	defer rows.Close()

	var leaderboard []LeaderboardEntry
	for rows.Next() {
		var entry LeaderboardEntry
		var lastActivity time.Time
		err := rows.Scan(&entry.Username, &entry.TotalScore, &entry.Submissions, &lastActivity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan leaderboard entry"})
			return
		}
		entry.LastActivity = lastActivity.Format("2006-01-02 15:04:05")
		leaderboard = append(leaderboard, entry)
	}

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": leaderboard,
		"total":       len(leaderboard),
	})
}

// Simplified submission processing (in production, this would use Docker containers)
func processSubmission(code, language string, challengeID int) (status string, score int, output string) {
	// Get test cases for the challenge
	var testCasesJSON string
	err := db.QueryRow("SELECT test_cases FROM challenges WHERE id = ?", challengeID).Scan(&testCasesJSON)
	if err != nil {
		return "failed", 0, "Error: Could not load test cases"
	}

	// Parse test cases
	var testCases []map[string]string
	if err := json.Unmarshal([]byte(testCasesJSON), &testCases); err != nil {
		return "failed", 0, "Error: Invalid test cases format"
	}

	// Simulate code execution (in production, this would run in a secure sandbox)
	// For demo purposes, we'll just return a success status
	passedTests := len(testCases) // Assume all tests pass for demo
	totalTests := len(testCases)

	if passedTests == totalTests {
		status = "passed"
		score = 100
		output = "All tests passed! Great job!"
	} else {
		status = "failed"
		score = (passedTests * 100) / totalTests
		output = "Some tests failed. Keep trying!"
	}

	return status, score, output
}
