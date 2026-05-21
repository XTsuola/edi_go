package routes

import (
	"go_project/models"
	"go_project/service"

	"github.com/gin-gonic/gin"
)

func register(c *gin.Context) {
	var params models.RegisterParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ParamsErr(err.Error(), c)
		return
	}
	err := service.Register(&params)
	if err != nil {
		MyErr(err.Error(), c)
		return
	}
	CreateOk("注册成功", c)
}
