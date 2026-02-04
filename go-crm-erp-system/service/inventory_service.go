package service

import (
	"crm-erp-system/database"
	"crm-erp-system/model"
	"errors"
)

type InventoryService struct{}

func (s *InventoryService) Create(inventory *model.Inventory) (int64, error) {
	// 检查产品是否存在
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", inventory.ProductID).Scan(&count)
	if err != nil || count == 0 {
		return 0, errors.New("产品不存在")
	}

	result, err := database.DB.Exec(
		"INSERT INTO inventory (product_id, quantity, warehouse) VALUES (?, ?, ?)",
		inventory.ProductID, inventory.Quantity, inventory.Warehouse,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *InventoryService) GetByProductID(productID int64) (*model.Inventory, error) {
	var inventory model.Inventory
	err := database.DB.QueryRow(
		"SELECT id, product_id, quantity, warehouse, updated_at FROM inventory WHERE product_id = ?",
		productID,
	).Scan(&inventory.ID, &inventory.ProductID, &inventory.Quantity, &inventory.Warehouse, &inventory.UpdatedAt)

	if err != nil {
		return nil, errors.New("库存记录不存在")
	}
	return &inventory, nil
}

func (s *InventoryService) Update(productID int64, quantity int) error {
	result, err := database.DB.Exec(
		"UPDATE inventory SET quantity=?, updated_at=CURRENT_TIMESTAMP WHERE product_id=?",
		quantity, productID,
	)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("库存记录不存在")
	}
	return nil
}

func (s *InventoryService) List() ([]model.Inventory, error) {
	rows, err := database.DB.Query(
		"SELECT id, product_id, quantity, warehouse, updated_at FROM inventory ORDER BY updated_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventories []model.Inventory
	for rows.Next() {
		var inventory model.Inventory
		if err := rows.Scan(&inventory.ID, &inventory.ProductID, &inventory.Quantity, &inventory.Warehouse, &inventory.UpdatedAt); err != nil {
			continue
		}
		inventories = append(inventories, inventory)
	}
	return inventories, nil
}
