package handlers

import (
	"net/http"

	"backend-challenge/internal/models"
	"backend-challenge/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(productService *service.ProductService) *OrderHandler {
	couponService := service.NewCouponService()
	return &OrderHandler{
		orderService: service.NewOrderService(productService, couponService),
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate items array is not empty
	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order must contain at least one item"})
		return
	}

	// Create order
	order, err := h.orderService.CreateOrder(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
