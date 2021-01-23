package application

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/boof/umg/auth"
	"github.com/boof/umg/controller"
	"github.com/boof/umg/settings"
)

func mapRoutes(e *echo.Echo) {
	mapPublicRoutes(e)
	mapUserRoutes(e)
	mapAdminRoutes(e)
}

func mapPublicRoutes(e *echo.Echo) {
	public := e.Group(settings.BaseURL)

	public.POST("login", controller.Login)
	public.POST("password/reset/email", controller.ResetPasswordByEmail)
	public.POST("password/validate", controller.ValidateResetPassToken)
	public.POST("password/reset", controller.ChangePassword)
}

func mapUserRoutes(e *echo.Echo) {
	user := e.Group(settings.BaseURL)

	user.Use(middleware.JWTWithConfig(auth.GetMiddlewareConfig()))
	user.Use(auth.UserHandler)

	user.GET("user/domains", controller.GetUserDomains)
	user.GET("user/properties", controller.GetProperties)
	user.GET("user/domain/:id/products", controller.GetUserProducts)
	user.GET("user/:id", controller.GetUser)
	user.PUT("user", controller.UpdateUser)
}

func mapAdminRoutes(e *echo.Echo) {
	admin := e.Group(settings.BaseURL)
	admin.Use(middleware.JWTWithConfig(auth.GetMiddlewareConfig()))
	admin.Use(auth.AdminHandler)

	admin.GET("domains", controller.GetDomains)
	admin.GET("domain/:id/products", controller.GetProducts)
	admin.GET("users", controller.GetUsers)
	admin.GET("roles", controller.GetRoles)
	admin.GET("role/:id/policies", controller.GetPolicies)

	admin.GET("user/:id/email/welcome_reset", controller.SendWelcomeAndReset)
	admin.GET("user/:id/email/history", controller.GetUserEmailHistory)

	admin.GET("search/users", controller.SearchUsers)
	admin.GET("search/roles", controller.SearchRoles)

	admin.POST("role", controller.AddRole)
	admin.POST("role/with-policy", controller.AddRoleWithPolicy)
	admin.POST("role/assign", controller.AssignRole)
	admin.POST("role/disallow", controller.DisallowRole)
	admin.POST("user", controller.AddUser)
	admin.POST("user/with-role", controller.AddUserWithRole)
	admin.POST("domain", controller.AddDomain)
	admin.POST("product", controller.AddProduct)
	admin.POST("policy/domain", controller.AddDomPolicy)
	admin.POST("policy/product", controller.AddProdPolicy)
	admin.POST("policy/product/all", controller.AddAllProdPolicy)

	admin.POST("user/access/expire", controller.AddAccessExpire)
	admin.PUT("user/access/expire", controller.EditAccessExpire)
	admin.DELETE("user/:id/access/expire", controller.DelAccessExpire)

	admin.POST("password/change", controller.ChangeUserPassword)

	admin.DELETE("policy/:id", controller.DelPolicy)
	admin.DELETE("role/:id", controller.DelRole)
	admin.DELETE("user/:id", controller.DelUser)

	admin.PUT("role", controller.EditRole)
}
