package service

import (
	"fmt"
	"sync"
	"time"

	"backend-challenge/internal/models"
)

type OrderService struct {
	productService *ProductService
	couponService  *CouponService
	mu             sync.RWMutex
	orders         map[string]models.Order
}

func NewOrderService(ps *ProductService, cs *CouponService) *OrderService {
	return &OrderService{
		productService: ps,
		couponService:  cs,
		orders:         make(map[string]models.Order),
	}
}

func (s *OrderService) CreateOrder(req models.OrderRequest) (*models.Order, error) {
	// Validate coupon code first
	if req.CouponCode != "" && !s.couponService.IsValidCoupon(req.CouponCode) {
		return nil, fmt.Errorf("invalid coupon code: %s", req.CouponCode)
	}
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
