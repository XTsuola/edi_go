package routes

import (
	jwt "go_project/middleware"

	"github.com/gin-gonic/gin"
)

var R = gin.Default()

func InitController() {
	gin.SetMode(gin.ReleaseMode) // 或 gin.DebugMode
	// 统一抽取公共前缀 /api/v1/
	v1Group := R.Group("/api/v1")
	{

		// 用户相关接口
		v1Group.POST("/auth/users/register/", register)
		v1Group.POST("/auth/users/login/", login)

		// 需要token认证的子分组
		authGroup := v1Group.Group("/")
		authGroup.Use(jwt.AuthMiddleware()) // 启用中间件
		{
			// 只有携带有效token才能访问
			authGroup.POST("/users/change-password/", changePassword)
			//authGroup.POST("/upload/file/", uploadFile)
			authGroup.POST("/file_service/upload/chunked/", uploadChunkHandler)
			authGroup.POST("/file_service/upload/process/", processFile)
			authGroup.POST("/auth/users/logout/", logout)
			authGroup.GET("/categories/level1/", categoriesLevel1List)
			authGroup.POST("/categories/level1/", categoriesLevel1Add)
			authGroup.PATCH("/categories/level1/:id/", categoriesLevel1Update)
			authGroup.DELETE("/categories/level1/:id/", categoriesLevel1Delete)
			authGroup.GET("/categories/level2/", categoriesLevel2List)
			authGroup.POST("/categories/level2/", categoriesLevel2Add)
			authGroup.PATCH("/categories/level2/:id/", categoriesLevel2Update)
			authGroup.DELETE("/categories/level2/:id/", categoriesLevel2Delete)
			authGroup.GET("/parameter-groups/", parameterGroupList)
			authGroup.POST("/parameter-groups/", parameterGroupAdd)
			authGroup.PATCH("/parameter-groups/:id/", parameterGroupUpdate)
			authGroup.DELETE("/parameter-groups/:id/", parameterGroupDelete)
			authGroup.GET("/parameter-templates/", parameterTemplatesList)
			authGroup.POST("/parameter-templates/", parameterTemplatesAdd)
			authGroup.PATCH("/parameter-templates/:id/", parameterTemplatesUpdate)
			authGroup.DELETE("/parameter-templates/:id/", parameterTemplatesDelete)
		}
	}
}
