package controllers

import (
	"codelearn-backend/db"
	"codelearn-backend/models"
	"codelearn-backend/utils"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         models.User `json:"user"`
}

func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingID int
	err := db.DB.QueryRow("SELECT id FROM users WHERE username = ? OR email = ?", req.Username, req.Email).Scan(&existingID)
	if err != sql.ErrNoRows {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	hashedPassword, err := utils.HashingPassword([]byte(req.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	result, err := db.DB.Exec(`
		INSERT INTO users (username, email, password)
		VALUES (?, ?, ?)
	`, req.Username, req.Email, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	userID, _ := result.LastInsertId()

	user := models.User{
		ID:       int(userID),
		Username: req.Username,
		Email:    req.Email,
	}

	token, refreshToken, err := utils.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	})
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	var hashedPassword string
	err := db.DB.QueryRow(`
		SELECT id, username, email, password, created_at, updated_at
		FROM users WHERE username = ?
	`, req.Username).Scan(&user.ID, &user.Username, &user.Email, &hashedPassword, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := utils.ComparePassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, refreshToken, err := utils.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	})
}

func RefreshTokenHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Refresh token not implemented yet"})
}

func GetProfileHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var user models.User
	err := db.DB.QueryRow(`
		SELECT id, username, email, created_at, updated_at
		FROM users WHERE id = ?
	`, userID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateProfileHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.DB.Exec(`
		UPDATE users SET email = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, req.Email, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func CLIAuthHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found in context"})
		return
	}

	cliClaims := &Claims{
		UserID:   userID.(int),
		Username: username.(string),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // 30 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	cliToken := jwt.NewWithClaims(jwt.SigningMethodHS256, cliClaims)
	cliTokenString, err := cliToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CLI token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cli_token":  cliTokenString,
		"user_id":    userID,
		"username":   username,
		"expires_at": time.Now().Add(30 * 24 * time.Hour).Unix(),
	})
}
