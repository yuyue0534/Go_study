package controller

import (
	"crm-erp-system/model"
	"crm-erp-system/service"
	"crm-erp-system/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: &service.UserService{},
	}
}

// Register 用户注册
func (ctrl *UserController) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.userService.Register(&req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "注册成功", nil)
}

// Login 用户登录
func (ctrl *UserController) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	token, err := ctrl.userService.Login(&req)
	if err != nil {
		utils.Error(c, 401, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"token": token,
	})
}

// GetUserInfo 获取用户信息
func (ctrl *UserController) GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")

	user, err := ctrl.userService.GetUserInfo(userID.(int64))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, user)
}
