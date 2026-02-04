package service

import (
	"crm-erp-system/database"
	"crm-erp-system/model"
	"errors"
)

type ProductService struct{}

func (s *ProductService) Create(product *model.Product) (int64, error) {
	// 检查SKU是否存在
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM products WHERE sku = ?", product.SKU).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("SKU已存在")
	}

	result, err := database.DB.Exec(
		"INSERT INTO products (name, sku, description, price, cost, category) VALUES (?, ?, ?, ?, ?, ?)",
		product.Name, product.SKU, product.Description, product.Price, product.Cost, product.Category,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *ProductService) GetByID(id int64) (*model.Product, error) {
	var product model.Product
	err := database.DB.QueryRow(
		"SELECT id, name, sku, description, price, cost, category, created_at, updated_at FROM products WHERE id = ?",
		id,
	).Scan(&product.ID, &product.Name, &product.SKU, &product.Description, &product.Price, &product.Cost, &product.Category, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return nil, errors.New("产品不存在")
	}
	return &product, nil
}

func (s *ProductService) List(page, pageSize int) ([]model.Product, error) {
	offset := (page - 1) * pageSize
	rows, err := database.DB.Query(
		"SELECT id, name, sku, description, price, cost, category, created_at, updated_at FROM products ORDER BY created_at DESC LIMIT ? OFFSET ?",
		pageSize, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.SKU, &product.Description, &product.Price, &product.Cost, &product.Category, &product.CreatedAt, &product.UpdatedAt); err != nil {
			continue
		}
		products = append(products, product)
	}
	return products, nil
}

func (s *ProductService) Update(id int64, product *model.Product) error {
	result, err := database.DB.Exec(
		"UPDATE products SET name=?, description=?, price=?, cost=?, category=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		product.Name, product.Description, product.Price, product.Cost, product.Category, id,
	)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("产品不存在")
	}
	return nil
}

func (s *ProductService) Delete(id int64) error {
	result, err := database.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("产品不存在")
	}
	return nil
}
