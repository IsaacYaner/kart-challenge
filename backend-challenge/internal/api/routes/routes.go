package routes

import (
	"backend-challenge/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Product routes
		api.GET("/product", handlers.ListProducts)
		api.GET("/product/:productId", handlers.GetProduct)
	}
} 