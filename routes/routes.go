package routes

import (
	"invar/controllers"
	"invar/middlewares"
	"invar/permission"

	"github.com/gin-gonic/gin"
)

func Setup(rg *gin.RouterGroup) {
	static := rg.Group("static", middlewares.LoggerToFile(), middlewares.Auth(), middlewares.CheckFolderPermission(permission.QueryUser))
	static.Static("/", "./static")

	public := rg.Group("public", middlewares.Auth(), middlewares.LoggerToFile())
	public.Static("/", "./public")

	api := rg.Group("/api", middlewares.LoggerToFile(), middlewares.RequestSizeLimit())
	v1 := api.Group("/v1")
	v1.POST("register", controllers.RegisterUser)
	v1.POST("get_email_code", controllers.GetEmailCode)
	v1.POST("auth", controllers.Authenticate)
	v1.POST("admin_auth", controllers.AuthenticateAdmin)
	v1.POST("refresh_token", controllers.RefreshToken)
	v1.POST("revoke_token", controllers.RevokeToken)
	v1.POST("forget_password", controllers.ForgetPassword)
	v1.GET("forget_password_access", controllers.ForgetPasswordAccess)
	v1.POST("reset_password_with_token", controllers.ResetPasswordWithToken)

	v1WithAuth := v1.Use(middlewares.Auth())

	v1WithAuth.POST("change_password", controllers.ChangeUserPassword)

	v1WithAuth.GET("kyc", controllers.GetKYC)
	v1WithAuth.POST("kyc", controllers.AddKYC)
	v1WithAuth.PATCH("kyc", controllers.UpdateKYC)

	v1WithAuth.GET("whitelist", controllers.GetWhiteLists)
	v1WithAuth.POST("whitelist", controllers.AddWhiteList)
	v1WithAuth.PATCH("whitelist/:id", controllers.UpdateWhiteList)
	v1WithAuth.DELETE("whitelist/:id", controllers.DeleteWhiteList)

	v1WithAuth.GET("product", controllers.GetProducts)
	v1WithAuth.GET("product/:id", controllers.GetProduct)

	v1WithAuth.GET("order", controllers.GetOrders)
	v1WithAuth.POST("order", controllers.AddOrder)
	v1WithAuth.PATCH("cancel_order/:id", controllers.CancelOrder)
	v1WithAuth.PATCH("payment_order/:id", controllers.PaymentOrder)

	v1WithAuth.GET("bank/:id", controllers.GetBank)

	v1WithAuth.GET("stack", controllers.GetStacks)
	v1WithAuth.GET("stack/:id", controllers.GetStack)
	v1WithAuth.GET("stack_record", controllers.GetStacksRecord)

	admin := v1.Group("/admin", middlewares.Auth())

	admin.POST("change_password", controllers.ChangeAdminPassword)
	admin.POST("change_password_by_admin", middlewares.CheckAdminPermission(permission.ModifyUser), controllers.ChangeUserPasswordByAdmin)

	admin.GET("whitelist/:id", middlewares.CheckAdminPermission(permission.QueryWhiteList), controllers.GetWhiteListsByAdmin)
	admin.POST("whitelist", middlewares.CheckAdminPermission(permission.ModifyWhiteList), controllers.AddWhiteList)
	admin.PATCH("whitelist/:id", middlewares.CheckAdminPermission(permission.ModifyWhiteList), controllers.UpdateWhiteList)
	admin.DELETE("whitelist/:id", middlewares.CheckAdminPermission(permission.ModifyWhiteList), controllers.DeleteWhiteList)

	admin.GET("kyc/:id", middlewares.CheckAdminPermission(permission.QueryWhiteList), controllers.GetWhiteListsByAdmin)
	admin.PATCH("kyc/:id", middlewares.CheckAdminPermission(permission.ModifyWhiteList), controllers.UpdateKYCByAdmin)

	admin.GET("bank", middlewares.CheckAdminPermission(permission.QueryBank), controllers.GetBanks)
	admin.GET("bank/:id", middlewares.CheckAdminPermission(permission.QueryBank), controllers.GetBank)

	admin.GET("product", middlewares.CheckAdminPermission(permission.QueryProduct), controllers.GetProductsByAdmin)
	admin.GET("product/:id", middlewares.CheckAdminPermission(permission.QueryProduct), controllers.GetProduct)
	admin.POST("product", middlewares.CheckAdminPermission(permission.ModifyProduct), controllers.AddProduct)
	admin.PATCH("product/:id", middlewares.CheckAdminPermission(permission.ModifyProduct), controllers.UpdateProduct)

	admin.GET("order", middlewares.CheckAdminPermission(permission.QueryOrder), controllers.GetOrdersByAdmin)
	admin.PATCH("cancel_order/:id", middlewares.CheckAdminPermission(permission.ModifyOrder), controllers.CancelOrder)
	admin.PATCH("completed_order/:id", middlewares.CheckAdminPermission(permission.ModifyOrder), controllers.CompletedOrder)

	admin.GET("stack", middlewares.CheckAdminPermission(permission.QueryStack), controllers.GetStacks)
	admin.GET("stack/:id", middlewares.CheckAdminPermission(permission.QueryStack), controllers.GetStack)
	admin.POST("stack", middlewares.CheckAdminPermission(permission.QueryStack), controllers.AddStack)
	admin.PATCH("stack/:id", middlewares.CheckAdminPermission(permission.QueryStack), controllers.UpdateStack)
	admin.GET("stack_record", middlewares.CheckAdminPermission(permission.QueryStack), controllers.GetStacksRecordByAdmin)
	admin.POST("stack_profit_record", middlewares.CheckAdminPermission(permission.QueryStack), controllers.AddStackProfitRecord)
	admin.DELETE("stack_profit_record", middlewares.CheckAdminPermission(permission.QueryStack), controllers.DeleteStackProfitRecord)
}
