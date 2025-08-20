package main

import (
	"log"

	"emoji-maker-backend/config"
	"emoji-maker-backend/models"
	"emoji-maker-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 初始化SQLite数据库，使用纯Go的SQLite驱动
	engine, err := xorm.NewEngine("sqlite", "./users.db")
	if err != nil {
		panic(err)
	}
	defer engine.Close()

	// 同步数据库表结构
	err = engine.Sync2(new(models.User))
	if err != nil {
		panic(err)
	}

	// 创建fiber应用实例
	app := fiber.New(fiber.Config{
		AppName: "Emoji Maker Backend",
	})

	// 添加中间件
	app.Use(logger.New())
	app.Use(cors.New())

	// 设置路由
	setupRoutes(app, engine)

	// 设置静态文件服务
	app.Static("/tasks", "./tasks")

	// 启动服务器
	log.Fatal(app.ListenTLS(":"+config.AppConfig.Server.Port, "cert.pem", "key.pem"))
}

func setupRoutes(app *fiber.App, engine *xorm.Engine) {
	// 设置视频相关路由
	routes.SetupVideoRoutes(app)

	// 设置用户相关路由
	routes.SetupUserRoutes(app, engine)

	// 默认路由
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Emoji Maker Backend is running!")
	})
}
