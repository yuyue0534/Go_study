package database

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite", "./ecommerce.db")
	if err != nil {
		return err
	}
	if err = DB.Ping(); err != nil {
		return err
	}
	if err = createTables(); err != nil {
		return err
	}
	if err = insertInitialData(); err != nil {
		return err
	}
	log.Println("数据库初始化成功")
	return nil
}

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			phone TEXT,
			role TEXT DEFAULT 'customer',
			status TEXT DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS addresses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			receiver_name TEXT NOT NULL,
			phone TEXT NOT NULL,
			province TEXT NOT NULL,
			city TEXT NOT NULL,
			district TEXT NOT NULL,
			detail_address TEXT NOT NULL,
			is_default INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			description TEXT,
			parent_id INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			price REAL NOT NULL,
			stock INTEGER DEFAULT 0,
			category_id INTEGER,
			seller_id INTEGER,
			brand TEXT,
			image_url TEXT,
			status TEXT DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (category_id) REFERENCES categories(id),
			FOREIGN KEY (seller_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS cart_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_no TEXT UNIQUE NOT NULL,
			user_id INTEGER NOT NULL,
			total_amount REAL NOT NULL,
			status TEXT DEFAULT 'pending_payment',
			payment_method TEXT,
			shipping_address TEXT NOT NULL,
			receiver_name TEXT NOT NULL,
			receiver_phone TEXT NOT NULL,
			shipping_method TEXT DEFAULT 'standard',
			tracking_number TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			paid_at TIMESTAMP,
			shipped_at TIMESTAMP,
			completed_at TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS order_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			product_name TEXT NOT NULL,
			product_price REAL NOT NULL,
			quantity INTEGER NOT NULL,
			subtotal REAL NOT NULL,
			FOREIGN KEY (order_id) REFERENCES orders(id),
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`,
		`CREATE TABLE IF NOT EXISTS reviews (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			order_id INTEGER NOT NULL,
			rating INTEGER NOT NULL CHECK(rating >= 1 AND rating <= 5),
			comment TEXT,
			status TEXT DEFAULT 'approved',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (product_id) REFERENCES products(id),
			FOREIGN KEY (order_id) REFERENCES orders(id)
		)`,
		`CREATE TABLE IF NOT EXISTS review_replies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			review_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (review_id) REFERENCES reviews(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
	}
	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("创建表失败: %v", err)
		}
	}
	return nil
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func insertInitialData() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	users := []struct {
		username, password, email, role string
	}{
		{"admin", "admin123", "admin@ecommerce.com", "admin"},
		{"seller1", "seller123", "seller1@ecommerce.com", "seller"},
		{"customer1", "customer123", "customer1@ecommerce.com", "customer"},
	}
	for _, u := range users {
		_, err := DB.Exec(
			"INSERT INTO users (username, password, email, role) VALUES (?, ?, ?, ?)",
			u.username, hashPassword(u.password), u.email, u.role,
		)
		if err != nil {
			return err
		}
	}
	categories := []struct {
		name, description string
	}{
		{"电子产品", "各类电子设备和配件"},
		{"服装鞋帽", "男女服装、鞋子、配饰"},
		{"家居生活", "家居用品、装饰品"},
		{"图书音像", "图书、电子书、音乐"},
		{"食品饮料", "零食、饮料、生鲜"},
		{"运动户外", "运动器材、户外用品"},
	}
	for _, c := range categories {
		_, err := DB.Exec(
			"INSERT INTO categories (name, description) VALUES (?, ?)",
			c.name, c.description,
		)
		if err != nil {
			return err
		}
	}
	products := []struct {
		name, description, brand, imageURL, status string
		price                                      float64
		stock, categoryID, sellerID                int
	}{
		{"iPhone 15 Pro", "苹果最新旗舰手机，性能强劲", "Apple", "https://via.placeholder.com/300x300?text=iPhone+15", "approved", 7999.00, 50, 1, 2},
		{"MacBook Air M2", "轻薄便携笔记本电脑", "Apple", "https://via.placeholder.com/300x300?text=MacBook", "approved", 8999.00, 30, 1, 2},
		{"AirPods Pro", "主动降噪无线耳机", "Apple", "https://via.placeholder.com/300x300?text=AirPods", "approved", 1899.00, 100, 1, 2},
		{"男士休闲T恤", "纯棉舒适透气", "Uniqlo", "https://via.placeholder.com/300x300?text=T-Shirt", "approved", 99.00, 200, 2, 2},
		{"运动鞋", "缓震跑步鞋", "Nike", "https://via.placeholder.com/300x300?text=Shoes", "approved", 399.00, 80, 2, 2},
		{"北欧风台灯", "简约设计护眼台灯", "IKEA", "https://via.placeholder.com/300x300?text=Lamp", "approved", 299.00, 60, 3, 2},
		{"Python编程从入门到实践", "编程学习经典教材", "人民邮电出版社", "https://via.placeholder.com/300x300?text=Book", "approved", 89.00, 150, 4, 2},
		{"瑜伽垫", "防滑环保瑜伽垫", "Lululemon", "https://via.placeholder.com/300x300?text=Yoga+Mat", "approved", 129.00, 120, 6, 2},
		{"有机绿茶", "精选茶叶，清香怡人", "天福茗茶", "https://via.placeholder.com/300x300?text=Tea", "approved", 158.00, 90, 5, 2},
		{"无线键盘鼠标套装", "办公游戏两用", "Logitech", "https://via.placeholder.com/300x300?text=Keyboard", "approved", 199.00, 70, 1, 2},
	}
	for _, p := range products {
		_, err := DB.Exec(
			`INSERT INTO products (name, description, price, stock, category_id, seller_id, brand, image_url, status, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			p.name, p.description, p.price, p.stock, p.categoryID, p.sellerID, p.brand, p.imageURL, p.status, time.Now(),
		)
		if err != nil {
			return err
		}
	}
	log.Println("初始数据插入成功")
	return nil
}
