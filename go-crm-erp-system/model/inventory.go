package model

import "time"

type Inventory struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,gte=0"`
	Warehouse string    `json:"warehouse"`
	UpdatedAt time.Time `json:"updated_at"`
}
