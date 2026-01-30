package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"ecommerce/internal/database"
	"ecommerce/internal/middleware"
	"ecommerce/internal/models"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.JSONError(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" || req.Email == "" {
		middleware.JSONError(w, "请填写完整信息", http.StatusBadRequest)
		return
	}
	var exists int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&exists)
	if err != nil {
		middleware.JSONError(w, "数据库错误", http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		middleware.JSONError(w, "用户名已存在", http.StatusBadRequest)
		return
	}
	err = database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.Email).Scan(&exists)
	if err != nil {
		middleware.JSONError(w, "数据库错误", http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		middleware.JSONError(w, "邮箱已被注册", http.StatusBadRequest)
		return
	}
	hash := sha256.Sum256([]byte(req.Password))
	hashedPassword := fmt.Sprintf("%x", hash)
	_, err = database.DB.Exec(
		"INSERT INTO users (username, password, email, phone, role) VALUES (?, ?, ?, ?, 'customer')",
		req.Username, hashedPassword, req.Email, req.Phone,
	)
	if err != nil {
		middleware.JSONError(w, "注册失败", http.StatusInternalServerError)
		return
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "注册成功",
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.JSONError(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	hash := sha256.Sum256([]byte(req.Password))
	hashedPassword := fmt.Sprintf("%x", hash)
	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, username, email, role, status FROM users WHERE username = ? AND password = ?",
		req.Username, hashedPassword,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Status)
	if err == sql.ErrNoRows {
		middleware.JSONError(w, "用户名或密码错误", http.StatusUnauthorized)
		return
	}
	if err != nil {
		middleware.JSONError(w, "数据库错误", http.StatusInternalServerError)
		return
	}
	if user.Status != "active" {
		middleware.JSONError(w, "账户已被冻结", http.StatusForbidden)
		return
	}
	sessionID := fmt.Sprintf("%d-%d", user.ID, time.Now().Unix())
	middleware.Sessions[sessionID] = map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	})
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "登录成功",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		delete(middleware.Sessions, cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "已退出登录",
	})
}

func CurrentUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		middleware.JSONError(w, "未登录", http.StatusUnauthorized)
		return
	}
	session, ok := middleware.Sessions[cookie.Value]
	if !ok {
		middleware.JSONError(w, "会话已过期", http.StatusUnauthorized)
		return
	}
	userID := session["user_id"].(int)
	var user models.User
	err = database.DB.QueryRow(
		"SELECT id, username, email, phone, role FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.Role)
	if err != nil {
		middleware.JSONError(w, "用户不存在", http.StatusNotFound)
		return
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"user":    user,
	})
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category_id")
	search := r.URL.Query().Get("search")
	sort := r.URL.Query().Get("sort")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage := 12
	offset := (page - 1) * perPage
	query := "SELECT id, name, description, price, stock, category_id, seller_id, brand, image_url, status, created_at, updated_at FROM products WHERE status = 'approved'"
	var args []interface{}
	if categoryID != "" {
		query += " AND category_id = ?"
		args = append(args, categoryID)
	}
	if search != "" {
		query += " AND (name LIKE ? OR description LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}
	switch sort {
	case "price_asc":
		query += " ORDER BY price ASC"
	case "price_desc":
		query += " ORDER BY price DESC"
	default:
		query += " ORDER BY created_at DESC"
	}
	countQuery := strings.Replace(query, "SELECT id, name, description, price, stock, category_id, seller_id, brand, image_url, status, created_at, updated_at", "SELECT COUNT(*)", 1)
	var total int
	database.DB.QueryRow(countQuery, args...).Scan(&total)
	query += " LIMIT ? OFFSET ?"
	args = append(args, perPage, offset)
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		middleware.JSONError(w, "数据库错误", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CategoryID, &p.SellerID, &p.Brand, &p.ImageURL, &p.Status, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			continue
		}
		products = append(products, p)
	}
	middleware.JSON(w, map[string]interface{}{
		"success":  true,
		"products": products,
		"total":    total,
		"page":     page,
		"pages":    (total + perPage - 1) / perPage,
	})
}

func GetProductDetail(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	var product models.Product
	err := database.DB.QueryRow(
		"SELECT id, name, description, price, stock, category_id, seller_id, brand, image_url, status, created_at, updated_at FROM products WHERE id = ?",
		productID,
	).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CategoryID, &product.SellerID, &product.Brand, &product.ImageURL, &product.Status, &product.CreatedAt, &product.UpdatedAt)
	if err == sql.ErrNoRows {
		middleware.JSONError(w, "商品不存在", http.StatusNotFound)
		return
	}
	if err != nil {
		middleware.JSONError(w, "数据库错误", http.StatusInternalServerError)
		return
	}
	rows, err := database.DB.Query(`
		SELECT r.id, r.user_id, r.product_id, r.order_id, r.rating, r.comment, r.status, r.created_at, u.username
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.product_id = ? AND r.status = 'approved'
		ORDER BY r.created_at DESC
	`, productID)
	var reviews []models.Review
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var r models.Review
			rows.Scan(&r.ID, &r.UserID, &r.ProductID, &r.OrderID, &r.Rating, &r.Comment, &r.Status, &r.CreatedAt, &r.Username)
			reviews = append(reviews, r)
		}
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"product": product,
		"reviews": reviews,
	})
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, description, parent_id, created_at FROM categories ORDER BY id")
	if err != nil {
		middleware.JSONError(w, "数据库错误", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var categories []models.Category
	for rows.Next() {
		var c models.Category
		rows.Scan(&c.ID, &c.Name, &c.Description, &c.ParentID, &c.CreatedAt)
		categories = append(categories, c)
	}
	middleware.JSON(w, map[string]interface{}{
		"success":    true,
		"categories": categories,
	})
}

func GetCart(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	userID := session["user_id"].(int)
	rows, err := database.DB.Query(`
		SELECT c.id, c.user_id, c.product_id, c.quantity, c.created_at,
			   p.name, p.price, p.stock, p.image_url
		FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = ?
	`, userID)
	if err != nil {
		middleware.JSONError(w, "数据库错误", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.CreatedAt,
			&item.Name, &item.Price, &item.Stock, &item.ImageURL)
		items = append(items, item)
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"items":   items,
	})
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	userID := session["user_id"].(int)
	var req struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Quantity < 1 {
		req.Quantity = 1
	}
	var stock int
	err := database.DB.QueryRow("SELECT stock FROM products WHERE id = ?", req.ProductID).Scan(&stock)
	if err == sql.ErrNoRows {
		middleware.JSONError(w, "商品不存在", http.StatusNotFound)
		return
	}
	if stock < req.Quantity {
		middleware.JSONError(w, "库存不足", http.StatusBadRequest)
		return
	}
	var existingID, existingQty int
	err = database.DB.QueryRow(
		"SELECT id, quantity FROM cart_items WHERE user_id = ? AND product_id = ?",
		userID, req.ProductID,
	).Scan(&existingID, &existingQty)
	if err == sql.ErrNoRows {
		_, err = database.DB.Exec(
			"INSERT INTO cart_items (user_id, product_id, quantity) VALUES (?, ?, ?)",
			userID, req.ProductID, req.Quantity,
		)
	} else {
		newQty := existingQty + req.Quantity
		if newQty > stock {
			middleware.JSONError(w, "库存不足", http.StatusBadRequest)
			return
		}
		_, err = database.DB.Exec(
			"UPDATE cart_items SET quantity = ? WHERE id = ?",
			newQty, existingID,
		)
	}
	if err != nil {
		middleware.JSONError(w, "操作失败", http.StatusInternalServerError)
		return
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "已添加到购物车",
	})
}

func UpdateCart(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	userID := session["user_id"].(int)
	var req struct {
		CartItemID int `json:"cart_item_id"`
		Quantity   int `json:"quantity"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Quantity < 1 {
		middleware.JSONError(w, "数量必须大于0", http.StatusBadRequest)
		return
	}
	var stock int
	err := database.DB.QueryRow(`
		SELECT p.stock FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.id = ? AND c.user_id = ?
	`, req.CartItemID, userID).Scan(&stock)
	if err == sql.ErrNoRows {
		middleware.JSONError(w, "购物车项不存在", http.StatusNotFound)
		return
	}
	if stock < req.Quantity {
		middleware.JSONError(w, "库存不足", http.StatusBadRequest)
		return
	}
	_, err = database.DB.Exec(
		"UPDATE cart_items SET quantity = ? WHERE id = ? AND user_id = ?",
		req.Quantity, req.CartItemID, userID,
	)
	if err != nil {
		middleware.JSONError(w, "更新失败", http.StatusInternalServerError)
		return
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "更新成功",
	})
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	userID := session["user_id"].(int)
	var req struct {
		CartItemID int `json:"cart_item_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	_, err := database.DB.Exec(
		"DELETE FROM cart_items WHERE id = ? AND user_id = ?",
		req.CartItemID, userID,
	)
	if err != nil {
		middleware.JSONError(w, "删除失败", http.StatusInternalServerError)
		return
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "已删除",
	})
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	userID := session["user_id"].(int)
	var req struct {
		CartItemIDs   []int  `json:"cart_item_ids"`
		AddressID     int    `json:"address_id"`
		PaymentMethod string `json:"payment_method"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	var addr models.Address
	err := database.DB.QueryRow(
		"SELECT receiver_name, phone, province, city, district, detail_address FROM addresses WHERE id = ? AND user_id = ?",
		req.AddressID, userID,
	).Scan(&addr.ReceiverName, &addr.Phone, &addr.Province, &addr.City, &addr.District, &addr.DetailAddress)
	if err == sql.ErrNoRows {
		middleware.JSONError(w, "地址不存在", http.StatusNotFound)
		return
	}
	placeholders := strings.Trim(strings.Repeat("?,", len(req.CartItemIDs)), ",")
	args := make([]interface{}, len(req.CartItemIDs)+1)
	args[0] = userID
	for i, id := range req.CartItemIDs {
		args[i+1] = id
	}
	query := fmt.Sprintf(`
		SELECT c.id, c.product_id, p.name, p.price, p.stock, c.quantity
		FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = ? AND c.id IN (%s)
	`, placeholders)
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		middleware.JSONError(w, "数据库错误", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	type cartItem struct {
		ID        int
		ProductID int
		Name      string
		Price     float64
		Stock     int
		Quantity  int
	}
	var items []cartItem
	var totalAmount float64
	for rows.Next() {
		var item cartItem
		rows.Scan(&item.ID, &item.ProductID, &item.Name, &item.Price, &item.Stock, &item.Quantity)
		if item.Stock < item.Quantity {
			middleware.JSONError(w, fmt.Sprintf("%s 库存不足", item.Name), http.StatusBadRequest)
			return
		}
		totalAmount += item.Price * float64(item.Quantity)
		items = append(items, item)
	}
	if len(items) == 0 {
		middleware.JSONError(w, "购物车为空", http.StatusBadRequest)
		return
	}
	orderNo := fmt.Sprintf("ORD%s%04d", time.Now().Format("20060102150405"), rand.Intn(10000))
	shippingAddress := fmt.Sprintf("%s %s %s %s", addr.Province, addr.City, addr.District, addr.DetailAddress)
	result, err := database.DB.Exec(`
		INSERT INTO orders (order_no, user_id, total_amount, status, payment_method, shipping_address, receiver_name, receiver_phone)
		VALUES (?, ?, ?, 'pending_payment', ?, ?, ?, ?)
	`, orderNo, userID, totalAmount, req.PaymentMethod, shippingAddress, addr.ReceiverName, addr.Phone)
	if err != nil {
		middleware.JSONError(w, "创建订单失败", http.StatusInternalServerError)
		return
	}
	orderID, _ := result.LastInsertId()
	for _, item := range items {
		subtotal := item.Price * float64(item.Quantity)
		database.DB.Exec(
			"INSERT INTO order_items (order_id, product_id, product_name, product_price, quantity, subtotal) VALUES (?, ?, ?, ?, ?, ?)",
			orderID, item.ProductID, item.Name, item.Price, item.Quantity, subtotal,
		)
		database.DB.Exec("UPDATE products SET stock = stock - ? WHERE id = ?", item.Quantity, item.ProductID)
	}
	database.DB.Exec(fmt.Sprintf("DELETE FROM cart_items WHERE user_id = ? AND id IN (%s)", placeholders), args...)
	middleware.JSON(w, map[string]interface{}{
		"success":  true,
		"message":  "订单创建成功",
		"order_no": orderNo,
		"order_id": orderID,
	})
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	userID := session["user_id"].(int)
	rows, err := database.DB.Query(`
		SELECT id, order_no, user_id, total_amount, status, payment_method, shipping_address, 
			   receiver_name, receiver_phone, shipping_method, tracking_number, created_at, paid_at, shipped_at, completed_at
		FROM orders
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)
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
		itemRows, _ := database.DB.Query(
			"SELECT id, order_id, product_id, product_name, product_price, quantity, subtotal FROM order_items WHERE order_id = ?",
			o.ID,
		)
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

func PayOrder(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	userID := session["user_id"].(int)
	orderID := r.URL.Query().Get("id")
	var status string
	err := database.DB.QueryRow("SELECT status FROM orders WHERE id = ? AND user_id = ?", orderID, userID).Scan(&status)
	if err == sql.ErrNoRows {
		middleware.JSONError(w, "订单不存在", http.StatusNotFound)
		return
	}
	if status != "pending_payment" {
		middleware.JSONError(w, "订单状态不正确", http.StatusBadRequest)
		return
	}
	now := time.Now()
	_, err = database.DB.Exec(
		"UPDATE orders SET status = 'pending_shipment', paid_at = ? WHERE id = ?",
		now, orderID,
	)
	if err != nil {
		middleware.JSONError(w, "支付失败", http.StatusInternalServerError)
		return
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "支付成功",
	})
}

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	userID := session["user_id"].(int)
	orderID := r.URL.Query().Get("id")
	var status string
	err := database.DB.QueryRow("SELECT status FROM orders WHERE id = ? AND user_id = ?", orderID, userID).Scan(&status)
	if err == sql.ErrNoRows {
		middleware.JSONError(w, "订单不存在", http.StatusNotFound)
		return
	}
	if status != "pending_payment" && status != "pending_shipment" {
		middleware.JSONError(w, "该订单不能取消", http.StatusBadRequest)
		return
	}
	rows, _ := database.DB.Query("SELECT product_id, quantity FROM order_items WHERE order_id = ?", orderID)
	for rows.Next() {
		var productID, quantity int
		rows.Scan(&productID, &quantity)
		database.DB.Exec("UPDATE products SET stock = stock + ? WHERE id = ?", quantity, productID)
	}
	rows.Close()
	_, err = database.DB.Exec("UPDATE orders SET status = 'cancelled' WHERE id = ?", orderID)
	if err != nil {
		middleware.JSONError(w, "取消失败", http.StatusInternalServerError)
		return
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "订单已取消",
	})
}

func GetAddresses(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	userID := session["user_id"].(int)
	rows, err := database.DB.Query(`
		SELECT id, user_id, receiver_name, phone, province, city, district, detail_address, is_default, created_at
		FROM addresses
		WHERE user_id = ?
		ORDER BY is_default DESC, created_at DESC
	`, userID)
	if err != nil {
		middleware.JSONError(w, "数据库错误", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var addresses []models.Address
	for rows.Next() {
		var a models.Address
		var isDefault int
		rows.Scan(&a.ID, &a.UserID, &a.ReceiverName, &a.Phone, &a.Province, &a.City, &a.District, &a.DetailAddress, &isDefault, &a.CreatedAt)
		a.IsDefault = isDefault == 1
		addresses = append(addresses, a)
	}
	middleware.JSON(w, map[string]interface{}{
		"success":   true,
		"addresses": addresses,
	})
}

func AddAddress(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	userID := session["user_id"].(int)
	var addr models.Address
	json.NewDecoder(r.Body).Decode(&addr)
	if addr.IsDefault {
		database.DB.Exec("UPDATE addresses SET is_default = 0 WHERE user_id = ?", userID)
	}
	isDefault := 0
	if addr.IsDefault {
		isDefault = 1
	}
	_, err := database.DB.Exec(`
		INSERT INTO addresses (user_id, receiver_name, phone, province, city, district, detail_address, is_default)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, userID, addr.ReceiverName, addr.Phone, addr.Province, addr.City, addr.District, addr.DetailAddress, isDefault)
	if err != nil {
		middleware.JSONError(w, "添加失败", http.StatusInternalServerError)
		return
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "地址添加成功",
	})
}

func DeleteAddress(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	session := middleware.Sessions[cookie.Value]
	userID := session["user_id"].(int)
	addressID := r.URL.Query().Get("id")
	_, err := database.DB.Exec("DELETE FROM addresses WHERE id = ? AND user_id = ?", addressID, userID)
	if err != nil {
		middleware.JSONError(w, "删除失败", http.StatusInternalServerError)
		return
	}
	middleware.JSON(w, map[string]interface{}{
		"success": true,
		"message": "地址已删除",
	})
}
