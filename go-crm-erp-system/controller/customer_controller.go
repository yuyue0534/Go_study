package controller

import (
	"crm-erp-system/model"
	"crm-erp-system/service"
	"crm-erp-system/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type CustomerController struct {
	customerService *service.CustomerService
}

func NewCustomerController() *CustomerController {
	return &CustomerController{
		customerService: &service.CustomerService{},
	}
}

// Create 创建客户
func (ctrl *CustomerController) Create(c *gin.Context) {
	var customer model.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	id, err := ctrl.customerService.Create(&customer, userID.(int64))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"id": id})
}

// Get 获取客户详情
func (ctrl *CustomerController) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的客户ID")
		return
	}

	customer, err := ctrl.customerService.GetByID(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, customer)
}

// List 获取客户列表
func (ctrl *CustomerController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	customers, err := ctrl.customerService.List(page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":      customers,
		"page":      page,
		"page_size": pageSize,
	})
}

// Update 更新客户
func (ctrl *CustomerController) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的客户ID")
		return
	}

	var customer model.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.customerService.Update(id, &customer); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除客户
func (ctrl *CustomerController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的客户ID")
		return
	}

	if err := ctrl.customerService.Delete(id); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
