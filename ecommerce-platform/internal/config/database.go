package config

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDatabase() {
	// 确保数据目录存在
	dataDir := "./data"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.MkdirAll(dataDir, 0755)
	}

	dbPath := filepath.Join(dataDir, "ecommerce.db")
	var err error
	DB, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 设置连接池
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	// 创建表
	createTables()
}

func createTables() {
	// 用户表
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			phone VARCHAR(20),
			role VARCHAR(20) DEFAULT 'customer',
			avatar VARCHAR(255),
			status INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("创建用户表失败:", err)
	}

	// 商家信息表
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS sellers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER UNIQUE NOT NULL,
			shop_name VARCHAR(100) NOT NULL,
			shop_description TEXT,
			shop_logo VARCHAR(255),
			status INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatal("创建商家表失败:", err)
	}

	// 商品分类表
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(50) NOT NULL,
			parent_id INTEGER DEFAULT 0,
			icon VARCHAR(100),
			sort_order INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("创建分类表失败:", err)
	}

	// 商品表
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			seller_id INTEGER NOT NULL,
			category_id INTEGER,
			name VARCHAR(200) NOT NULL,
			description TEXT,
			price DECIMAL(10,2) NOT NULL,
			original_price DECIMAL(10,2),
			stock INTEGER DEFAULT 0,
			sales INTEGER DEFAULT 0,
			images TEXT,
			brand VARCHAR(100),
			status INTEGER DEFAULT 0,
			rating DECIMAL(2,1) DEFAULT 5.0,
			rating_count INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (seller_id) REFERENCES sellers(id),
			FOREIGN KEY (category_id) REFERENCES categories(id)
		)
	`)
	if err != nil {
		log.Fatal("创建商品表失败:", err)
	}

	// 购物车表
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS cart_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER DEFAULT 1,
			selected INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (product_id) REFERENCES products(id),
			UNIQUE(user_id, product_id)
		)
	`)
	if err != nil {
		log.Fatal("创建购物车表失败:", err)
	}

	// 收货地址表
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS addresses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			receiver_name VARCHAR(50) NOT NULL,
			phone VARCHAR(20) NOT NULL,
			province VARCHAR(50),
			city VARCHAR(50),
			district VARCHAR(50),
			detail_address TEXT NOT NULL,
			is_default INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatal("创建地址表失败:", err)
	}

	// 订单表
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_no VARCHAR(50) UNIQUE NOT NULL,
			user_id INTEGER NOT NULL,
			seller_id INTEGER NOT NULL,
			total_amount DECIMAL(10,2) NOT NULL,
			pay_amount DECIMAL(10,2) NOT NULL,
			freight_amount DECIMAL(10,2) DEFAULT 0,
			status INTEGER DEFAULT 0,
			pay_type VARCHAR(20),
			pay_time DATETIME,
			ship_time DATETIME,
			receive_time DATETIME,
			finish_time DATETIME,
			receiver_name VARCHAR(50),
			receiver_phone VARCHAR(20),
			receiver_address TEXT,
			remark TEXT,
			tracking_no VARCHAR(50),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (seller_id) REFERENCES sellers(id)
		)
	`)
	if err != nil {
		log.Fatal("创建订单表失败:", err)
	}

	// 订单商品表
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS order_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			product_name VARCHAR(200),
			product_image VARCHAR(255),
			price DECIMAL(10,2) NOT NULL,
			quantity INTEGER NOT NULL,
			total_price DECIMAL(10,2) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (order_id) REFERENCES orders(id),
			FOREIGN KEY (product_id) REFERENCES products(id)
		)
	`)
	if err != nil {
		log.Fatal("创建订单商品表失败:", err)
	}

	// 商品评价表
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS reviews (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			rating INTEGER NOT NULL,
			content TEXT,
			images TEXT,
			reply TEXT,
			reply_time DATETIME,
			status INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (order_id) REFERENCES orders(id),
			FOREIGN KEY (product_id) REFERENCES products(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatal("创建评价表失败:", err)
	}

	// 系统消息表
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title VARCHAR(200),
			content TEXT,
			type VARCHAR(20),
			is_read INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatal("创建消息表失败:", err)
	}

	log.Println("数据库表创建成功")
}

func CloseDatabase() {
	if DB != nil {
		DB.Close()
	}
}
