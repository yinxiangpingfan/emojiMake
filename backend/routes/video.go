package routes

import (
	"emoji-maker-backend/controllers"
	"emoji-maker-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupVideoRoutes(app *fiber.App) {
	// 视频相关路由
	video := app.Group("/api/v1/video", middleware.Protected())

	// 创建视频生成任务
	video.Post("/create", controllers.CreateVideoTask)

	// 创建视频生成任务 (带提示词处理)
	video.Post("/create_with_prompt", controllers.CreateVideoTaskWithPromptProcessing)

	// 查询任务结果
	video.Get("/query/:job_id", controllers.GetVideoTaskResult)
}
