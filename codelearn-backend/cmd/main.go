package main

import (
	"codelearn-backend/api"
	"codelearn-backend/db"
	"log"
	"os"
)

func init() {
	if err := db.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
}

func main() {

	r := api.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting CodeLearn Backend on port %s", port)
	log.Fatal(r.Run(":" + port))
}
