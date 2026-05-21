package service

import (
	"errors"
	my "go_project/config"
	jwt "go_project/middleware"
	models "go_project/models"
	"go_project/utils"

	"github.com/google/uuid"
)

// Login 登录模块
func Login(params *models.LoginParams) (*models.LoginResult, error) {
	// 根用户名查用户
	var user models.UserAll
	if err := my.DB.Table("users").Where("username = ?", params.Identity).First(&user).Error; err != nil {
		return nil, errors.New("用户名不存在")
	}
	// 校验密码
	if !jwt.CheckPassword(user.Password, params.Password) {
		return nil, errors.New("密码错误")
	}
	// 生成accessToken
	accessToken, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}
	// 生成refreshToken
	refreshToken, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}
	// 修改最后登录时间和登录状态
	if err2 := my.DB.Table("users").Where("id = ?", user.ID).Updates(map[string]any{
		"last_login": utils.NowTimestamptz(),
		"is_active":  true,
	}).Error; err2 != nil {
		return nil, err2
	}
	// 返回token和用户信息
	var userObj models.UserLogin
	userObj.ID = user.ID
	userObj.Username = user.Username
	userObj.Role = user.Role
	userObj.LastLogin = user.LastLogin
	userObj.IsActive = user.IsActive
	data := models.LoginResult{
		Access:  accessToken,
		Refresh: refreshToken,
		User:    userObj,
	}
	return &data, nil
}

// Logout 登出模块
func Logout(params *models.LogoutParams) error {
	// 解析 token
	obj, err := jwt.ParseToken(params.RefreshToken)
	if err != nil {
		return err
	}
	// 修改当前登录状态
	if err2 := my.DB.Table("users").Where("id = ?", obj.UserID).Updates(map[string]any{
		"is_active": false,
	}).Error; err2 != nil {
		return err2
	}
	return nil
}

// Register 注册模块
func Register(params *models.RegisterParams) error {
	// 判断邮箱是否已存在
	var count int64
	my.DB.Table("users").Where("email = ?", params.Email).Count(&count)
	if count > 0 {
		return errors.New("邮箱已被注册")
	}
	// 密码加密
	hashedPwd, err := jwt.HashPassword(params.Password)
	if err != nil {
		return errors.New("加密失败")
	}
	user := models.RegisterUserData{
		ID:          uuid.New(),
		Username:    params.Username,
		Email:       params.Email,
		Password:    hashedPwd, // 密文
		IsActive:    true,
		IsSuperuser: false,
		IsStaff:     false,
		Role:        "user",
		DataJoined:  utils.NowTimestamptz(),
	}
	if err2 := my.DB.Table("users").Create(&user).Error; err2 != nil {
		return err2
	}
	return nil
}
