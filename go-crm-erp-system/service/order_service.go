package service

import (
	"crm-erp-system/database"
	"crm-erp-system/model"
	"errors"
	"fmt"
	"time"
)

type OrderService struct{}

func (s *OrderService) Create(order *model.Order, userID int64) (int64, error) {
	// 生成订单号
	order.OrderNo = fmt.Sprintf("ORD%d", time.Now().UnixNano())
	order.UserID = userID
	order.Status = "pending"
	order.TotalAmount = order.UnitPrice * float64(order.Quantity)

	// 检查客户是否存在
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM customers WHERE id = ?", order.CustomerID).Scan(&count)
	if err != nil || count == 0 {
		return 0, errors.New("客户不存在")
	}

	// 检查产品是否存在
	err = database.DB.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", order.ProductID).Scan(&count)
	if err != nil || count == 0 {
		return 0, errors.New("产品不存在")
	}

	result, err := database.DB.Exec(
		"INSERT INTO orders (order_no, customer_id, product_id, quantity, unit_price, total_amount, status, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		order.OrderNo, order.CustomerID, order.ProductID, order.Quantity, order.UnitPrice, order.TotalAmount, order.Status, order.UserID,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *OrderService) GetByID(id int64) (*model.Order, error) {
	var order model.Order
	err := database.DB.QueryRow(
		"SELECT id, order_no, customer_id, product_id, quantity, unit_price, total_amount, status, user_id, created_at, updated_at FROM orders WHERE id = ?",
		id,
	).Scan(&order.ID, &order.OrderNo, &order.CustomerID, &order.ProductID, &order.Quantity, &order.UnitPrice, &order.TotalAmount, &order.Status, &order.UserID, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		return nil, errors.New("订单不存在")
	}
	return &order, nil
}

func (s *OrderService) List(page, pageSize int) ([]model.Order, error) {
	offset := (page - 1) * pageSize
	rows, err := database.DB.Query(
		"SELECT id, order_no, customer_id, product_id, quantity, unit_price, total_amount, status, user_id, created_at, updated_at FROM orders ORDER BY created_at DESC LIMIT ? OFFSET ?",
		pageSize, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.ID, &order.OrderNo, &order.CustomerID, &order.ProductID, &order.Quantity, &order.UnitPrice, &order.TotalAmount, &order.Status, &order.UserID, &order.CreatedAt, &order.UpdatedAt); err != nil {
			continue
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (s *OrderService) UpdateStatus(id int64, status string) error {
	validStatuses := map[string]bool{
		"pending":   true,
		"confirmed": true,
		"shipped":   true,
		"completed": true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return errors.New("无效的订单状态")
	}

	result, err := database.DB.Exec(
		"UPDATE orders SET status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		status, id,
	)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("订单不存在")
	}
	return nil
}

func (s *OrderService) Delete(id int64) error {
	result, err := database.DB.Exec("DELETE FROM orders WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("订单不存在")
	}
	return nil
}
