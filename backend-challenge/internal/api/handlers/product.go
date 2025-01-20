package handlers

import (
	"net/http"
	"strconv"

	"backend-challenge/internal/models"
	"github.com/gin-gonic/gin"
)

// Sample data, which might come from a database later.
var products = []models.Product{
	{
		ID:       "1",
		Name:     "Chicken burger",
		Price:    13.3,
		Category: "Burger",
	},
	// TODO: Add more sample products later
}

func ListProducts(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

func GetProduct(c *gin.Context) {
	productID := c.Param("productId")
	
	// Convert string ID to int64 for validation
	_, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Find product by ID
	for _, product := range products {
		if product.ID == productID {
			c.JSON(http.StatusOK, product)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
} 