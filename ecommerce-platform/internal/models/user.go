package models

import (
	"database/sql"
	"ecommerce-platform/internal/config"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	Avatar    string    `json:"avatar"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Seller struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	ShopName        string    `json:"shop_name"`
	ShopDescription string    `json:"shop_description"`
	ShopLogo        string    `json:"shop_logo"`
	Status          int       `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

// CreateUser 创建用户
func CreateUser(user *User) (int64, error) {
	result, err := config.DB.Exec(`
		INSERT INTO users (username, password, email, phone, role, avatar, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, user.Username, user.Password, user.Email, user.Phone, user.Role, user.Avatar, user.Status)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetUserByUsername 通过用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := config.DB.QueryRow(`
		SELECT id, username, password, email, phone, role, avatar, status, created_at, updated_at
		FROM users WHERE username = ?
	`, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone,
		&user.Role, &user.Avatar, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID 通过ID获取用户
func GetUserByID(id int64) (*User, error) {
	user := &User{}
	err := config.DB.QueryRow(`
		SELECT id, username, password, email, phone, role, avatar, status, created_at, updated_at
		FROM users WHERE id = ?
	`, id).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone,
		&user.Role, &user.Avatar, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail 通过邮箱获取用户
func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := config.DB.QueryRow(`
		SELECT id, username, password, email, phone, role, avatar, status, created_at, updated_at
		FROM users WHERE email = ?
	`, email).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone,
		&user.Role, &user.Avatar, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(user *User) error {
	_, err := config.DB.Exec(`
		UPDATE users SET email = ?, phone = ?, avatar = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, user.Email, user.Phone, user.Avatar, user.ID)
	return err
}

// UpdateUserPassword 更新密码
func UpdateUserPassword(userID int64, password string) error {
	_, err := config.DB.Exec(`
		UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
	`, password, userID)
	return err
}

// GetAllUsers 获取所有用户（管理员用）
func GetAllUsers(page, pageSize int, role string) ([]*User, int, error) {
	offset := (page - 1) * pageSize
	
	var total int
	var rows *sql.Rows
	var err error
	
	if role != "" {
		err = config.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE role = ?`, role).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
		rows, err = config.DB.Query(`
			SELECT id, username, password, email, phone, role, avatar, status, created_at, updated_at
			FROM users WHERE role = ? ORDER BY created_at DESC LIMIT ? OFFSET ?
		`, role, pageSize, offset)
	} else {
		err = config.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
		rows, err = config.DB.Query(`
			SELECT id, username, password, email, phone, role, avatar, status, created_at, updated_at
			FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?
		`, pageSize, offset)
	}
	
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone,
			&user.Role, &user.Avatar, &user.Status, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	return users, total, nil
}

// UpdateUserStatus 更新用户状态
func UpdateUserStatus(userID int64, status int) error {
	_, err := config.DB.Exec(`UPDATE users SET status = ? WHERE id = ?`, status, userID)
	return err
}

// UpdateUserRole 更新用户角色
func UpdateUserRole(userID int64, role string) error {
	_, err := config.DB.Exec(`UPDATE users SET role = ? WHERE id = ?`, role, userID)
	return err
}

// CreateSeller 创建商家信息
func CreateSeller(seller *Seller) (int64, error) {
	result, err := config.DB.Exec(`
		INSERT INTO sellers (user_id, shop_name, shop_description, shop_logo, status)
		VALUES (?, ?, ?, ?, ?)
	`, seller.UserID, seller.ShopName, seller.ShopDescription, seller.ShopLogo, seller.Status)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetSellerByUserID 通过用户ID获取商家信息
func GetSellerByUserID(userID int64) (*Seller, error) {
	seller := &Seller{}
	err := config.DB.QueryRow(`
		SELECT id, user_id, shop_name, shop_description, shop_logo, status, created_at
		FROM sellers WHERE user_id = ?
	`, userID).Scan(&seller.ID, &seller.UserID, &seller.ShopName, &seller.ShopDescription,
		&seller.ShopLogo, &seller.Status, &seller.CreatedAt)
	if err != nil {
		return nil, err
	}
	return seller, nil
}

// GetSellerByID 通过ID获取商家信息
func GetSellerByID(id int64) (*Seller, error) {
	seller := &Seller{}
	err := config.DB.QueryRow(`
		SELECT id, user_id, shop_name, shop_description, shop_logo, status, created_at
		FROM sellers WHERE id = ?
	`, id).Scan(&seller.ID, &seller.UserID, &seller.ShopName, &seller.ShopDescription,
		&seller.ShopLogo, &seller.Status, &seller.CreatedAt)
	if err != nil {
		return nil, err
	}
	return seller, nil
}

// UpdateSeller 更新商家信息
func UpdateSeller(seller *Seller) error {
	_, err := config.DB.Exec(`
		UPDATE sellers SET shop_name = ?, shop_description = ?, shop_logo = ?
		WHERE id = ?
	`, seller.ShopName, seller.ShopDescription, seller.ShopLogo, seller.ID)
	return err
}

// UpdateSellerStatus 更新商家状态
func UpdateSellerStatus(sellerID int64, status int) error {
	_, err := config.DB.Exec(`UPDATE sellers SET status = ? WHERE id = ?`, status, sellerID)
	return err
}

// GetAllSellers 获取所有商家
func GetAllSellers(page, pageSize int) ([]*Seller, int, error) {
	offset := (page - 1) * pageSize
	
	var total int
	err := config.DB.QueryRow(`SELECT COUNT(*) FROM sellers`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := config.DB.Query(`
		SELECT id, user_id, shop_name, shop_description, shop_logo, status, created_at
		FROM sellers ORDER BY created_at DESC LIMIT ? OFFSET ?
	`, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sellers []*Seller
	for rows.Next() {
		seller := &Seller{}
		err := rows.Scan(&seller.ID, &seller.UserID, &seller.ShopName, &seller.ShopDescription,
			&seller.ShopLogo, &seller.Status, &seller.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		sellers = append(sellers, seller)
	}
	return sellers, total, nil
}
