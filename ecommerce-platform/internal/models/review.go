package models

import (
	"database/sql"
	"ecommerce-platform/internal/config"
	"time"
)

type Review struct {
	ID        int64      `json:"id"`
	OrderID   int64      `json:"order_id"`
	ProductID int64      `json:"product_id"`
	UserID    int64      `json:"user_id"`
	Rating    int        `json:"rating"`
	Content   string     `json:"content"`
	Images    string     `json:"images"`
	Reply     string     `json:"reply"`
	ReplyTime *time.Time `json:"reply_time"`
	Status    int        `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	// 关联信息
	Username    string `json:"username,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	ProductName string `json:"product_name,omitempty"`
}

// CreateReview 创建评价
func CreateReview(review *Review) (int64, error) {
	result, err := config.DB.Exec(`
		INSERT INTO reviews (order_id, product_id, user_id, rating, content, images, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, review.OrderID, review.ProductID, review.UserID, review.Rating, review.Content, review.Images, review.Status)
	if err != nil {
		return 0, err
	}

	// 更新商品评分
	UpdateProductRating(review.ProductID)

	return result.LastInsertId()
}

// GetReviewByID 通过ID获取评价
func GetReviewByID(id int64) (*Review, error) {
	review := &Review{}
	var replyTime sql.NullTime
	err := config.DB.QueryRow(`
		SELECT r.id, r.order_id, r.product_id, r.user_id, r.rating, r.content, r.images,
			r.reply, r.reply_time, r.status, r.created_at, u.username, u.avatar
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.id = ?
	`, id).Scan(&review.ID, &review.OrderID, &review.ProductID, &review.UserID, &review.Rating,
		&review.Content, &review.Images, &review.Reply, &replyTime, &review.Status,
		&review.CreatedAt, &review.Username, &review.Avatar)
	if err != nil {
		return nil, err
	}
	if replyTime.Valid {
		review.ReplyTime = &replyTime.Time
	}
	return review, nil
}

// GetProductReviews 获取商品评价列表
func GetProductReviews(productID int64, page, pageSize int) ([]*Review, int, error) {
	offset := (page - 1) * pageSize

	var total int
	err := config.DB.QueryRow(`SELECT COUNT(*) FROM reviews WHERE product_id = ? AND status = 1`, productID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := config.DB.Query(`
		SELECT r.id, r.order_id, r.product_id, r.user_id, r.rating, r.content, r.images,
			r.reply, r.reply_time, r.status, r.created_at, u.username, u.avatar
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.product_id = ? AND r.status = 1
		ORDER BY r.created_at DESC LIMIT ? OFFSET ?
	`, productID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []*Review
	for rows.Next() {
		review := &Review{}
		var replyTime sql.NullTime
		err := rows.Scan(&review.ID, &review.OrderID, &review.ProductID, &review.UserID, &review.Rating,
			&review.Content, &review.Images, &review.Reply, &replyTime, &review.Status,
			&review.CreatedAt, &review.Username, &review.Avatar)
		if err != nil {
			return nil, 0, err
		}
		if replyTime.Valid {
			review.ReplyTime = &replyTime.Time
		}
		reviews = append(reviews, review)
	}
	return reviews, total, nil
}

// GetUserReviews 获取用户评价列表
func GetUserReviews(userID int64, page, pageSize int) ([]*Review, int, error) {
	offset := (page - 1) * pageSize

	var total int
	err := config.DB.QueryRow(`SELECT COUNT(*) FROM reviews WHERE user_id = ?`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := config.DB.Query(`
		SELECT r.id, r.order_id, r.product_id, r.user_id, r.rating, r.content, r.images,
			r.reply, r.reply_time, r.status, r.created_at, p.name
		FROM reviews r
		LEFT JOIN products p ON r.product_id = p.id
		WHERE r.user_id = ?
		ORDER BY r.created_at DESC LIMIT ? OFFSET ?
	`, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []*Review
	for rows.Next() {
		review := &Review{}
		var replyTime sql.NullTime
		err := rows.Scan(&review.ID, &review.OrderID, &review.ProductID, &review.UserID, &review.Rating,
			&review.Content, &review.Images, &review.Reply, &replyTime, &review.Status,
			&review.CreatedAt, &review.ProductName)
		if err != nil {
			return nil, 0, err
		}
		if replyTime.Valid {
			review.ReplyTime = &replyTime.Time
		}
		reviews = append(reviews, review)
	}
	return reviews, total, nil
}

// GetSellerReviews 获取商家评价列表
func GetSellerReviews(sellerID int64, page, pageSize int) ([]*Review, int, error) {
	offset := (page - 1) * pageSize

	var total int
	err := config.DB.QueryRow(`
		SELECT COUNT(*) FROM reviews r
		JOIN products p ON r.product_id = p.id
		WHERE p.seller_id = ?
	`, sellerID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := config.DB.Query(`
		SELECT r.id, r.order_id, r.product_id, r.user_id, r.rating, r.content, r.images,
			r.reply, r.reply_time, r.status, r.created_at, u.username, u.avatar, p.name
		FROM reviews r
		JOIN products p ON r.product_id = p.id
		LEFT JOIN users u ON r.user_id = u.id
		WHERE p.seller_id = ?
		ORDER BY r.created_at DESC LIMIT ? OFFSET ?
	`, sellerID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []*Review
	for rows.Next() {
		review := &Review{}
		var replyTime sql.NullTime
		err := rows.Scan(&review.ID, &review.OrderID, &review.ProductID, &review.UserID, &review.Rating,
			&review.Content, &review.Images, &review.Reply, &replyTime, &review.Status,
			&review.CreatedAt, &review.Username, &review.Avatar, &review.ProductName)
		if err != nil {
			return nil, 0, err
		}
		if replyTime.Valid {
			review.ReplyTime = &replyTime.Time
		}
		reviews = append(reviews, review)
	}
	return reviews, total, nil
}

// ReplyReview 商家回复评价
func ReplyReview(id int64, reply string) error {
	_, err := config.DB.Exec(`
		UPDATE reviews SET reply = ?, reply_time = CURRENT_TIMESTAMP WHERE id = ?
	`, reply, id)
	return err
}

// UpdateReviewStatus 更新评价状态
func UpdateReviewStatus(id int64, status int) error {
	_, err := config.DB.Exec(`UPDATE reviews SET status = ? WHERE id = ?`, status, id)
	return err
}

// DeleteReview 删除评价
func DeleteReview(id int64) error {
	// 先获取商品ID
	var productID int64
	config.DB.QueryRow(`SELECT product_id FROM reviews WHERE id = ?`, id).Scan(&productID)

	_, err := config.DB.Exec(`DELETE FROM reviews WHERE id = ?`, id)
	if err != nil {
		return err
	}

	// 更新商品评分
	if productID > 0 {
		UpdateProductRating(productID)
	}
	return nil
}

// CheckReviewExists 检查是否已评价
func CheckReviewExists(orderID, productID int64) bool {
	var count int
	config.DB.QueryRow(`SELECT COUNT(*) FROM reviews WHERE order_id = ? AND product_id = ?`, orderID, productID).Scan(&count)
	return count > 0
}
