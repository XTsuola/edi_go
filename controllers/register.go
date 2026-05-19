package controllers

import (
	"github.com/gin-gonic/gin"
	my "go_project/config"
	"go_project/models"
)

func register(c *gin.Context) {
	var params models.RegisterParams
	if err := c.ShouldBindJSON(&params); err != nil {
		MyErr(err.Error(), c)
		return
	}
	// 判断邮箱是否已存在
	var count int64
	my.DB.Table("users").Where("email = ?", params.Email).Count(&count)
	if count > 0 {
		MyErr("邮箱已被注册", c)
		return
	}
	// 密码加密
	hashedPwd, err := HashPassword(params.Password)
	if err != nil {
		c.JSON(500, gin.H{"msg": "加密失败"})
		return
	}
	user := models.RegisterUserData{
		Username:    params.Username,
		Email:       params.Email,
		Password:    hashedPwd, // 密文
		IsActive:    true,
		IsSuperuser: false,
		IsStaff:     false,
		Role:        "user",
		DataJoined:  NowTimestamptz(),
	}
	if err2 := my.DB.Table("users").Create(&user).Error; err2 != nil {
		MyErr(err2.Error(), c)
		return
	}
	CreateOk("注册成功", c)
}
