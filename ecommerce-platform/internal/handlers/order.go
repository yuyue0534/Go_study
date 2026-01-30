package handlers

import (
	"encoding/json"
	"ecommerce-platform/internal/middleware"
	"ecommerce-platform/internal/models"
	"ecommerce-platform/internal/utils"
	"net/http"
	"strconv"
)

// CreateOrder 创建订单
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var req struct {
		AddressID int64  `json:"address_id"`
		Remark    string `json:"remark"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	// 获取收货地址
	var address *models.Address
	var err error
	if req.AddressID > 0 {
		address, err = models.GetAddressByID(req.AddressID)
	} else {
		address, err = models.GetDefaultAddress(session.UserID)
	}
	if err != nil || address == nil {
		utils.BadRequest(w, "请选择收货地址")
		return
	}

	// 获取选中的购物车商品
	cartItems, err := models.GetSelectedCartItems(session.UserID)
	if err != nil || len(cartItems) == 0 {
		utils.BadRequest(w, "请选择要购买的商品")
		return
	}

	// 按商家分组
	sellerOrders := make(map[int64][]*models.CartItem)
	for _, item := range cartItems {
		sellerOrders[item.Product.SellerID] = append(sellerOrders[item.Product.SellerID], item)
	}

	var orderIDs []int64

	// 为每个商家创建订单
	for sellerID, items := range sellerOrders {
		var totalAmount float64
		var orderItems []*models.OrderItem

		for _, item := range items {
			// 检查库存
			product, _ := models.GetProductByID(item.ProductID)
			if product == nil || product.Stock < item.Quantity {
				utils.BadRequest(w, "商品"+item.Product.Name+"库存不足")
				return
			}

			itemTotal := item.Product.Price * float64(item.Quantity)
			totalAmount += itemTotal

			orderItems = append(orderItems, &models.OrderItem{
				ProductID:    item.ProductID,
				ProductName:  item.Product.Name,
				ProductImage: utils.GetFirstImage(item.Product.Images),
				Price:        item.Product.Price,
				Quantity:     item.Quantity,
				TotalPrice:   itemTotal,
			})
		}

		receiverAddress := address.Province + address.City + address.District + address.DetailAddress

		order := &models.Order{
			OrderNo:         models.GenerateOrderNo(),
			UserID:          session.UserID,
			SellerID:        sellerID,
			TotalAmount:     totalAmount,
			PayAmount:       totalAmount,
			FreightAmount:   0,
			Status:          0, // 待支付
			ReceiverName:    address.ReceiverName,
			ReceiverPhone:   address.Phone,
			ReceiverAddress: receiverAddress,
			Remark:          req.Remark,
			Items:           orderItems,
		}

		orderID, err := models.CreateOrder(order)
		if err != nil {
			utils.InternalError(w, "创建订单失败: "+err.Error())
			return
		}
		orderIDs = append(orderIDs, orderID)
	}

	// 清空已购买的购物车商品
	models.ClearSelectedCart(session.UserID)

	utils.Success(w, map[string]interface{}{
		"order_ids": orderIDs,
	})
}

// GetOrders 获取订单列表
func GetOrders(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	query := r.URL.Query()
	page := utils.ParseInt(query.Get("page"), 1)
	size := utils.ParseInt(query.Get("size"), 10)
	status := utils.ParseInt(query.Get("status"), -1)

	orders, total, err := models.GetOrders(page, size, session.UserID, 0, status)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.SuccessPage(w, orders, total, page, size)
}

// GetOrder 获取订单详情
func GetOrder(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	order, err := models.GetOrderByID(id)
	if err != nil {
		utils.NotFound(w, "订单不存在")
		return
	}

	// 验证订单归属
	if session.Role != "admin" && order.UserID != session.UserID {
		// 检查是否是商家查看自己的订单
		seller, _ := models.GetSellerByUserID(session.UserID)
		if seller == nil || order.SellerID != seller.ID {
			utils.Forbidden(w, "没有权限")
			return
		}
	}

	utils.Success(w, order)
}

// PayOrder 支付订单
func PayOrder(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var req struct {
		OrderID int64  `json:"order_id"`
		PayType string `json:"pay_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	order, err := models.GetOrderByID(req.OrderID)
	if err != nil {
		utils.NotFound(w, "订单不存在")
		return
	}

	if order.UserID != session.UserID {
		utils.Forbidden(w, "没有权限")
		return
	}

	if order.Status != 0 {
		utils.BadRequest(w, "订单状态不允许支付")
		return
	}

	// 模拟支付成功
	if err := models.PayOrder(req.OrderID, req.PayType); err != nil {
		utils.InternalError(w, "支付失败")
		return
	}

	utils.SuccessMessage(w, "支付成功")
}

// CancelOrder 取消订单
func CancelOrder(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	order, err := models.GetOrderByID(id)
	if err != nil {
		utils.NotFound(w, "订单不存在")
		return
	}

	if order.UserID != session.UserID && session.Role != "admin" {
		utils.Forbidden(w, "没有权限")
		return
	}

	if err := models.CancelOrder(id); err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.SuccessMessage(w, "取消成功")
}

// ReceiveOrder 确认收货
func ReceiveOrder(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	order, err := models.GetOrderByID(id)
	if err != nil {
		utils.NotFound(w, "订单不存在")
		return
	}

	if order.UserID != session.UserID {
		utils.Forbidden(w, "没有权限")
		return
	}

	if order.Status != 2 {
		utils.BadRequest(w, "订单状态不允许确认收货")
		return
	}

	if err := models.ReceiveOrder(id); err != nil {
		utils.InternalError(w, "操作失败")
		return
	}

	utils.SuccessMessage(w, "确认收货成功")
}

// GetSellerOrders 商家获取订单列表
func GetSellerOrders(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	seller, err := models.GetSellerByUserID(session.UserID)
	if err != nil {
		utils.Forbidden(w, "您还不是商家")
		return
	}

	query := r.URL.Query()
	page := utils.ParseInt(query.Get("page"), 1)
	size := utils.ParseInt(query.Get("size"), 10)
	status := utils.ParseInt(query.Get("status"), -1)

	orders, total, err := models.GetOrders(page, size, 0, seller.ID, status)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.SuccessPage(w, orders, total, page, size)
}

// ShipOrder 商家发货
func ShipOrder(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var req struct {
		OrderID    int64  `json:"order_id"`
		TrackingNo string `json:"tracking_no"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	order, err := models.GetOrderByID(req.OrderID)
	if err != nil {
		utils.NotFound(w, "订单不存在")
		return
	}

	seller, _ := models.GetSellerByUserID(session.UserID)
	if session.Role != "admin" && (seller == nil || order.SellerID != seller.ID) {
		utils.Forbidden(w, "没有权限")
		return
	}

	if order.Status != 1 {
		utils.BadRequest(w, "订单状态不允许发货")
		return
	}

	if err := models.ShipOrder(req.OrderID, req.TrackingNo); err != nil {
		utils.InternalError(w, "发货失败")
		return
	}

	utils.SuccessMessage(w, "发货成功")
}

// GetSellerOrderStats 商家订单统计
func GetSellerOrderStats(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	seller, err := models.GetSellerByUserID(session.UserID)
	if err != nil {
		utils.Forbidden(w, "您还不是商家")
		return
	}

	stats, err := models.GetOrderStats(seller.ID)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.Success(w, stats)
}

// AdminGetOrders 管理员获取所有订单
func AdminGetOrders(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := utils.ParseInt(query.Get("page"), 1)
	size := utils.ParseInt(query.Get("size"), 20)
	status := utils.ParseInt(query.Get("status"), -1)

	orders, total, err := models.GetOrders(page, size, 0, 0, status)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.SuccessPage(w, orders, total, page, size)
}

// AdminGetOrderStats 管理员订单统计
func AdminGetOrderStats(w http.ResponseWriter, r *http.Request) {
	stats, err := models.GetOrderStats(0)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.Success(w, stats)
}
