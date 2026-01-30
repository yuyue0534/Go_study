package handlers

import (
	"encoding/json"
	"ecommerce-platform/internal/middleware"
	"ecommerce-platform/internal/models"
	"ecommerce-platform/internal/utils"
	"net/http"
)

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register 用户注册
func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if req.Username == "" || req.Password == "" || req.Email == "" {
		utils.BadRequest(w, "用户名、密码和邮箱不能为空")
		return
	}

	if _, err := models.GetUserByUsername(req.Username); err == nil {
		utils.BadRequest(w, "用户名已存在")
		return
	}

	if _, err := models.GetUserByEmail(req.Email); err == nil {
		utils.BadRequest(w, "邮箱已被注册")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalError(w, "密码加密失败")
		return
	}

	role := "customer"
	if req.Role == "seller" {
		role = "seller"
	}

	user := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     role,
		Status:   1,
	}

	userID, err := models.CreateUser(user)
	if err != nil {
		utils.InternalError(w, "注册失败: "+err.Error())
		return
	}
	user.ID = userID

	if role == "seller" {
		seller := &models.Seller{
			UserID:   userID,
			ShopName: req.Username + "的店铺",
			Status:   0,
		}
		models.CreateSeller(seller)
	}

	token := middleware.CreateSession(user)

	utils.Success(w, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// Login 用户登录
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.BadRequest(w, "用户名和密码不能为空")
		return
	}

	user, err := models.GetUserByUsername(req.Username)
	if err != nil {
		utils.BadRequest(w, "用户名或密码错误")
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		utils.BadRequest(w, "用户名或密码错误")
		return
	}

	if user.Status != 1 {
		utils.Forbidden(w, "账号已被禁用")
		return
	}

	token := middleware.CreateSession(user)

	// 设置Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	})

	var sellerInfo interface{}
	if user.Role == "seller" {
		seller, _ := models.GetSellerByUserID(user.ID)
		if seller != nil {
			sellerInfo = map[string]interface{}{
				"id":        seller.ID,
				"shop_name": seller.ShopName,
				"status":    seller.Status,
			}
		}
	}

	utils.Success(w, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"phone":    user.Phone,
			"role":     user.Role,
			"avatar":   user.Avatar,
		},
		"seller": sellerInfo,
	})
}

// Logout 用户登出
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err == nil {
		middleware.DeleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	utils.SuccessMessage(w, "退出成功")
}

// GetCurrentUser 获取当前用户信息
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)
	if session == nil {
		utils.Unauthorized(w, "请先登录")
		return
	}

	user, err := models.GetUserByID(session.UserID)
	if err != nil {
		utils.NotFound(w, "用户不存在")
		return
	}

	var sellerInfo interface{}
	if user.Role == "seller" {
		seller, _ := models.GetSellerByUserID(user.ID)
		if seller != nil {
			sellerInfo = map[string]interface{}{
				"id":          seller.ID,
				"shop_name":   seller.ShopName,
				"description": seller.ShopDescription,
				"logo":        seller.ShopLogo,
				"status":      seller.Status,
			}
		}
	}

	utils.Success(w, map[string]interface{}{
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"phone":    user.Phone,
			"role":     user.Role,
			"avatar":   user.Avatar,
		},
		"seller": sellerInfo,
	})
}

// UpdateProfile 更新用户资料
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var req struct {
		Email  string `json:"email"`
		Phone  string `json:"phone"`
		Avatar string `json:"avatar"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	user, _ := models.GetUserByID(session.UserID)
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := models.UpdateUser(user); err != nil {
		utils.InternalError(w, "更新失败")
		return
	}

	utils.SuccessMessage(w, "更新成功")
}

// ChangePassword 修改密码
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	user, _ := models.GetUserByID(session.UserID)
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		utils.BadRequest(w, "原密码错误")
		return
	}

	hashedPassword, _ := utils.HashPassword(req.NewPassword)
	if err := models.UpdateUserPassword(session.UserID, hashedPassword); err != nil {
		utils.InternalError(w, "修改失败")
		return
	}

	utils.SuccessMessage(w, "密码修改成功")
}

// UpdateSellerInfo 更新商家信息
func UpdateSellerInfo(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetCurrentSession(r)

	var req struct {
		ShopName        string `json:"shop_name"`
		ShopDescription string `json:"shop_description"`
		ShopLogo        string `json:"shop_logo"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	seller, err := models.GetSellerByUserID(session.UserID)
	if err != nil {
		utils.NotFound(w, "商家信息不存在")
		return
	}

	if req.ShopName != "" {
		seller.ShopName = req.ShopName
	}
	if req.ShopDescription != "" {
		seller.ShopDescription = req.ShopDescription
	}
	if req.ShopLogo != "" {
		seller.ShopLogo = req.ShopLogo
	}

	if err := models.UpdateSeller(seller); err != nil {
		utils.InternalError(w, "更新失败")
		return
	}

	utils.SuccessMessage(w, "更新成功")
}

// AdminGetUsers 管理员获取用户列表
func AdminGetUsers(w http.ResponseWriter, r *http.Request) {
	page := utils.ParseInt(r.URL.Query().Get("page"), 1)
	size := utils.ParseInt(r.URL.Query().Get("size"), 20)
	role := r.URL.Query().Get("role")

	users, total, err := models.GetAllUsers(page, size, role)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.SuccessPage(w, users, total, page, size)
}

// AdminUpdateUserStatus 管理员更新用户状态
func AdminUpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int64 `json:"user_id"`
		Status int   `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if err := models.UpdateUserStatus(req.UserID, req.Status); err != nil {
		utils.InternalError(w, "更新失败")
		return
	}

	utils.SuccessMessage(w, "更新成功")
}

// AdminGetSellers 管理员获取商家列表
func AdminGetSellers(w http.ResponseWriter, r *http.Request) {
	page := utils.ParseInt(r.URL.Query().Get("page"), 1)
	size := utils.ParseInt(r.URL.Query().Get("size"), 20)

	sellers, total, err := models.GetAllSellers(page, size)
	if err != nil {
		utils.InternalError(w, "获取失败")
		return
	}

	utils.SuccessPage(w, sellers, total, page, size)
}

// AdminUpdateSellerStatus 管理员更新商家状态
func AdminUpdateSellerStatus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SellerID int64 `json:"seller_id"`
		Status   int   `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "请求参数错误")
		return
	}

	if err := models.UpdateSellerStatus(req.SellerID, req.Status); err != nil {
		utils.InternalError(w, "更新失败")
		return
	}

	utils.SuccessMessage(w, "更新成功")
}
