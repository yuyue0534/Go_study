package controller

import (
	"crm-erp-system/model"
	"crm-erp-system/service"
	"crm-erp-system/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type OrderController struct {
	orderService *service.OrderService
}

func NewOrderController() *OrderController {
	return &OrderController{
		orderService: &service.OrderService{},
	}
}

// Create 创建订单
func (ctrl *OrderController) Create(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	id, err := ctrl.orderService.Create(&order, userID.(int64))
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, gin.H{"id": id})
}

// Get 获取订单详情
func (ctrl *OrderController) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的订单ID")
		return
	}

	order, err := ctrl.orderService.GetByID(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, order)
}

// List 获取订单列表
func (ctrl *OrderController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	orders, err := ctrl.orderService.List(page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":      orders,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateStatus 更新订单状态
func (ctrl *OrderController) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的订单ID")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.orderService.UpdateStatus(id, req.Status); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除订单
func (ctrl *OrderController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的订单ID")
		return
	}

	if err := ctrl.orderService.Delete(id); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
