#!/bin/bash
set -e

echo "开始创建完整的Go电商平台项目..."

# 复制前端文件
echo "1. 复制前端文件..."
cp /mnt/user-data/outputs/ecommerce-platform/static/css/style.css web/static/css/
cp /mnt/user-data/outputs/ecommerce-platform/static/js/common.js web/static/js/
cp /mnt/user-data/outputs/ecommerce-platform/templates/*.html web/templates/

echo "2. 创建 go.mod..."
cat > go.mod << 'EOF'
module ecommerce

go 1.19

require github.com/mattn/go-sqlite3 v1.14.18
EOF

echo "3. 创建 internal/models/models.go..."
cat > internal/models/models.go << 'EOF'
package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Address struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	ReceiverName  string    `json:"receiver_name"`
	Phone         string    `json:"phone"`
	Province      string    `json:"province"`
	City          string    `json:"city"`
	District      string    `json:"district"`
	DetailAddress string    `json:"detail_address"`
	IsDefault     bool      `json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
}

type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    *int      `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CategoryID  int       `json:"category_id"`
	SellerID    int       `json:"seller_id"`
	Brand       string    `json:"brand"`
	ImageURL    string    `json:"image_url"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CartItem struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name,omitempty"`
	Price     float64   `json:"price,omitempty"`
	Stock     int       `json:"stock,omitempty"`
	ImageURL  string    `json:"image_url,omitempty"`
}

type Order struct {
	ID              int          `json:"id"`
	OrderNo         string       `json:"order_no"`
	UserID          int          `json:"user_id"`
	TotalAmount     float64      `json:"total_amount"`
	Status          string       `json:"status"`
	PaymentMethod   string       `json:"payment_method"`
	ShippingAddress string       `json:"shipping_address"`
	ReceiverName    string       `json:"receiver_name"`
	ReceiverPhone   string       `json:"receiver_phone"`
	ShippingMethod  string       `json:"shipping_method"`
	TrackingNumber  string       `json:"tracking_number"`
	CreatedAt       time.Time    `json:"created_at"`
	PaidAt          *time.Time   `json:"paid_at"`
	ShippedAt       *time.Time   `json:"shipped_at"`
	CompletedAt     *time.Time   `json:"completed_at"`
	Items           []OrderItem  `json:"items,omitempty"`
}

type OrderItem struct {
	ID           int     `json:"id"`
	OrderID      int     `json:"order_id"`
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
	Quantity     int     `json:"quantity"`
	Subtotal     float64 `json:"subtotal"`
}

type Review struct {
	ID        int           `json:"id"`
	UserID    int           `json:"user_id"`
	ProductID int           `json:"product_id"`
	OrderID   int           `json:"order_id"`
	Rating    int           `json:"rating"`
	Comment   string        `json:"comment"`
	Status    string        `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	Username  string        `json:"username,omitempty"`
	Replies   []ReviewReply `json:"replies,omitempty"`
}

type ReviewReply struct {
	ID        int       `json:"id"`
	ReviewID  int       `json:"review_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username,omitempty"`
}
EOF

echo "所有文件创建完成！"
