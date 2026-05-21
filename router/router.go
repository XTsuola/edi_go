package router

import (
	"go_project/config"
	"go_project/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	routes.R.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有域，生产环境建议改为具体域名
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	routes.R.Static("/static", "./dist")
	routes.R.Static("/upload", "./upload")
	routes.R.NoRoute(func(c *gin.Context) { // 兼容 Vue/React 单页应用路由
		c.File("./dist/index.html")
	})
	config.InitDB()
	routes.InitController()
	err := routes.R.Run(":8002")
	if err != nil {
		return
	}

	//StartHTTPS()
}
