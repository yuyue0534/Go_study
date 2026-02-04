package controller

import (
	"crm-erp-system/model"
	"crm-erp-system/service"
	"crm-erp-system/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type InventoryController struct {
	inventoryService *service.InventoryService
}

func NewInventoryController() *InventoryController {
	return &InventoryController{
		inventoryService: &service.InventoryService{},
	}
}

// Create 创建库存记录
func (ctrl *InventoryController) Create(c *gin.Context) {
	var inventory model.Inventory
	if err := c.ShouldBindJSON(&inventory); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	id, err := ctrl.inventoryService.Create(&inventory)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, gin.H{"id": id})
}

// GetByProductID 根据产品ID获取库存
func (ctrl *InventoryController) GetByProductID(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的产品ID")
		return
	}

	inventory, err := ctrl.inventoryService.GetByProductID(productID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, inventory)
}

// Update 更新库存
func (ctrl *InventoryController) Update(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的产品ID")
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,gte=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.inventoryService.Update(productID, req.Quantity); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// List 获取所有库存列表
func (ctrl *InventoryController) List(c *gin.Context) {
	inventories, err := ctrl.inventoryService.List()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list": inventories,
	})
}
