package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	if err := InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Create Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // In production, specify your frontend domain
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"message": "CodeLearn Backend API is running",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", RegisterHandler)
			auth.POST("/login", LoginHandler)
			auth.POST("/refresh", RefreshTokenHandler)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(AuthMiddleware())
		{
			// User profile
			protected.GET("/profile", GetProfileHandler)
			protected.PUT("/profile", UpdateProfileHandler)

			// Challenges
			protected.GET("/challenges", GetChallengesHandler)
			protected.GET("/challenges/:id", GetChallengeHandler)
			protected.POST("/challenges/:id/submit", SubmitSolutionHandler)

			// Submissions
			protected.GET("/submissions", GetSubmissionsHandler)
			protected.GET("/submissions/:id", GetSubmissionHandler)

			// Leaderboard
			protected.GET("/leaderboard", GetLeaderboardHandler)

			// CLI authentication
			protected.POST("/cli/auth", CLIAuthHandler)
		}
	}

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting CodeLearn Backend on port %s", port)
	log.Fatal(r.Run(":" + port))
}
