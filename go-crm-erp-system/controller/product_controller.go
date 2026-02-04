package controller

import (
	"crm-erp-system/model"
	"crm-erp-system/service"
	"crm-erp-system/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ProductController struct {
	productService *service.ProductService
}

func NewProductController() *ProductController {
	return &ProductController{
		productService: &service.ProductService{},
	}
}

// Create 创建产品
func (ctrl *ProductController) Create(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	id, err := ctrl.productService.Create(&product)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, gin.H{"id": id})
}

// Get 获取产品详情
func (ctrl *ProductController) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的产品ID")
		return
	}

	product, err := ctrl.productService.GetByID(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, product)
}

// List 获取产品列表
func (ctrl *ProductController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	products, err := ctrl.productService.List(page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":      products,
		"page":      page,
		"page_size": pageSize,
	})
}

// Update 更新产品
func (ctrl *ProductController) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的产品ID")
		return
	}

	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.productService.Update(id, &product); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除产品
func (ctrl *ProductController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "无效的产品ID")
		return
	}

	if err := ctrl.productService.Delete(id); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
