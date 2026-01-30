package main

import (
	"ecommerce-platform/internal/config"
	"ecommerce-platform/internal/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 初始化数据库
	config.InitDatabase()
	defer config.CloseDatabase()

	// 确保上传目录存在
	os.MkdirAll("./uploads", 0755)

	// 初始化示例数据
	initSampleData()

	// 设置路由
	handler := routes.SetupRoutes()

	// 启动服务器
	port := ":8080"
	fmt.Printf("🚀 电商平台启动成功！\n")
	fmt.Printf("📍 访问地址: http://localhost%s\n", port)
	fmt.Printf("👤 管理员账号: admin / admin123\n")
	fmt.Printf("🏪 商家账号: seller / seller123\n")
	fmt.Printf("🛒 顾客账号: customer / customer123\n")
	fmt.Println("----------------------------------------")

	log.Fatal(http.ListenAndServe(port, handler))
}

func initSampleData() {
	// 检查是否已初始化
	var count int
	config.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count > 0 {
		return
	}

	log.Println("初始化示例数据...")

	// 创建管理员
	adminPassword, _ := hashPassword("admin123")
	config.DB.Exec(`
		INSERT INTO users (username, password, email, phone, role, status)
		VALUES ('admin', ?, 'admin@shop.com', '13800000000', 'admin', 1)
	`, adminPassword)

	// 创建商家用户
	sellerPassword, _ := hashPassword("seller123")
	result, _ := config.DB.Exec(`
		INSERT INTO users (username, password, email, phone, role, status)
		VALUES ('seller', ?, 'seller@shop.com', '13800000001', 'seller', 1)
	`, sellerPassword)
	sellerUserID, _ := result.LastInsertId()

	// 创建商家信息
	result, _ = config.DB.Exec(`
		INSERT INTO sellers (user_id, shop_name, shop_description, status)
		VALUES (?, '优品旗舰店', '专注品质，用心服务', 1)
	`, sellerUserID)
	sellerID, _ := result.LastInsertId()

	// 创建顾客
	customerPassword, _ := hashPassword("customer123")
	result, _ = config.DB.Exec(`
		INSERT INTO users (username, password, email, phone, role, status)
		VALUES ('customer', ?, 'customer@shop.com', '13800000002', 'customer', 1)
	`, customerPassword)
	customerUserID, _ := result.LastInsertId()

	// 创建顾客地址
	config.DB.Exec(`
		INSERT INTO addresses (user_id, receiver_name, phone, province, city, district, detail_address, is_default)
		VALUES (?, '张三', '13800000002', '广东省', '深圳市', '南山区', '科技园路100号', 1)
	`, customerUserID)

	// 创建分类
	categories := []struct {
		name string
		icon string
	}{
		{"手机数码", "📱"},
		{"电脑办公", "💻"},
		{"家用电器", "🏠"},
		{"服装鞋包", "👔"},
		{"美妆护肤", "💄"},
		{"食品饮料", "🍔"},
		{"图书文具", "📚"},
		{"运动户外", "⚽"},
	}

	var categoryIDs []int64
	for _, cat := range categories {
		result, _ := config.DB.Exec(`
			INSERT INTO categories (name, icon) VALUES (?, ?)
		`, cat.name, cat.icon)
		id, _ := result.LastInsertId()
		categoryIDs = append(categoryIDs, id)
	}

	// 创建示例商品
	products := []struct {
		name          string
		description   string
		price         float64
		originalPrice float64
		stock         int
		categoryIndex int
		images        string
	}{
		{"iPhone 15 Pro Max 256GB", "Apple iPhone 15 Pro Max，钛金属设计，A17 Pro芯片，强大的相机系统", 9999.00, 10999.00, 100, 0, "https://picsum.photos/400/400?random=1"},
		{"MacBook Pro 14英寸", "Apple M3 Pro芯片，18GB内存，512GB存储，Liquid Retina XDR显示屏", 16999.00, 18999.00, 50, 1, "https://picsum.photos/400/400?random=2"},
		{"索尼 WH-1000XM5 耳机", "行业领先降噪，30小时续航，高解析度音频，舒适佩戴", 2699.00, 2999.00, 200, 0, "https://picsum.photos/400/400?random=3"},
		{"戴森 V15 Detect 吸尘器", "激光探测灰尘，LCD屏幕显示，60分钟续航，整机过滤", 5999.00, 6499.00, 80, 2, "https://picsum.photos/400/400?random=4"},
		{"Nike Air Jordan 1 Retro", "经典复刻，优质皮革，舒适缓震，时尚百搭", 1299.00, 1499.00, 150, 3, "https://picsum.photos/400/400?random=5"},
		{"SK-II 神仙水 230ml", "日本原装进口，改善肌肤纹理，提亮肤色，补水保湿", 1590.00, 1790.00, 300, 4, "https://picsum.photos/400/400?random=6"},
		{"三只松鼠坚果大礼包", "10袋装混合坚果，新鲜美味，送礼自用皆宜", 168.00, 199.00, 500, 5, "https://picsum.photos/400/400?random=7"},
		{"Kindle Paperwhite 电子书", "6.8英寸显示屏，可调节暖光，防水设计，海量书籍", 1068.00, 1199.00, 120, 6, "https://picsum.photos/400/400?random=8"},
		{"小米智能手表 S3", "1.43英寸AMOLED屏，150+运动模式，血氧检测，15天续航", 899.00, 999.00, 200, 0, "https://picsum.photos/400/400?random=9"},
		{"联想 ThinkPad X1 Carbon", "14英寸2.8K OLED屏，英特尔酷睿Ultra，32GB内存", 12999.00, 14999.00, 40, 1, "https://picsum.photos/400/400?random=10"},
		{"海尔 对开门冰箱 520L", "一级能效，变频风冷无霜，智能杀菌，大容量", 3999.00, 4599.00, 60, 2, "https://picsum.photos/400/400?random=11"},
		{"优衣库 男士羽绒服", "高保暖轻量设计，防泼水面料，多色可选", 599.00, 799.00, 300, 3, "https://picsum.photos/400/400?random=12"},
	}

	for _, p := range products {
		config.DB.Exec(`
			INSERT INTO products (seller_id, category_id, name, description, price, original_price, stock, images, status, rating, sales)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?)
		`, sellerID, categoryIDs[p.categoryIndex], p.name, p.description, p.price, p.originalPrice, p.stock, p.images, 4.5+float64(p.categoryIndex%5)*0.1, p.stock/2)
	}

	log.Println("示例数据初始化完成")
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// 引入bcrypt包

