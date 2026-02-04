package model

import "time"

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" binding:"required"`
	SKU         string    `json:"sku" binding:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"required,gt=0"`
	Cost        float64   `json:"cost"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
