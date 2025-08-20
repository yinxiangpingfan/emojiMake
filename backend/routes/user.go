package routes

import (
	"emoji-maker-backend/controllers"
	"emoji-maker-backend/middleware"
	"emoji-maker-backend/repositories"
	"emoji-maker-backend/services"

	"github.com/gofiber/fiber/v2"
	"xorm.io/xorm"
)

// SetupUserRoutes 设置用户相关路由
func SetupUserRoutes(app *fiber.App, engine *xorm.Engine) {
	// 初始化依赖
	userRepo := repositories.NewXormUserRepository(engine)
	jwtService := services.NewJWTService()
	userService := services.NewUserService(userRepo, jwtService)
	userHandler := controllers.NewUserHandler(userService)

	// 创建API路由组
	api := app.Group("/api")
	v1 := api.Group("/v1")
	users := v1.Group("/users")

	// 设置路由
	users.Post("/register", userHandler.Register)
	users.Post("/login", userHandler.Login)
	users.Post("/change-password", middleware.Protected(), userHandler.ChangePassword)
}
