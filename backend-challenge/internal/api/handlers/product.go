package handlers

import (
	"net/http"
	"strconv"

	"backend-challenge/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler() (*ProductHandler, error) {
	ps, err := service.NewProductService()
	if err != nil {
		return nil, err
	}
	return &ProductHandler{productService: ps}, nil
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	products := h.productService.GetProducts()
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("productId")

	// Validate ID is a proper number
	_, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID supplied"})
		return
	}

	product, found := h.productService.GetProductByID(productID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetProductService() *service.ProductService {
	return h.productService
}
