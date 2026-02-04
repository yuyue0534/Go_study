package model

import "time"

type Order struct {
	ID          int64     `json:"id"`
	OrderNo     string    `json:"order_no"`
	CustomerID  int64     `json:"customer_id" binding:"required"`
	ProductID   int64     `json:"product_id" binding:"required"`
	Quantity    int       `json:"quantity" binding:"required,gt=0"`
	UnitPrice   float64   `json:"unit_price" binding:"required,gt=0"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	UserID      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
