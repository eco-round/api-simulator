package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"api-simulator/db"
	"api-simulator/handlers"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	db.Init()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "ecoround-api-simulator"})
	})

	// Source endpoints (what CRE calls)
	matchHandler := handlers.NewMatchHandler()
	matchHandler.RegisterRoutes(r)

	// Admin endpoints (for managing matches)
	adminHandler := handlers.NewAdminHandler()
	adminHandler.RegisterRoutes(r)

	log.Println("EcoRound API Simulator starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
