package routes

import (
	"go_project/models"
	"go_project/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	var params models.LoginParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ParamsErr(err.Error(), c)
		return
	}
	data, err := service.Login(&params)
	if err != nil {
		AuthorizedErr(err.Error(), c)
		return
	}
	LoginRes("登录成功", c, *data)
}

func logout(c *gin.Context) {
	var params models.LogoutParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ParamsErr(err.Error(), c)
		return
	}
	err := service.Logout(&params)
	if err != nil {
		if strings.Contains(err.Error(), "token") || strings.Contains(err.Error(), "parse") {
			AuthorizedErr("登录已过期，请重新登录", c)
			return
		}
		MyErr("登出失败，服务器异常", c)
	}
	LogoutRes("登出成功", c)
}
