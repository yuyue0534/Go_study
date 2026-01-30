package models

import (
	"ecommerce-platform/internal/config"
	"time"
)

type CartItem struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Selected  int       `json:"selected"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 商品信息
	Product *Product `json:"product,omitempty"`
}

// AddToCart 添加商品到购物车
func AddToCart(userID, productID int64, quantity int) error {
	// 检查是否已存在
	var existID int64
	err := config.DB.QueryRow(`SELECT id FROM cart_items WHERE user_id = ? AND product_id = ?`, userID, productID).Scan(&existID)
	
	if err == nil {
		// 已存在，更新数量
		_, err = config.DB.Exec(`
			UPDATE cart_items SET quantity = quantity + ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`, quantity, existID)
		return err
	}
	
	// 不存在，新增
	_, err = config.DB.Exec(`
		INSERT INTO cart_items (user_id, product_id, quantity, selected)
		VALUES (?, ?, ?, 1)
	`, userID, productID, quantity)
	return err
}

// GetCartItems 获取用户购物车
func GetCartItems(userID int64) ([]*CartItem, error) {
	rows, err := config.DB.Query(`
		SELECT ci.id, ci.user_id, ci.product_id, ci.quantity, ci.selected, ci.created_at, ci.updated_at,
			p.name, p.price, p.original_price, p.stock, p.images, s.shop_name
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		LEFT JOIN sellers s ON p.seller_id = s.id
		WHERE ci.user_id = ?
		ORDER BY ci.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*CartItem
	for rows.Next() {
		item := &CartItem{Product: &Product{}}
		err := rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.Selected,
			&item.CreatedAt, &item.UpdatedAt, &item.Product.Name, &item.Product.Price,
			&item.Product.OriginalPrice, &item.Product.Stock, &item.Product.Images, &item.Product.SellerName)
		if err != nil {
			return nil, err
		}
		item.Product.ID = item.ProductID
		items = append(items, item)
	}
	return items, nil
}

// UpdateCartItemQuantity 更新购物车商品数量
func UpdateCartItemQuantity(id int64, quantity int) error {
	_, err := config.DB.Exec(`
		UPDATE cart_items SET quantity = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
	`, quantity, id)
	return err
}

// UpdateCartItemSelected 更新购物车商品选中状态
func UpdateCartItemSelected(id int64, selected int) error {
	_, err := config.DB.Exec(`
		UPDATE cart_items SET selected = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
	`, selected, id)
	return err
}

// UpdateCartAllSelected 更新用户购物车全部选中状态
func UpdateCartAllSelected(userID int64, selected int) error {
	_, err := config.DB.Exec(`
		UPDATE cart_items SET selected = ?, updated_at = CURRENT_TIMESTAMP WHERE user_id = ?
	`, selected, userID)
	return err
}

// DeleteCartItem 删除购物车商品
func DeleteCartItem(id int64) error {
	_, err := config.DB.Exec(`DELETE FROM cart_items WHERE id = ?`, id)
	return err
}

// DeleteCartItemByUserAndProduct 删除用户购物车中的指定商品
func DeleteCartItemByUserAndProduct(userID, productID int64) error {
	_, err := config.DB.Exec(`DELETE FROM cart_items WHERE user_id = ? AND product_id = ?`, userID, productID)
	return err
}

// ClearCart 清空用户购物车
func ClearCart(userID int64) error {
	_, err := config.DB.Exec(`DELETE FROM cart_items WHERE user_id = ?`, userID)
	return err
}

// ClearSelectedCart 清空用户选中的购物车商品
func ClearSelectedCart(userID int64) error {
	_, err := config.DB.Exec(`DELETE FROM cart_items WHERE user_id = ? AND selected = 1`, userID)
	return err
}

// GetSelectedCartItems 获取用户选中的购物车商品
func GetSelectedCartItems(userID int64) ([]*CartItem, error) {
	rows, err := config.DB.Query(`
		SELECT ci.id, ci.user_id, ci.product_id, ci.quantity, ci.selected, ci.created_at, ci.updated_at,
			p.name, p.price, p.original_price, p.stock, p.images, p.seller_id, s.shop_name
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		LEFT JOIN sellers s ON p.seller_id = s.id
		WHERE ci.user_id = ? AND ci.selected = 1
		ORDER BY ci.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*CartItem
	for rows.Next() {
		item := &CartItem{Product: &Product{}}
		err := rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.Selected,
			&item.CreatedAt, &item.UpdatedAt, &item.Product.Name, &item.Product.Price,
			&item.Product.OriginalPrice, &item.Product.Stock, &item.Product.Images,
			&item.Product.SellerID, &item.Product.SellerName)
		if err != nil {
			return nil, err
		}
		item.Product.ID = item.ProductID
		items = append(items, item)
	}
	return items, nil
}

// GetCartCount 获取用户购物车商品数量
func GetCartCount(userID int64) (int, error) {
	var count int
	err := config.DB.QueryRow(`SELECT COALESCE(SUM(quantity), 0) FROM cart_items WHERE user_id = ?`, userID).Scan(&count)
	return count, err
}
