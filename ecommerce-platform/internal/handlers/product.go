package handlers

import (
	"encoding/json"
	"ecommerce-platform/internal/middleware"
	"ecommerce-platform/internal/models"
	"ecommerce-platform/internal/utils"
	"net/http"
	"strconv"
	"strings"
)

// GetCategories 获取所有分类
func GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := models.GetAllCategories()
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}
	utils.Success(w, categories)
}

// CreateCategory 创建分类（管理员）
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if category.Name == "" {
		utils.BadRequest(w, "分类名称不能为空")
		return
	}

	id, err := models.CreateCategory(&category)
	if err != nil {
		utils.InternalError(w, "创建失败")
		return
	}
	category.ID = id

	utils.Success(w, category)
}

// UpdateCategory 更新分类（管理员）
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if err := models.UpdateCategory(&category); err != nil {
		utils.InternalError(w, "更新失败")
		return
	}

	utils.SuccessMessage(w, "更新成功")
}

// DeleteCategory 删除分类（管理员）
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := models.DeleteCategory(id); err != nil {
		utils.InternalError(w, "删除失败")
		return
	}

	utils.SuccessMessage(w, "删除成功")
}

// GetProducts 获取商品列表
func GetProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := utils.ParseInt(query.Get("page"), 1)
	size := utils.ParseInt(query.Get("size"), 20)
	categoryID := utils.ParseInt64(query.Get("category_id"), 0)
	sellerID := utils.ParseInt64(query.Get("seller_id"), 0)
	keyword := query.Get("keyword")
	sortBy := query.Get("sort")
	status := utils.ParseInt(query.Get("status"), 1) // 默认只显示上架商品

	products, total, err := models.GetProducts(page, size, categoryID, sellerID, keyword, status, sortBy)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.SuccessPage(w, products, total, page, size)
}

// GetProduct 获取商品详情
func GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	product, err := models.GetProductByID(id)
	if err != nil {
		utils.NotFound(w, "商品不存在")
		return
	}

	utils.Success(w, product)
}

// GetHotProducts 获取热销商品
func GetHotProducts(w http.ResponseWriter, r *http.Request) {
	limit := utils.ParseInt(r.URL.Query().Get("limit"), 10)
	products, err := models.GetHotProducts(limit)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}
	utils.Success(w, products)
}

// GetNewProducts 获取新品
func GetNewProducts(w http.ResponseWriter, r *http.Request) {
	limit := utils.ParseInt(r.URL.Query().Get("limit"), 10)
	products, err := models.GetNewProducts(limit)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}
	utils.Success(w, products)
}

// CreateProduct 创建商品（商家）
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	seller, err := models.GetSellerByUserID(session.UserID)
	if err != nil {
		utils.Forbidden(w, "您还不是商家")
		return
	}

	if seller.Status != 1 {
		utils.Forbidden(w, "商家账号未审核通过")
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if product.Name == "" {
		utils.BadRequest(w, "商品名称不能为空")
		return
	}

	product.SellerID = seller.ID
	product.Status = 0 // 待审核

	id, err := models.CreateProduct(&product)
	if err != nil {
		utils.InternalError(w, "创建失败: "+err.Error())
		return
	}
	product.ID = id

	utils.Success(w, product)
}

// UpdateProduct 更新商品（商家）
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	// 验证商品归属
	existProduct, err := models.GetProductByID(product.ID)
	if err != nil {
		utils.NotFound(w, "商品不存在")
		return
	}

	seller, _ := models.GetSellerByUserID(session.UserID)
	if session.Role != "admin" && (seller == nil || existProduct.SellerID != seller.ID) {
		utils.Forbidden(w, "没有权限")
		return
	}

	// 保留原有的seller_id
	product.SellerID = existProduct.SellerID

	if err := models.UpdateProduct(&product); err != nil {
		utils.InternalError(w, "更新失败")
		return
	}

	utils.SuccessMessage(w, "更新成功")
}

// DeleteProduct 删除商品（商家）
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	product, err := models.GetProductByID(id)
	if err != nil {
		utils.NotFound(w, "商品不存在")
		return
	}

	seller, _ := models.GetSellerByUserID(session.UserID)
	if session.Role != "admin" && (seller == nil || product.SellerID != seller.ID) {
		utils.Forbidden(w, "没有权限")
		return
	}

	if err := models.DeleteProduct(id); err != nil {
		utils.InternalError(w, "删除失败")
		return
	}

	utils.SuccessMessage(w, "删除成功")
}

// GetSellerProducts 获取商家的商品列表
func GetSellerProducts(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	seller, err := models.GetSellerByUserID(session.UserID)
	if err != nil {
		utils.Forbidden(w, "您还不是商家")
		return
	}

	query := r.URL.Query()
	page := utils.ParseInt(query.Get("page"), 1)
	size := utils.ParseInt(query.Get("size"), 20)
	status := utils.ParseInt(query.Get("status"), -1)
	keyword := query.Get("keyword")

	products, total, err := models.GetProducts(page, size, 0, seller.ID, keyword, status, "")
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.SuccessPage(w, products, total, page, size)
}

// UpdateProductStock 更新库存（商家）
func UpdateProductStock(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var req struct {
		ProductID int64 `json:"product_id"`
		Quantity  int   `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	product, err := models.GetProductByID(req.ProductID)
	if err != nil {
		utils.NotFound(w, "商品不存在")
		return
	}

	seller, _ := models.GetSellerByUserID(session.UserID)
	if session.Role != "admin" && (seller == nil || product.SellerID != seller.ID) {
		utils.Forbidden(w, "没有权限")
		return
	}

	if err := models.UpdateProductStock(req.ProductID, req.Quantity); err != nil {
		utils.InternalError(w, "更新失败")
		return
	}

	utils.SuccessMessage(w, "更新成功")
}

// AdminGetProducts 管理员获取所有商品
func AdminGetProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := utils.ParseInt(query.Get("page"), 1)
	size := utils.ParseInt(query.Get("size"), 20)
	status := utils.ParseInt(query.Get("status"), -1)
	keyword := query.Get("keyword")

	products, total, err := models.GetProducts(page, size, 0, 0, keyword, status, "")
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.SuccessPage(w, products, total, page, size)
}

// AdminUpdateProductStatus 管理员更新商品状态
func AdminUpdateProductStatus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProductID int64 `json:"product_id"`
		Status    int   `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if err := models.UpdateProductStatus(req.ProductID, req.Status); err != nil {
		utils.InternalError(w, "更新失败")
		return
	}

	utils.SuccessMessage(w, "更新成功")
}

// UploadImage 上传图片
func UploadImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.BadRequest(w, "文件上传失败")
		return
	}
	defer file.Close()

	// 检查文件类型
	contentType := handler.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		utils.BadRequest(w, "只能上传图片文件")
		return
	}

	// 生成文件名
	ext := ".jpg"
	if strings.Contains(handler.Filename, ".") {
		parts := strings.Split(handler.Filename, ".")
		ext = "." + parts[len(parts)-1]
	}
	filename := utils.GenerateToken()[:16] + ext

	// 保存文件
	// 这里简化处理，实际应该保存到文件服务器或云存储
	// 返回一个模拟的URL
	imageURL := "/uploads/" + filename

	utils.Success(w, map[string]string{
		"url": imageURL,
	})
}
