package routes

import (
	my "go_project/config"
	jwt "go_project/middleware"
	"go_project/models"

	"github.com/gin-gonic/gin"
)

func changePassword(c *gin.Context) {
	var params models.ChangePasswordParams
	if err := c.ShouldBindJSON(&params); err != nil {
		MyErr(err.Error(), c)
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		TokenNull("token已失效", c)
		return
	}
	// 根用户名ID查用户
	var user models.UserAll
	if err := my.DB.Table("users").Where("id = ?", userID).First(&user).Error; err != nil {
		MyErr("用户名不存在", c)
		return
	}
	// 校验密码
	if !jwt.CheckPassword(user.Password, params.OldPassword) {
		MyErr("原密码错误", c)
		return
	}
	// 密码加密
	hashedPwd, err := jwt.HashPassword(params.NewPassword)
	if err != nil {
		MyErr("加密失败", c)
		return
	}
	if err2 := my.DB.Table("users").Where("id = ?", userID).Updates(map[string]interface{}{
		"password": hashedPwd,
	}).Error; err2 != nil {
		MyErr(err2.Error(), c)
		return
	}
	HandleOk("修改成功", c)
}
