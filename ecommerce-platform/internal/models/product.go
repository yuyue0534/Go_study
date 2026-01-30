package models

import (
	"ecommerce-platform/internal/config"
	"strings"
	"time"
)

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  int64     `json:"parent_id"`
	Icon      string    `json:"icon"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID            int64     `json:"id"`
	SellerID      int64     `json:"seller_id"`
	CategoryID    int64     `json:"category_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	OriginalPrice float64   `json:"original_price"`
	Stock         int       `json:"stock"`
	Sales         int       `json:"sales"`
	Images        string    `json:"images"`
	Brand         string    `json:"brand"`
	Status        int       `json:"status"`
	Rating        float64   `json:"rating"`
	RatingCount   int       `json:"rating_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	// 关联信息
	SellerName   string `json:"seller_name,omitempty"`
	CategoryName string `json:"category_name,omitempty"`
}

// CreateCategory 创建分类
func CreateCategory(category *Category) (int64, error) {
	result, err := config.DB.Exec(`
		INSERT INTO categories (name, parent_id, icon, sort_order)
		VALUES (?, ?, ?, ?)
	`, category.Name, category.ParentID, category.Icon, category.SortOrder)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetAllCategories 获取所有分类
func GetAllCategories() ([]*Category, error) {
	rows, err := config.DB.Query(`
		SELECT id, name, parent_id, icon, sort_order, created_at
		FROM categories ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		cat := &Category{}
		err := rows.Scan(&cat.ID, &cat.Name, &cat.ParentID, &cat.Icon, &cat.SortOrder, &cat.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// GetCategoryByID 通过ID获取分类
func GetCategoryByID(id int64) (*Category, error) {
	cat := &Category{}
	err := config.DB.QueryRow(`
		SELECT id, name, parent_id, icon, sort_order, created_at
		FROM categories WHERE id = ?
	`, id).Scan(&cat.ID, &cat.Name, &cat.ParentID, &cat.Icon, &cat.SortOrder, &cat.CreatedAt)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

// UpdateCategory 更新分类
func UpdateCategory(category *Category) error {
	_, err := config.DB.Exec(`
		UPDATE categories SET name = ?, parent_id = ?, icon = ?, sort_order = ?
		WHERE id = ?
	`, category.Name, category.ParentID, category.Icon, category.SortOrder, category.ID)
	return err
}

// DeleteCategory 删除分类
func DeleteCategory(id int64) error {
	_, err := config.DB.Exec(`DELETE FROM categories WHERE id = ?`, id)
	return err
}

// CreateProduct 创建商品
func CreateProduct(product *Product) (int64, error) {
	result, err := config.DB.Exec(`
		INSERT INTO products (seller_id, category_id, name, description, price, original_price, stock, images, brand, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, product.SellerID, product.CategoryID, product.Name, product.Description,
		product.Price, product.OriginalPrice, product.Stock, product.Images, product.Brand, product.Status)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetProductByID 通过ID获取商品
func GetProductByID(id int64) (*Product, error) {
	product := &Product{}
	err := config.DB.QueryRow(`
		SELECT p.id, p.seller_id, COALESCE(p.category_id, 0), p.name, COALESCE(p.description, ''), 
			p.price, COALESCE(p.original_price, 0), p.stock, COALESCE(p.sales, 0), COALESCE(p.images, ''), 
			COALESCE(p.brand, ''), p.status, COALESCE(p.rating, 5.0), COALESCE(p.rating_count, 0),
			p.created_at, p.updated_at, COALESCE(s.shop_name, ''), COALESCE(c.name, '')
		FROM products p
		LEFT JOIN sellers s ON p.seller_id = s.id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = ?
	`, id).Scan(&product.ID, &product.SellerID, &product.CategoryID, &product.Name, &product.Description,
		&product.Price, &product.OriginalPrice, &product.Stock, &product.Sales, &product.Images,
		&product.Brand, &product.Status, &product.Rating, &product.RatingCount,
		&product.CreatedAt, &product.UpdatedAt, &product.SellerName, &product.CategoryName)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// UpdateProduct 更新商品
func UpdateProduct(product *Product) error {
	_, err := config.DB.Exec(`
		UPDATE products SET category_id = ?, name = ?, description = ?, price = ?,
			original_price = ?, stock = ?, images = ?, brand = ?, status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, product.CategoryID, product.Name, product.Description, product.Price,
		product.OriginalPrice, product.Stock, product.Images, product.Brand, product.Status, product.ID)
	return err
}

// DeleteProduct 删除商品
func DeleteProduct(id int64) error {
	_, err := config.DB.Exec(`DELETE FROM products WHERE id = ?`, id)
	return err
}

// UpdateProductStatus 更新商品状态
func UpdateProductStatus(id int64, status int) error {
	_, err := config.DB.Exec(`UPDATE products SET status = ? WHERE id = ?`, status, id)
	return err
}

// GetProducts 获取商品列表
func GetProducts(page, pageSize int, categoryID int64, sellerID int64, keyword string, status int, sortBy string) ([]*Product, int, error) {
	offset := (page - 1) * pageSize
	
	// 构建查询条件
	conditions := []string{"1=1"}
	args := []interface{}{}
	
	if categoryID > 0 {
		conditions = append(conditions, "p.category_id = ?")
		args = append(args, categoryID)
	}
	if sellerID > 0 {
		conditions = append(conditions, "p.seller_id = ?")
		args = append(args, sellerID)
	}
	if keyword != "" {
		conditions = append(conditions, "(p.name LIKE ? OR p.description LIKE ?)")
		args = append(args, "%"+keyword+"%", "%"+keyword+"%")
	}
	if status >= 0 {
		conditions = append(conditions, "p.status = ?")
		args = append(args, status)
	}
	
	whereClause := strings.Join(conditions, " AND ")
	
	// 统计总数
	var total int
	countSQL := `SELECT COUNT(*) FROM products p WHERE ` + whereClause
	err := config.DB.QueryRow(countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	// 排序
	orderClause := "ORDER BY p.created_at DESC"
	switch sortBy {
	case "price_asc":
		orderClause = "ORDER BY p.price ASC"
	case "price_desc":
		orderClause = "ORDER BY p.price DESC"
	case "sales":
		orderClause = "ORDER BY p.sales DESC"
	case "rating":
		orderClause = "ORDER BY p.rating DESC"
	}
	
	// 查询商品
	querySQL := `
		SELECT p.id, p.seller_id, COALESCE(p.category_id, 0), p.name, COALESCE(p.description, ''), 
			p.price, COALESCE(p.original_price, 0), p.stock, COALESCE(p.sales, 0), COALESCE(p.images, ''), 
			COALESCE(p.brand, ''), p.status, COALESCE(p.rating, 5.0), COALESCE(p.rating_count, 0),
			p.created_at, p.updated_at, COALESCE(s.shop_name, ''), COALESCE(c.name, '')
		FROM products p
		LEFT JOIN sellers s ON p.seller_id = s.id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE ` + whereClause + ` ` + orderClause + ` LIMIT ? OFFSET ?`
	
	args = append(args, pageSize, offset)
	rows, err := config.DB.Query(querySQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(&product.ID, &product.SellerID, &product.CategoryID, &product.Name, &product.Description,
			&product.Price, &product.OriginalPrice, &product.Stock, &product.Sales, &product.Images,
			&product.Brand, &product.Status, &product.Rating, &product.RatingCount,
			&product.CreatedAt, &product.UpdatedAt, &product.SellerName, &product.CategoryName)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}
	return products, total, nil
}

// UpdateProductStock 更新商品库存
func UpdateProductStock(id int64, quantity int) error {
	_, err := config.DB.Exec(`
		UPDATE products SET stock = stock + ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
	`, quantity, id)
	return err
}

// UpdateProductSales 更新商品销量
func UpdateProductSales(id int64, quantity int) error {
	_, err := config.DB.Exec(`
		UPDATE products SET sales = sales + ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
	`, quantity, id)
	return err
}

// UpdateProductRating 更新商品评分
func UpdateProductRating(productID int64) error {
	_, err := config.DB.Exec(`
		UPDATE products SET 
			rating = (SELECT COALESCE(AVG(rating), 5.0) FROM reviews WHERE product_id = ? AND status = 1),
			rating_count = (SELECT COUNT(*) FROM reviews WHERE product_id = ? AND status = 1)
		WHERE id = ?
	`, productID, productID, productID)
	return err
}

// GetHotProducts 获取热销商品
func GetHotProducts(limit int) ([]*Product, error) {
	rows, err := config.DB.Query(`
		SELECT p.id, p.seller_id, COALESCE(p.category_id, 0), p.name, COALESCE(p.description, ''), 
			p.price, COALESCE(p.original_price, 0), p.stock, p.sales, COALESCE(p.images, ''), 
			COALESCE(p.brand, ''), p.status, COALESCE(p.rating, 5.0), COALESCE(p.rating_count, 0),
			p.created_at, p.updated_at, COALESCE(s.shop_name, ''), COALESCE(c.name, '')
		FROM products p
		LEFT JOIN sellers s ON p.seller_id = s.id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.status = 1
		ORDER BY p.sales DESC LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(&product.ID, &product.SellerID, &product.CategoryID, &product.Name, &product.Description,
			&product.Price, &product.OriginalPrice, &product.Stock, &product.Sales, &product.Images,
			&product.Brand, &product.Status, &product.Rating, &product.RatingCount,
			&product.CreatedAt, &product.UpdatedAt, &product.SellerName, &product.CategoryName)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

// GetNewProducts 获取新品
func GetNewProducts(limit int) ([]*Product, error) {
	rows, err := config.DB.Query(`
		SELECT p.id, p.seller_id, COALESCE(p.category_id, 0), p.name, COALESCE(p.description, ''), 
			p.price, COALESCE(p.original_price, 0), p.stock, p.sales, COALESCE(p.images, ''), 
			COALESCE(p.brand, ''), p.status, COALESCE(p.rating, 5.0), COALESCE(p.rating_count, 0),
			p.created_at, p.updated_at, COALESCE(s.shop_name, ''), COALESCE(c.name, '')
		FROM products p
		LEFT JOIN sellers s ON p.seller_id = s.id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.status = 1
		ORDER BY p.created_at DESC LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(&product.ID, &product.SellerID, &product.CategoryID, &product.Name, &product.Description,
			&product.Price, &product.OriginalPrice, &product.Stock, &product.Sales, &product.Images,
			&product.Brand, &product.Status, &product.Rating, &product.RatingCount,
			&product.CreatedAt, &product.UpdatedAt, &product.SellerName, &product.CategoryName)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
