package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"api-simulator/handlers"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "ecoround-api-simulator"})
	})

	matchHandler := handlers.NewMatchHandler()
	matchHandler.RegisterRoutes(r)

	log.Println("EcoRound API Simulator starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
