package models

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type OrderRequest struct {
	CouponCode string      `json:"couponCode,omitempty"`
	Items      []OrderItem `json:"items"`
}

type Order struct {
	ID       string    `json:"id"`
	Items    []OrderItem `json:"items"`
	Products []Product `json:"products"`
} 