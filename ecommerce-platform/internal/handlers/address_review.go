package handlers

import (
	"encoding/json"
	"ecommerce-platform/internal/middleware"
	"ecommerce-platform/internal/models"
	"ecommerce-platform/internal/utils"
	"net/http"
	"strconv"
)

// GetAddresses 获取收货地址列表
func GetAddresses(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	addresses, err := models.GetUserAddresses(session.UserID)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.Success(w, addresses)
}

// CreateAddress 创建收货地址
func CreateAddress(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var address models.Address
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if address.ReceiverName == "" || address.Phone == "" || address.DetailAddress == "" {
		utils.BadRequest(w, "收货人、电话和详细地址不能为空")
		return
	}

	address.UserID = session.UserID

	id, err := models.CreateAddress(&address)
	if err != nil {
		utils.InternalError(w, "创建失败")
		return
	}
	address.ID = id

	utils.Success(w, address)
}

// UpdateAddress 更新收货地址
func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var address models.Address
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	// 验证地址归属
	existAddr, err := models.GetAddressByID(address.ID)
	if err != nil || existAddr.UserID != session.UserID {
		utils.Forbidden(w, "没有权限")
		return
	}

	address.UserID = session.UserID

	if err := models.UpdateAddress(&address); err != nil {
		utils.InternalError(w, "更新失败")
		return
	}

	utils.SuccessMessage(w, "更新成功")
}

// DeleteAddress 删除收货地址
func DeleteAddress(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	// 验证地址归属
	address, err := models.GetAddressByID(id)
	if err != nil || address.UserID != session.UserID {
		utils.Forbidden(w, "没有权限")
		return
	}

	if err := models.DeleteAddress(id); err != nil {
		utils.InternalError(w, "删除失败")
		return
	}

	utils.SuccessMessage(w, "删除成功")
}

// SetDefaultAddress 设置默认地址
func SetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var req struct {
		AddressID int64 `json:"address_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if err := models.SetDefaultAddress(session.UserID, req.AddressID); err != nil {
		utils.InternalError(w, "设置失败")
		return
	}

	utils.SuccessMessage(w, "设置成功")
}

// GetProductReviews 获取商品评价
func GetProductReviews(w http.ResponseWriter, r *http.Request) {
	productID := utils.ParseInt64(r.URL.Query().Get("product_id"), 0)
	page := utils.ParseInt(r.URL.Query().Get("page"), 1)
	size := utils.ParseInt(r.URL.Query().Get("size"), 10)

	reviews, total, err := models.GetProductReviews(productID, page, size)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.SuccessPage(w, reviews, total, page, size)
}

// CreateReview 创建评价
func CreateReview(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if review.Rating < 1 || review.Rating > 5 {
		utils.BadRequest(w, "评分必须在1-5之间")
		return
	}

	// 检查订单是否存在且已完成
	order, err := models.GetOrderByID(review.OrderID)
	if err != nil || order.UserID != session.UserID {
		utils.BadRequest(w, "订单不存在")
		return
	}

	if order.Status != 3 {
		utils.BadRequest(w, "订单未完成，不能评价")
		return
	}

	// 检查是否已评价
	if models.CheckReviewExists(review.OrderID, review.ProductID) {
		utils.BadRequest(w, "已评价过该商品")
		return
	}

	review.UserID = session.UserID
	review.Status = 1

	id, err := models.CreateReview(&review)
	if err != nil {
		utils.InternalError(w, "评价失败")
		return
	}
	review.ID = id

	utils.Success(w, review)
}

// GetUserReviews 获取用户的评价
func GetUserReviews(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	page := utils.ParseInt(r.URL.Query().Get("page"), 1)
	size := utils.ParseInt(r.URL.Query().Get("size"), 10)

	reviews, total, err := models.GetUserReviews(session.UserID, page, size)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.SuccessPage(w, reviews, total, page, size)
}

// GetSellerReviews 商家获取评价列表
func GetSellerReviews(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	seller, err := models.GetSellerByUserID(session.UserID)
	if err != nil {
		utils.Forbidden(w, "您还不是商家")
		return
	}

	page := utils.ParseInt(r.URL.Query().Get("page"), 1)
	size := utils.ParseInt(r.URL.Query().Get("size"), 10)

	reviews, total, err := models.GetSellerReviews(seller.ID, page, size)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.SuccessPage(w, reviews, total, page, size)
}

// ReplyReview 商家回复评价
func ReplyReview(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var req struct {
		ReviewID int64  `json:"review_id"`
		Reply    string `json:"reply"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	review, err := models.GetReviewByID(req.ReviewID)
	if err != nil {
		utils.NotFound(w, "评价不存在")
		return
	}

	// 验证商品归属
	product, _ := models.GetProductByID(review.ProductID)
	seller, _ := models.GetSellerByUserID(session.UserID)
	if session.Role != "admin" && (seller == nil || product.SellerID != seller.ID) {
		utils.Forbidden(w, "没有权限")
		return
	}

	if err := models.ReplyReview(req.ReviewID, req.Reply); err != nil {
		utils.InternalError(w, "回复失败")
		return
	}

	utils.SuccessMessage(w, "回复成功")
}
