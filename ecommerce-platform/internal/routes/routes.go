package routes

import (
	"ecommerce-platform/internal/handlers"
	"ecommerce-platform/internal/middleware"
	"net/http"
)

// SetupRoutes 配置路由
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// 静态文件 - 必须在其他路由之前
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// 上传文件
	uploadFs := http.FileServer(http.Dir("./uploads"))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", uploadFs))

	// API路由
	// 用户认证
	mux.HandleFunc("/api/auth/register", handlers.Register)
	mux.HandleFunc("/api/auth/login", handlers.Login)
	mux.HandleFunc("/api/auth/logout", handlers.Logout)
	mux.HandleFunc("/api/auth/user", middleware.Auth(handlers.GetCurrentUser))

	// 用户管理
	mux.HandleFunc("/api/user/profile", middleware.Auth(handlers.UpdateProfile))
	mux.HandleFunc("/api/user/password", middleware.Auth(handlers.ChangePassword))
	mux.HandleFunc("/api/user/seller", middleware.RequireSeller(handlers.UpdateSellerInfo))

	// 商品分类
	mux.HandleFunc("/api/categories", handlers.GetCategories)
	mux.HandleFunc("/api/admin/category/create", middleware.RequireAdmin(handlers.CreateCategory))
	mux.HandleFunc("/api/admin/category/update", middleware.RequireAdmin(handlers.UpdateCategory))
	mux.HandleFunc("/api/admin/category/delete", middleware.RequireAdmin(handlers.DeleteCategory))

	// 商品
	mux.HandleFunc("/api/products", handlers.GetProducts)
	mux.HandleFunc("/api/product", handlers.GetProduct)
	mux.HandleFunc("/api/products/hot", handlers.GetHotProducts)
	mux.HandleFunc("/api/products/new", handlers.GetNewProducts)
	mux.HandleFunc("/api/upload/image", middleware.Auth(handlers.UploadImage))

	// 商家商品管理
	mux.HandleFunc("/api/seller/products", middleware.RequireSeller(handlers.GetSellerProducts))
	mux.HandleFunc("/api/seller/product/create", middleware.RequireSeller(handlers.CreateProduct))
	mux.HandleFunc("/api/seller/product/update", middleware.RequireSeller(handlers.UpdateProduct))
	mux.HandleFunc("/api/seller/product/delete", middleware.RequireSeller(handlers.DeleteProduct))
	mux.HandleFunc("/api/seller/product/stock", middleware.RequireSeller(handlers.UpdateProductStock))

	// 购物车
	mux.HandleFunc("/api/cart", middleware.Auth(handlers.GetCart))
	mux.HandleFunc("/api/cart/add", middleware.Auth(handlers.AddToCart))
	mux.HandleFunc("/api/cart/update", middleware.Auth(handlers.UpdateCartItem))
	mux.HandleFunc("/api/cart/select", middleware.Auth(handlers.UpdateCartItemSelected))
	mux.HandleFunc("/api/cart/select-all", middleware.Auth(handlers.SelectAllCart))
	mux.HandleFunc("/api/cart/delete", middleware.Auth(handlers.DeleteCartItem))
	mux.HandleFunc("/api/cart/clear", middleware.Auth(handlers.ClearCart))
	mux.HandleFunc("/api/cart/count", middleware.Auth(handlers.GetCartCount))

	// 收货地址
	mux.HandleFunc("/api/addresses", middleware.Auth(handlers.GetAddresses))
	mux.HandleFunc("/api/address/create", middleware.Auth(handlers.CreateAddress))
	mux.HandleFunc("/api/address/update", middleware.Auth(handlers.UpdateAddress))
	mux.HandleFunc("/api/address/delete", middleware.Auth(handlers.DeleteAddress))
	mux.HandleFunc("/api/address/default", middleware.Auth(handlers.SetDefaultAddress))

	// 订单
	mux.HandleFunc("/api/order/create", middleware.Auth(handlers.CreateOrder))
	mux.HandleFunc("/api/orders", middleware.Auth(handlers.GetOrders))
	mux.HandleFunc("/api/order", middleware.Auth(handlers.GetOrder))
	mux.HandleFunc("/api/order/pay", middleware.Auth(handlers.PayOrder))
	mux.HandleFunc("/api/order/cancel", middleware.Auth(handlers.CancelOrder))
	mux.HandleFunc("/api/order/receive", middleware.Auth(handlers.ReceiveOrder))

	// 商家订单管理
	mux.HandleFunc("/api/seller/orders", middleware.RequireSeller(handlers.GetSellerOrders))
	mux.HandleFunc("/api/seller/order/ship", middleware.RequireSeller(handlers.ShipOrder))
	mux.HandleFunc("/api/seller/order/stats", middleware.RequireSeller(handlers.GetSellerOrderStats))

	// 评价
	mux.HandleFunc("/api/reviews", handlers.GetProductReviews)
	mux.HandleFunc("/api/review/create", middleware.Auth(handlers.CreateReview))
	mux.HandleFunc("/api/user/reviews", middleware.Auth(handlers.GetUserReviews))
	mux.HandleFunc("/api/seller/reviews", middleware.RequireSeller(handlers.GetSellerReviews))
	mux.HandleFunc("/api/seller/review/reply", middleware.RequireSeller(handlers.ReplyReview))

	// 管理员
	mux.HandleFunc("/api/admin/users", middleware.RequireAdmin(handlers.AdminGetUsers))
	mux.HandleFunc("/api/admin/user/status", middleware.RequireAdmin(handlers.AdminUpdateUserStatus))
	mux.HandleFunc("/api/admin/sellers", middleware.RequireAdmin(handlers.AdminGetSellers))
	mux.HandleFunc("/api/admin/seller/status", middleware.RequireAdmin(handlers.AdminUpdateSellerStatus))
	mux.HandleFunc("/api/admin/products", middleware.RequireAdmin(handlers.AdminGetProducts))
	mux.HandleFunc("/api/admin/product/status", middleware.RequireAdmin(handlers.AdminUpdateProductStatus))
	mux.HandleFunc("/api/admin/orders", middleware.RequireAdmin(handlers.AdminGetOrders))
	mux.HandleFunc("/api/admin/order/stats", middleware.RequireAdmin(handlers.AdminGetOrderStats))

	// 页面路由
	mux.HandleFunc("/login", serveIndex)
	mux.HandleFunc("/register", serveIndex)
	mux.HandleFunc("/product/", serveIndex)
	mux.HandleFunc("/cart", serveIndex)
	mux.HandleFunc("/checkout", serveIndex)
	mux.HandleFunc("/orders", serveIndex)
	mux.HandleFunc("/order/", serveIndex)
	mux.HandleFunc("/user/", serveIndex)
	mux.HandleFunc("/seller/", serveIndex)
	mux.HandleFunc("/admin/", serveIndex)
	
	// 根路由必须放在最后，且需要精确匹配
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 只处理精确的根路径
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		serveIndex(w, r)
	})

	// 应用中间件
	handler := middleware.Logger(middleware.CORS(mux))

	return handler
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}
