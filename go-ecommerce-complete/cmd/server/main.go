package main

import (
	"ecommerce/internal/database"
	"ecommerce/internal/handlers"
	"ecommerce/internal/middleware"
	"log"
	"net/http"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}
	defer database.DB.Close()

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", servePage("web/templates/index.html"))
	http.HandleFunc("/login", servePage("web/templates/login.html"))
	http.HandleFunc("/register", servePage("web/templates/register.html"))
	http.HandleFunc("/products", servePage("web/templates/products.html"))
	http.HandleFunc("/product/", servePage("web/templates/product_detail.html"))
	http.HandleFunc("/cart", servePage("web/templates/cart.html"))
	http.HandleFunc("/checkout", servePage("web/templates/checkout.html"))
	http.HandleFunc("/orders", servePage("web/templates/orders.html"))
	http.HandleFunc("/profile", servePage("web/templates/profile.html"))
	http.HandleFunc("/seller", servePage("web/templates/seller.html"))
	http.HandleFunc("/admin", servePage("web/templates/admin.html"))

	http.HandleFunc("/api/register", handlers.Register)
	http.HandleFunc("/api/login", handlers.Login)
	http.HandleFunc("/api/logout", handlers.Logout)
	http.HandleFunc("/api/current-user", handlers.CurrentUser)

	http.HandleFunc("/api/products", handlers.GetProducts)
	http.HandleFunc("/api/product/", handlers.GetProductDetail)
	http.HandleFunc("/api/categories", handlers.GetCategories)

	http.HandleFunc("/api/cart", middleware.AuthMiddleware(handlers.GetCart))
	http.HandleFunc("/api/cart/add", middleware.AuthMiddleware(handlers.AddToCart))
	http.HandleFunc("/api/cart/update", middleware.AuthMiddleware(handlers.UpdateCart))
	http.HandleFunc("/api/cart/remove", middleware.AuthMiddleware(handlers.RemoveFromCart))

	http.HandleFunc("/api/orders", middleware.AuthMiddleware(handlers.GetOrders))
	http.HandleFunc("/api/orders/create", middleware.AuthMiddleware(handlers.CreateOrder))
	http.HandleFunc("/api/orders/pay", middleware.AuthMiddleware(handlers.PayOrder))
	http.HandleFunc("/api/orders/cancel", middleware.AuthMiddleware(handlers.CancelOrder))

	http.HandleFunc("/api/addresses", middleware.AuthMiddleware(handlers.GetAddresses))
	http.HandleFunc("/api/addresses/add", middleware.AuthMiddleware(handlers.AddAddress))
	http.HandleFunc("/api/addresses/delete", middleware.AuthMiddleware(handlers.DeleteAddress))

	http.HandleFunc("/api/seller/products", middleware.RoleMiddleware("seller", "admin")(handlers.SellerGetProducts))
	http.HandleFunc("/api/seller/products/add", middleware.RoleMiddleware("seller")(handlers.SellerAddProduct))
	http.HandleFunc("/api/seller/products/update", middleware.RoleMiddleware("seller")(handlers.SellerUpdateProduct))
	http.HandleFunc("/api/seller/products/delete", middleware.RoleMiddleware("seller")(handlers.SellerDeleteProduct))
	http.HandleFunc("/api/seller/orders", middleware.RoleMiddleware("seller")(handlers.SellerGetOrders))

	http.HandleFunc("/api/admin/products/approve", middleware.RoleMiddleware("admin")(handlers.AdminApproveProduct))
	http.HandleFunc("/api/admin/products/reject", middleware.RoleMiddleware("admin")(handlers.AdminRejectProduct))
	http.HandleFunc("/api/admin/users", middleware.RoleMiddleware("admin")(handlers.AdminGetUsers))
	http.HandleFunc("/api/admin/stats", middleware.RoleMiddleware("admin")(handlers.AdminGetStats))

	log.Println("====================================")
	log.Println("电商平台服务器启动成功！")
	log.Println("====================================")
	log.Println("访问地址: http://localhost:8080")
	log.Println("")
	log.Println("测试账号：")
	log.Println("  管理员: admin / admin123")
	log.Println("  商家: seller1 / seller123")
	log.Println("  用户: customer1 / customer123")
	log.Println("====================================")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

func servePage(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}
