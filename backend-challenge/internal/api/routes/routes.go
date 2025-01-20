package routes

import (
	"log"

	"backend-challenge/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	productHandler, err := handlers.NewProductHandler()
	if err != nil {
		log.Fatalf("Failed to initialize product handler: %v", err)
	}

	orderHandler := handlers.NewOrderHandler(productHandler.GetProductService())

	api := r.Group("/api")
	{
		// Product routes
		api.GET("/product", productHandler.ListProducts)
		api.GET("/product/:productId", productHandler.GetProduct)

		// Order routes
		api.POST("/order", orderHandler.CreateOrder)
	}
}
