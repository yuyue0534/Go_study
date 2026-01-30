package models

import (
	"ecommerce-platform/internal/config"
	"time"
)

type Address struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	ReceiverName  string    `json:"receiver_name"`
	Phone         string    `json:"phone"`
	Province      string    `json:"province"`
	City          string    `json:"city"`
	District      string    `json:"district"`
	DetailAddress string    `json:"detail_address"`
	IsDefault     int       `json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateAddress 创建收货地址
func CreateAddress(address *Address) (int64, error) {
	// 如果设为默认，先取消其他默认
	if address.IsDefault == 1 {
		config.DB.Exec(`UPDATE addresses SET is_default = 0 WHERE user_id = ?`, address.UserID)
	}

	result, err := config.DB.Exec(`
		INSERT INTO addresses (user_id, receiver_name, phone, province, city, district, detail_address, is_default)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, address.UserID, address.ReceiverName, address.Phone, address.Province,
		address.City, address.District, address.DetailAddress, address.IsDefault)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetAddressByID 通过ID获取地址
func GetAddressByID(id int64) (*Address, error) {
	addr := &Address{}
	err := config.DB.QueryRow(`
		SELECT id, user_id, receiver_name, phone, province, city, district, detail_address, is_default, created_at
		FROM addresses WHERE id = ?
	`, id).Scan(&addr.ID, &addr.UserID, &addr.ReceiverName, &addr.Phone, &addr.Province,
		&addr.City, &addr.District, &addr.DetailAddress, &addr.IsDefault, &addr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// GetUserAddresses 获取用户所有地址
func GetUserAddresses(userID int64) ([]*Address, error) {
	rows, err := config.DB.Query(`
		SELECT id, user_id, receiver_name, phone, province, city, district, detail_address, is_default, created_at
		FROM addresses WHERE user_id = ? ORDER BY is_default DESC, created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []*Address
	for rows.Next() {
		addr := &Address{}
		err := rows.Scan(&addr.ID, &addr.UserID, &addr.ReceiverName, &addr.Phone, &addr.Province,
			&addr.City, &addr.District, &addr.DetailAddress, &addr.IsDefault, &addr.CreatedAt)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}
	return addresses, nil
}

// GetDefaultAddress 获取用户默认地址
func GetDefaultAddress(userID int64) (*Address, error) {
	addr := &Address{}
	err := config.DB.QueryRow(`
		SELECT id, user_id, receiver_name, phone, province, city, district, detail_address, is_default, created_at
		FROM addresses WHERE user_id = ? AND is_default = 1
	`, userID).Scan(&addr.ID, &addr.UserID, &addr.ReceiverName, &addr.Phone, &addr.Province,
		&addr.City, &addr.District, &addr.DetailAddress, &addr.IsDefault, &addr.CreatedAt)
	if err != nil {
		// 如果没有默认地址，返回第一个
		err = config.DB.QueryRow(`
			SELECT id, user_id, receiver_name, phone, province, city, district, detail_address, is_default, created_at
			FROM addresses WHERE user_id = ? ORDER BY created_at DESC LIMIT 1
		`, userID).Scan(&addr.ID, &addr.UserID, &addr.ReceiverName, &addr.Phone, &addr.Province,
			&addr.City, &addr.District, &addr.DetailAddress, &addr.IsDefault, &addr.CreatedAt)
		if err != nil {
			return nil, err
		}
	}
	return addr, nil
}

// UpdateAddress 更新地址
func UpdateAddress(address *Address) error {
	// 如果设为默认，先取消其他默认
	if address.IsDefault == 1 {
		config.DB.Exec(`UPDATE addresses SET is_default = 0 WHERE user_id = ? AND id != ?`, address.UserID, address.ID)
	}

	_, err := config.DB.Exec(`
		UPDATE addresses SET receiver_name = ?, phone = ?, province = ?, city = ?,
			district = ?, detail_address = ?, is_default = ?
		WHERE id = ?
	`, address.ReceiverName, address.Phone, address.Province, address.City,
		address.District, address.DetailAddress, address.IsDefault, address.ID)
	return err
}

// DeleteAddress 删除地址
func DeleteAddress(id int64) error {
	_, err := config.DB.Exec(`DELETE FROM addresses WHERE id = ?`, id)
	return err
}

// SetDefaultAddress 设置默认地址
func SetDefaultAddress(userID, addressID int64) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE addresses SET is_default = 0 WHERE user_id = ?`, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE addresses SET is_default = 1 WHERE id = ? AND user_id = ?`, addressID, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
