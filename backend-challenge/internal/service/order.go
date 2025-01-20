package service

import (
	"fmt"
	"sync"
	"time"

	"backend-challenge/internal/models"
)

type OrderService struct {
	productService *ProductService
	mu             sync.RWMutex
	orders         map[string]models.Order
}

func NewOrderService(ps *ProductService) *OrderService {
	return &OrderService{
		productService: ps,
		orders:         make(map[string]models.Order),
	}
}

func (s *OrderService) CreateOrder(req models.OrderRequest) (*models.Order, error) {
	// Validate all products exist
	var products []models.Product
	for _, item := range req.Items {
		product, exists := s.productService.GetProductByID(item.ProductID)
		if !exists {
			return nil, fmt.Errorf("product not found: %s", item.ProductID)
		}
		products = append(products, *product)
	}

	// Create order
	order := models.Order{
		ID:       generateOrderID(),
		Items:    req.Items,
		Products: products,
	}

	// Store order
	s.mu.Lock()
	s.orders[order.ID] = order
	s.mu.Unlock()

	return &order, nil
}

func generateOrderID() string {
	// Simple implementation - in production, use UUID or other robust ID generation
	return fmt.Sprintf("ORDER-%d", time.Now().UnixNano())
}
