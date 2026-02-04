package router

import (
	"crm-erp-system/controller"
	"crm-erp-system/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "CRM+ERP System API",
			"version": "1.0.0",
		})
	})

	// 初始化控制器
	userCtrl := controller.NewUserController()
	customerCtrl := controller.NewCustomerController()
	productCtrl := controller.NewProductController()
	inventoryCtrl := controller.NewInventoryController()
	orderCtrl := controller.NewOrderController()

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 用户模块（无需鉴权）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userCtrl.Register)
			auth.POST("/login", userCtrl.Login)
		}

		// 需要鉴权的路由
		authorized := v1.Group("/")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 用户信息
			authorized.GET("/user/info", userCtrl.GetUserInfo)

			// CRM客户管理
			customers := authorized.Group("/customers")
			{
				customers.POST("", customerCtrl.Create)
				customers.GET("", customerCtrl.List)
				customers.GET("/:id", customerCtrl.Get)
				customers.PUT("/:id", customerCtrl.Update)
				customers.DELETE("/:id", customerCtrl.Delete)
			}

			// ERP产品管理
			products := authorized.Group("/products")
			{
				products.POST("", productCtrl.Create)
				products.GET("", productCtrl.List)
				products.GET("/:id", productCtrl.Get)
				products.PUT("/:id", productCtrl.Update)
				products.DELETE("/:id", productCtrl.Delete)
			}

			// ERP库存管理
			inventory := authorized.Group("/inventory")
			{
				inventory.POST("", inventoryCtrl.Create)
				inventory.GET("", inventoryCtrl.List)
				inventory.GET("/product/:product_id", inventoryCtrl.GetByProductID)
				inventory.PUT("/product/:product_id", inventoryCtrl.Update)
			}

			// ERP订单管理
			orders := authorized.Group("/orders")
			{
				orders.POST("", orderCtrl.Create)
				orders.GET("", orderCtrl.List)
				orders.GET("/:id", orderCtrl.Get)
				orders.PUT("/:id/status", orderCtrl.UpdateStatus)
				orders.DELETE("/:id", orderCtrl.Delete)
			}
		}
	}

	return r
}
