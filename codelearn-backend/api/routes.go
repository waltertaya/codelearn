package api

import (
	"codelearn-backend/controllers"
	"codelearn-backend/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.CORSMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "CodeLearn Backend API is running",
		})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.RegisterHandler)
			auth.POST("/login", controllers.LoginHandler)
			auth.POST("/refresh", controllers.RefreshTokenHandler)
		}

		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware())
		{
			protected.GET("/profile", controllers.GetProfileHandler)
			protected.PUT("/profile", controllers.UpdateProfileHandler)

			protected.GET("/challenges", controllers.GetChallengesHandler)
			protected.GET("/challenges/:id", controllers.GetChallengeHandler)
			protected.POST("/challenges/:id/submit", controllers.SubmitSolutionHandler)

			protected.GET("/submissions", controllers.GetSubmissionsHandler)
			protected.GET("/submissions/:id", controllers.GetSubmissionHandler)

			protected.GET("/leaderboard", controllers.GetLeaderboardHandler)

			protected.POST("/cli/auth", controllers.CLIAuthHandler)
		}
	}

	return r
}
