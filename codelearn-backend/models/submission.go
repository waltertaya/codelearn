package models

import "time"

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
