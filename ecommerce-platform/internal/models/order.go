package models

import (
	"database/sql"
	"ecommerce-platform/internal/config"
	"fmt"
	"strings"
	"time"
)

type Order struct {
	ID             int64     `json:"id"`
	OrderNo        string    `json:"order_no"`
	UserID         int64     `json:"user_id"`
	SellerID       int64     `json:"seller_id"`
	TotalAmount    float64   `json:"total_amount"`
	PayAmount      float64   `json:"pay_amount"`
	FreightAmount  float64   `json:"freight_amount"`
	Status         int       `json:"status"` // 0待支付 1待发货 2已发货 3已完成 4已取消 5退款中 6已退款
	PayType        string    `json:"pay_type"`
	PayTime        *time.Time `json:"pay_time"`
	ShipTime       *time.Time `json:"ship_time"`
	ReceiveTime    *time.Time `json:"receive_time"`
	FinishTime     *time.Time `json:"finish_time"`
	ReceiverName   string    `json:"receiver_name"`
	ReceiverPhone  string    `json:"receiver_phone"`
	ReceiverAddress string   `json:"receiver_address"`
	Remark         string    `json:"remark"`
	TrackingNo     string    `json:"tracking_no"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	// 关联信息
	Items      []*OrderItem `json:"items,omitempty"`
	SellerName string       `json:"seller_name,omitempty"`
	Username   string       `json:"username,omitempty"`
}

type OrderItem struct {
	ID           int64     `json:"id"`
	OrderID      int64     `json:"order_id"`
	ProductID    int64     `json:"product_id"`
	ProductName  string    `json:"product_name"`
	ProductImage string    `json:"product_image"`
	Price        float64   `json:"price"`
	Quantity     int       `json:"quantity"`
	TotalPrice   float64   `json:"total_price"`
	CreatedAt    time.Time `json:"created_at"`
}

// GenerateOrderNo 生成订单号
func GenerateOrderNo() string {
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), time.Now().Unix()%1000)
}

// CreateOrder 创建订单
func CreateOrder(order *Order) (int64, error) {
	tx, err := config.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// 创建订单
	result, err := tx.Exec(`
		INSERT INTO orders (order_no, user_id, seller_id, total_amount, pay_amount, freight_amount,
			status, receiver_name, receiver_phone, receiver_address, remark)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, order.OrderNo, order.UserID, order.SellerID, order.TotalAmount, order.PayAmount,
		order.FreightAmount, order.Status, order.ReceiverName, order.ReceiverPhone,
		order.ReceiverAddress, order.Remark)
	if err != nil {
		return 0, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 创建订单商品
	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO order_items (order_id, product_id, product_name, product_image, price, quantity, total_price)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, orderID, item.ProductID, item.ProductName, item.ProductImage, item.Price, item.Quantity, item.TotalPrice)
		if err != nil {
			return 0, err
		}

		// 减少库存，增加销量
		_, err = tx.Exec(`UPDATE products SET stock = stock - ?, sales = sales + ? WHERE id = ?`,
			item.Quantity, item.Quantity, item.ProductID)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// GetOrderByID 通过ID获取订单
func GetOrderByID(id int64) (*Order, error) {
	order := &Order{}
	var payTime, shipTime, receiveTime, finishTime sql.NullTime
	err := config.DB.QueryRow(`
		SELECT o.id, o.order_no, o.user_id, o.seller_id, o.total_amount, o.pay_amount, o.freight_amount,
			o.status, o.pay_type, o.pay_time, o.ship_time, o.receive_time, o.finish_time,
			o.receiver_name, o.receiver_phone, o.receiver_address, o.remark, o.tracking_no,
			o.created_at, o.updated_at, s.shop_name, u.username
		FROM orders o
		LEFT JOIN sellers s ON o.seller_id = s.id
		LEFT JOIN users u ON o.user_id = u.id
		WHERE o.id = ?
	`, id).Scan(&order.ID, &order.OrderNo, &order.UserID, &order.SellerID, &order.TotalAmount,
		&order.PayAmount, &order.FreightAmount, &order.Status, &order.PayType,
		&payTime, &shipTime, &receiveTime, &finishTime,
		&order.ReceiverName, &order.ReceiverPhone, &order.ReceiverAddress, &order.Remark,
		&order.TrackingNo, &order.CreatedAt, &order.UpdatedAt, &order.SellerName, &order.Username)
	if err != nil {
		return nil, err
	}
	
	if payTime.Valid {
		order.PayTime = &payTime.Time
	}
	if shipTime.Valid {
		order.ShipTime = &shipTime.Time
	}
	if receiveTime.Valid {
		order.ReceiveTime = &receiveTime.Time
	}
	if finishTime.Valid {
		order.FinishTime = &finishTime.Time
	}

	// 获取订单商品
	order.Items, err = GetOrderItems(order.ID)
	return order, err
}

// GetOrderByOrderNo 通过订单号获取订单
func GetOrderByOrderNo(orderNo string) (*Order, error) {
	var id int64
	err := config.DB.QueryRow(`SELECT id FROM orders WHERE order_no = ?`, orderNo).Scan(&id)
	if err != nil {
		return nil, err
	}
	return GetOrderByID(id)
}

// GetOrderItems 获取订单商品
func GetOrderItems(orderID int64) ([]*OrderItem, error) {
	rows, err := config.DB.Query(`
		SELECT id, order_id, product_id, product_name, product_image, price, quantity, total_price, created_at
		FROM order_items WHERE order_id = ?
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*OrderItem
	for rows.Next() {
		item := &OrderItem{}
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductName,
			&item.ProductImage, &item.Price, &item.Quantity, &item.TotalPrice, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// GetOrders 获取订单列表
func GetOrders(page, pageSize int, userID, sellerID int64, status int) ([]*Order, int, error) {
	offset := (page - 1) * pageSize
	
	conditions := []string{"1=1"}
	args := []interface{}{}
	
	if userID > 0 {
		conditions = append(conditions, "o.user_id = ?")
		args = append(args, userID)
	}
	if sellerID > 0 {
		conditions = append(conditions, "o.seller_id = ?")
		args = append(args, sellerID)
	}
	if status >= 0 {
		conditions = append(conditions, "o.status = ?")
		args = append(args, status)
	}
	
	whereClause := strings.Join(conditions, " AND ")
	
	var total int
	err := config.DB.QueryRow(`SELECT COUNT(*) FROM orders o WHERE `+whereClause, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	queryArgs := append(args, pageSize, offset)
	rows, err := config.DB.Query(`
		SELECT o.id, o.order_no, o.user_id, o.seller_id, o.total_amount, o.pay_amount, o.freight_amount,
			o.status, o.pay_type, o.pay_time, o.ship_time, o.receive_time, o.finish_time,
			o.receiver_name, o.receiver_phone, o.receiver_address, o.remark, o.tracking_no,
			o.created_at, o.updated_at, s.shop_name, u.username
		FROM orders o
		LEFT JOIN sellers s ON o.seller_id = s.id
		LEFT JOIN users u ON o.user_id = u.id
		WHERE `+whereClause+` ORDER BY o.created_at DESC LIMIT ? OFFSET ?`, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		order := &Order{}
		var payTime, shipTime, receiveTime, finishTime sql.NullTime
		err := rows.Scan(&order.ID, &order.OrderNo, &order.UserID, &order.SellerID, &order.TotalAmount,
			&order.PayAmount, &order.FreightAmount, &order.Status, &order.PayType,
			&payTime, &shipTime, &receiveTime, &finishTime,
			&order.ReceiverName, &order.ReceiverPhone, &order.ReceiverAddress, &order.Remark,
			&order.TrackingNo, &order.CreatedAt, &order.UpdatedAt, &order.SellerName, &order.Username)
		if err != nil {
			return nil, 0, err
		}
		
		if payTime.Valid {
			order.PayTime = &payTime.Time
		}
		if shipTime.Valid {
			order.ShipTime = &shipTime.Time
		}
		if receiveTime.Valid {
			order.ReceiveTime = &receiveTime.Time
		}
		if finishTime.Valid {
			order.FinishTime = &finishTime.Time
		}

		// 获取订单商品
		order.Items, _ = GetOrderItems(order.ID)
		orders = append(orders, order)
	}
	return orders, total, nil
}

// UpdateOrderStatus 更新订单状态
func UpdateOrderStatus(id int64, status int) error {
	now := time.Now()
	var err error
	switch status {
	case 1: // 待发货（已支付）
		_, err = config.DB.Exec(`UPDATE orders SET status = ?, pay_time = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, now, id)
	case 2: // 已发货
		_, err = config.DB.Exec(`UPDATE orders SET status = ?, ship_time = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, now, id)
	case 3: // 已完成
		_, err = config.DB.Exec(`UPDATE orders SET status = ?, receive_time = ?, finish_time = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, now, now, id)
	default:
		_, err = config.DB.Exec(`UPDATE orders SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
	}
	return err
}

// PayOrder 支付订单
func PayOrder(id int64, payType string) error {
	_, err := config.DB.Exec(`
		UPDATE orders SET status = 1, pay_type = ?, pay_time = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND status = 0
	`, payType, id)
	return err
}

// ShipOrder 发货
func ShipOrder(id int64, trackingNo string) error {
	_, err := config.DB.Exec(`
		UPDATE orders SET status = 2, tracking_no = ?, ship_time = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND status = 1
	`, trackingNo, id)
	return err
}

// ReceiveOrder 确认收货
func ReceiveOrder(id int64) error {
	_, err := config.DB.Exec(`
		UPDATE orders SET status = 3, receive_time = CURRENT_TIMESTAMP, finish_time = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND status = 2
	`, id)
	return err
}

// CancelOrder 取消订单
func CancelOrder(id int64) error {
	// 获取订单信息
	order, err := GetOrderByID(id)
	if err != nil {
		return err
	}
	if order.Status != 0 {
		return fmt.Errorf("订单状态不允许取消")
	}

	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 更新订单状态
	_, err = tx.Exec(`UPDATE orders SET status = 4, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	if err != nil {
		return err
	}

	// 恢复库存
	for _, item := range order.Items {
		_, err = tx.Exec(`UPDATE products SET stock = stock + ?, sales = sales - ? WHERE id = ?`,
			item.Quantity, item.Quantity, item.ProductID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// UpdateOrderAddress 更新订单地址
func UpdateOrderAddress(id int64, name, phone, address string) error {
	_, err := config.DB.Exec(`
		UPDATE orders SET receiver_name = ?, receiver_phone = ?, receiver_address = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND status = 0
	`, name, phone, address, id)
	return err
}

// GetOrderStats 获取订单统计
func GetOrderStats(sellerID int64) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	var condition string
	var args []interface{}
	if sellerID > 0 {
		condition = " WHERE seller_id = ?"
		args = append(args, sellerID)
	}

	// 待支付
	var pending int
	config.DB.QueryRow(`SELECT COUNT(*) FROM orders`+condition+` AND status = 0`, args...).Scan(&pending)
	stats["pending"] = pending

	// 待发货
	var toShip int
	config.DB.QueryRow(`SELECT COUNT(*) FROM orders`+condition+` AND status = 1`, args...).Scan(&toShip)
	stats["to_ship"] = toShip

	// 已发货
	var shipped int
	config.DB.QueryRow(`SELECT COUNT(*) FROM orders`+condition+` AND status = 2`, args...).Scan(&shipped)
	stats["shipped"] = shipped

	// 已完成
	var completed int
	config.DB.QueryRow(`SELECT COUNT(*) FROM orders`+condition+` AND status = 3`, args...).Scan(&completed)
	stats["completed"] = completed

	// 总销售额
	var totalSales float64
	config.DB.QueryRow(`SELECT COALESCE(SUM(pay_amount), 0) FROM orders`+condition+` AND status >= 1`, args...).Scan(&totalSales)
	stats["total_sales"] = totalSales

	return stats, nil
}
