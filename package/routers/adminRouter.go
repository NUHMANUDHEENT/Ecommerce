package routers

import (
	controller "project1/package/controller/admin"
	"project1/package/middleware"

	"github.com/gin-gonic/gin"
)

var roleAdmin = "admin"

func AdminRouter(r *gin.RouterGroup) {
	//================ admin authentication=======================
	r.GET("/login", controller.AdminLogin)
	r.GET("/logout", controller.AdminLogout)
	r.GET("/", middleware.AuthMiddleware(roleAdmin), controller.AdminPage)

	//================User managment=======================
	r.GET("/user_managment", middleware.AuthMiddleware(roleAdmin), controller.UserList)
	r.PATCH("/user_managment/user_edit/:ID", middleware.AuthMiddleware(roleAdmin), controller.EditUserDetails)
	r.PATCH("/user_managment/user_block/:ID", middleware.AuthMiddleware(roleAdmin), controller.BlockUser)
	r.DELETE("/user_managment/user_delete/:ID", middleware.AuthMiddleware(roleAdmin), controller.DeleteUser)

	//================product managment=======================
	r.GET("/products", middleware.AuthMiddleware(roleAdmin), controller.ProductList)
	r.GET("/products/add_products", middleware.AuthMiddleware(roleAdmin), controller.AddProducts)
	r.POST("/products/add_products", middleware.AuthMiddleware(roleAdmin), controller.UploadImage)
	r.PATCH("products/edit_products/:ID", middleware.AuthMiddleware(roleAdmin), controller.EditProducts)
	r.DELETE("products/delete_products/:ID", middleware.AuthMiddleware(roleAdmin), controller.DeleteProducts)

	//================category managment=======================
	r.GET("/categories", middleware.AuthMiddleware(roleAdmin), controller.CategoryList)
	r.POST("/categories/add_category", middleware.AuthMiddleware(roleAdmin), controller.AddCategory)
	r.PATCH("/categories/edit_category/:ID", middleware.AuthMiddleware(roleAdmin), controller.EditCategories)
	r.DELETE("/categories/delete_category/:ID", middleware.AuthMiddleware(roleAdmin), controller.DeleteCategories)
	r.PATCH("/categories/block_category/:ID", middleware.AuthMiddleware(roleAdmin), controller.BlockCategory)

	//===================== Coupon managment ====================
	r.POST("/coupon", controller.CouponStore)

	// =================== order managment ==============
	r.GET("/orders", controller.AdminOrdersView)
	r.PATCH("/orderstatus/:ID", controller.AdminOrderStatus)
	r.PATCH("/ordercancel/:ID", controller.AdminCancelOrder)

}
