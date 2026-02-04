package service

import (
	"crm-erp-system/database"
	"crm-erp-system/model"
	"errors"
)

type CustomerService struct{}

func (s *CustomerService) Create(customer *model.Customer, userID int64) (int64, error) {
	customer.UserID = userID
	customer.Status = "active"

	result, err := database.DB.Exec(
		"INSERT INTO customers (name, company, email, phone, address, status, user_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		customer.Name, customer.Company, customer.Email, customer.Phone, customer.Address, customer.Status, customer.UserID,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *CustomerService) GetByID(id int64) (*model.Customer, error) {
	var customer model.Customer
	err := database.DB.QueryRow(
		"SELECT id, name, company, email, phone, address, status, user_id, created_at, updated_at FROM customers WHERE id = ?",
		id,
	).Scan(&customer.ID, &customer.Name, &customer.Company, &customer.Email, &customer.Phone, &customer.Address, &customer.Status, &customer.UserID, &customer.CreatedAt, &customer.UpdatedAt)

	if err != nil {
		return nil, errors.New("客户不存在")
	}
	return &customer, nil
}

func (s *CustomerService) List(page, pageSize int) ([]model.Customer, error) {
	offset := (page - 1) * pageSize
	rows, err := database.DB.Query(
		"SELECT id, name, company, email, phone, address, status, user_id, created_at, updated_at FROM customers ORDER BY created_at DESC LIMIT ? OFFSET ?",
		pageSize, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []model.Customer
	for rows.Next() {
		var customer model.Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Company, &customer.Email, &customer.Phone, &customer.Address, &customer.Status, &customer.UserID, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
			continue
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (s *CustomerService) Update(id int64, customer *model.Customer) error {
	result, err := database.DB.Exec(
		"UPDATE customers SET name=?, company=?, email=?, phone=?, address=?, status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		customer.Name, customer.Company, customer.Email, customer.Phone, customer.Address, customer.Status, id,
	)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("客户不存在")
	}
	return nil
}

func (s *CustomerService) Delete(id int64) error {
	result, err := database.DB.Exec("DELETE FROM customers WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("客户不存在")
	}
	return nil
}
