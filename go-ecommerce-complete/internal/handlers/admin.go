package handlers

import (
	"database/sql"
	"ecommerce/internal/database"
	"ecommerce/internal/middleware"
	"ecommerce/internal/models"
	"encoding/json"
	"net/http"
	"time"
)

func SellerGetProducts(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	userID := session["user_id"].(int)
	role := session["role"].(string)
	var query string
	var args []interface{}
	if role == "seller" {
		query = "SELECT id, name, description, price, stock, category_id, seller_id, brand, image_url, status, created_at, updated_at FROM products WHERE seller_id = ? ORDER BY created_at DESC"
		args = append(args, userID)
	} else if role == "admin" {
		query = "SELECT id, name, description, price, stock, category_id, seller_id, brand, image_url, status, created_at, updated_at FROM products ORDER BY created_at DESC"
	} else {
		middleware.JSONError(w, "权限不足", http.StatusForbidden)
		return
	}
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		middleware.JSONError(w, "数据库错误", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var p models.Product
		rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CategoryID, &p.SellerID, &p.Brand, &p.ImageURL, &p.Status, &p.CreatedAt, &p.UpdatedAt)
		products = append(products, p)
	}
	middleware.JSON(w, map[string]interface{}{
		"success":  true,
		"products": products,
	})
}

func SellerAddProduct(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	sellerID := session["user_id"].(int)
	var product models.Product
	json.NewDecoder(r.Body).Decode(&product)
	if product.ImageURL == "" {
		product.ImageURL = "https://via.placeholder.com/300x300?text=Product"
	}
	_, err := database.DB.Exec(`
		INSERT INTO products (name, description, price, stock, category_id, seller_id, brand, image_url, status, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'pending', ?)
	`, product.Name, product.Description, product.Price, product.Stock, product.CategoryID, sellerID, product.Brand, product.ImageURL, time.Now())
	if err != nil {
		middleware.JSONError(w, "添加失败", http.StatusInternalServerError)
		return
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "商品已提交审核",
	})
}

func SellerUpdateProduct(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	sellerID := session["user_id"].(int)
	productID := r.URL.Query().Get("id")
	var product models.Product
	json.NewDecoder(r.Body).Decode(&product)
	_, err := database.DB.Exec(`
		UPDATE products
		SET name = ?, description = ?, price = ?, stock = ?, category_id = ?, brand = ?, image_url = ?, updated_at = ?
		WHERE id = ? AND seller_id = ?
	`, product.Name, product.Description, product.Price, product.Stock, product.CategoryID, product.Brand, product.ImageURL, time.Now(), productID, sellerID)
	if err != nil {
		middleware.JSONError(w, "更新失败", http.StatusInternalServerError)
		return
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "商品更新成功",
	})
}

func SellerDeleteProduct(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	sellerID := session["user_id"].(int)
	productID := r.URL.Query().Get("id")
	_, err := database.DB.Exec("DELETE FROM products WHERE id = ? AND seller_id = ?", productID, sellerID)
	if err != nil {
		middleware.JSONError(w, "删除失败", http.StatusInternalServerError)
		return
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "商品已删除",
	})
}

func SellerGetOrders(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	sellerID := session["user_id"].(int)
	rows, err := database.DB.Query(`
		SELECT DISTINCT o.id, o.order_no, o.user_id, o.total_amount, o.status, o.payment_method,
			   o.shipping_address, o.receiver_name, o.receiver_phone, o.shipping_method,
			   o.tracking_number, o.created_at, o.paid_at, o.shipped_at, o.completed_at
		FROM orders o
		JOIN order_items oi ON o.id = oi.order_id
		JOIN products p ON oi.product_id = p.id
		WHERE p.seller_id = ?
		ORDER BY o.created_at DESC
	`, sellerID)
	if err != nil {
		middleware.JSONError(w, "数据库错误", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var orders []models.Order
	for rows.Next() {
		var o models.Order
		rows.Scan(&o.ID, &o.OrderNo, &o.UserID, &o.TotalAmount, &o.Status, &o.PaymentMethod,
			&o.ShippingAddress, &o.ReceiverName, &o.ReceiverPhone, &o.ShippingMethod,
			&o.TrackingNumber, &o.CreatedAt, &o.PaidAt, &o.ShippedAt, &o.CompletedAt)
		itemRows, _ := database.DB.Query(`
			SELECT oi.id, oi.order_id, oi.product_id, oi.product_name, oi.product_price, oi.quantity, oi.subtotal
			FROM order_items oi
			JOIN products p ON oi.product_id = p.id
			WHERE oi.order_id = ? AND p.seller_id = ?
		`, o.ID, sellerID)
		for itemRows.Next() {
			var item models.OrderItem
			itemRows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductName, &item.ProductPrice, &item.Quantity, &item.Subtotal)
			o.Items = append(o.Items, item)
		}
		itemRows.Close()
		orders = append(orders, o)
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"orders":  orders,
	})
}

func AdminApproveProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	_, err := database.DB.Exec("UPDATE products SET status = 'approved' WHERE id = ?", productID)
	if err != nil {
		middleware.JSONError(w, "操作失败", http.StatusInternalServerError)
		return
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "商品已审核通过",
	})
}

func AdminRejectProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	_, err := database.DB.Exec("UPDATE products SET status = 'rejected' WHERE id = ?", productID)
	if err != nil {
		middleware.JSONError(w, "操作失败", http.StatusInternalServerError)
		return
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "商品已拒绝",
	})
}

func AdminGetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, username, email, phone, role, status, created_at FROM users")
	if err != nil {
		middleware.JSONError(w, "数据库错误", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		var phone sql.NullString
		rows.Scan(&u.ID, &u.Username, &u.Email, &phone, &u.Role, &u.Status, &u.CreatedAt)
		if phone.Valid {
			u.Phone = phone.String
		}
		users = append(users, u)
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"users":   users,
	})
}

func AdminGetStats(w http.ResponseWriter, r *http.Request) {
	var stats struct {
		TotalUsers    int     `json:"total_users"`
		TotalProducts int     `json:"total_products"`
		TotalOrders   int     `json:"total_orders"`
		TotalSales    float64 `json:"total_sales"`
	}
	database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	database.DB.QueryRow("SELECT COUNT(*) FROM products WHERE status = 'approved'").Scan(&stats.TotalProducts)
	database.DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&stats.TotalOrders)
	database.DB.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE status != 'cancelled'").Scan(&stats.TotalSales)
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"stats":   stats,
	})
}
