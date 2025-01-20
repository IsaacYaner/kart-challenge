package service

import (
	"encoding/json"
	"os"
	"sync"

	"backend-challenge/internal/models"
)

type ProductService struct {
	products []models.Product
	mu       sync.RWMutex
}

type productData struct {
	Products []models.Product `json:"products"`
}

func NewProductService() (*ProductService, error) {
	ps := &ProductService{}
	if err := ps.loadProducts(); err != nil {
		return nil, err
	}
	return ps, nil
}

func (ps *ProductService) loadProducts() error {
	file, err := os.ReadFile("internal/data/products.json")
	if err != nil {
		return err
	}

	var data productData
	if err := json.Unmarshal(file, &data); err != nil {
		return err
	}

	ps.mu.Lock()
	ps.products = data.Products
	ps.mu.Unlock()

	return nil
}

func (ps *ProductService) GetProducts() []models.Product {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.products
}

func (ps *ProductService) GetProductByID(id string) (*models.Product, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, product := range ps.products {
		if product.ID == id {
			return &product, true
		}
	}
	return nil, false
}
