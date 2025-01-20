package main

import (
	"log"
	"net/http"

	"backend-challenge/internal/api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, api_key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// API key middleware
	r.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		apiKey := c.GetHeader("api_key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
			c.Abort()
			return
		}

		if apiKey != "apitest" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	})

	// Setup routes
	routes.SetupRoutes(r)
	
	log.Printf("Server starting on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
} 