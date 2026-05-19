package controllers

import (
	my "go_project/config"
	"go_project/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	var params models.LoginParams
	if err := c.ShouldBindJSON(&params); err != nil {
		MyErr(err.Error(), c)
		return
	}
	// 根用户名查用户
	var user models.UserAll
	if err := my.DB.Table("users").Where("username = ?", params.Identity).First(&user).Error; err != nil {
		MyErr("用户名不存在", c)
		return
	}
	// 校验密码
	if !CheckPassword(user.Password, params.Password) {
		MyErr("密码错误", c)
		return
	}
	// 生成token
	token, err := GenerateToken(user.ID)
	if err != nil {
		MyErr("生成token失败", c)
		return
	}
	// 将上次登录状态修改为当前时间
	if err2 := my.DB.Table("users").Where("id = ?", user.ID).Updates(map[string]any{
		"last_login": NowTimestamptz(),
		"is_active":  true,
	}).Error; err2 != nil {
		MyErr(err2.Error(), c)
		return
	}
	var userObj models.UserLogin
	userObj.ID = user.ID
	userObj.Username = user.Username
	userObj.Role = user.Role
	userObj.LastLogin = user.LastLogin
	userObj.IsActive = user.IsActive
	data := models.LoginResult{
		Access:  token,
		Refresh: token,
		User:    userObj,
	}
	LoginRes("登录成功", c, data)
}

func logout(c *gin.Context) {
	var params models.LogoutParams
	if err := c.ShouldBindJSON(&params); err != nil {
		MyErr(err.Error(), c)
		return
	}
	// 解析 token
	obj, err := ParseToken(params.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "token无效或已过期",
		})
		c.Abort()
		return
	}
	if err2 := my.DB.Table("users").Where("id = ?", obj.UserID).Updates(map[string]any{
		"is_active": false,
	}).Error; err2 != nil {
		MyErr(err2.Error(), c)
		return
	}
	LogoutRes("退出成功", c)
}
