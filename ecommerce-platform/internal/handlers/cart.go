package handlers

import (
	"encoding/json"
	"ecommerce-platform/internal/middleware"
	"ecommerce-platform/internal/models"
	"ecommerce-platform/internal/utils"
	"net/http"
	"strconv"
)

// GetCart 获取购物车
func GetCart(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	items, err := models.GetCartItems(session.UserID)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	// 计算总价
	var totalPrice float64
	var selectedCount int
	for _, item := range items {
		if item.Selected == 1 {
			totalPrice += item.Product.Price * float64(item.Quantity)
			selectedCount += item.Quantity
		}
	}

	utils.Success(w, map[string]interface{}{
		"items":          items,
		"total_price":    totalPrice,
		"selected_count": selectedCount,
	})
}

// AddToCart 添加商品到购物车
func AddToCart(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var req struct {
		ProductID int64 `json:"product_id"`
		Quantity  int   `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	// 检查商品是否存在
	product, err := models.GetProductByID(req.ProductID)
	if err != nil {
		utils.NotFound(w, "商品不存在")
		return
	}

	if product.Status != 1 {
		utils.BadRequest(w, "商品已下架")
		return
	}

	if product.Stock < req.Quantity {
		utils.BadRequest(w, "库存不足")
		return
	}

	if err := models.AddToCart(session.UserID, req.ProductID, req.Quantity); err != nil {
		utils.InternalError(w, "添加失败")
		return
	}

	// 获取购物车数量
	count, _ := models.GetCartCount(session.UserID)

	utils.Success(w, map[string]interface{}{
		"cart_count": count,
	})
}

// UpdateCartItem 更新购物车商品
func UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID       int64 `json:"id"`
		Quantity int   `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if req.Quantity <= 0 {
		utils.BadRequest(w, "数量必须大于0")
		return
	}

	if err := models.UpdateCartItemQuantity(req.ID, req.Quantity); err != nil {
		utils.InternalError(w, "更新失败")
		return
	}

	utils.SuccessMessage(w, "更新成功")
}

// UpdateCartItemSelected 更新购物车选中状态
func UpdateCartItemSelected(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID       int64 `json:"id"`
		Selected int   `json:"selected"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if err := models.UpdateCartItemSelected(req.ID, req.Selected); err != nil {
		utils.InternalError(w, "更新失败")
		return
	}

	utils.SuccessMessage(w, "更新成功")
}

// SelectAllCart 全选/取消全选购物车
func SelectAllCart(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var req struct {
		Selected int `json:"selected"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if err := models.UpdateCartAllSelected(session.UserID, req.Selected); err != nil {
		utils.InternalError(w, "更新失败")
		return
	}

	utils.SuccessMessage(w, "更新成功")
}

// DeleteCartItem 删除购物车商品
func DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := models.DeleteCartItem(id); err != nil {
		utils.InternalError(w, "删除失败")
		return
	}

	utils.SuccessMessage(w, "删除成功")
}

// ClearCart 清空购物车
func ClearCart(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	if err := models.ClearCart(session.UserID); err != nil {
		utils.InternalError(w, "清空失败")
		return
	}

	utils.SuccessMessage(w, "清空成功")
}

// GetCartCount 获取购物车数量
func GetCartCount(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	count, err := models.GetCartCount(session.UserID)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.Success(w, map[string]int{
		"count": count,
	})
}
