package routes

import (
	"go_project/models"
	"go_project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
200 = http.StatusOK        // 成功
201 = http.StatusCreated   // 创建成功
204 = http.StatusNoContent // 无内容返回
205 = http.StatusResetContent // 清空表单/重置页面
400 = http.StatusBadRequest // 参数错误
401 = http.StatusUnauthorized // 未登录
403 = http.StatusForbidden  // 无权限
500 = http.StatusInternalServerError // 服务器错误
*/

// MyErr 接口500报错
func MyErr(err string, c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": err,
	})
}

// ParamsErr 接口参数400报错
func ParamsErr(err string, c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": err,
	})
}

// AuthorizedErr 接口参数400报错
func AuthorizedErr(err string, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": err,
	})
}

// LoginRes 登录
func LoginRes[T models.LoginResult](msg string, c *gin.Context, data models.LoginResult) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    data,
		"message": msg,
	})
}

// LogoutRes 登出
func LogoutRes(msg string, c *gin.Context) {
	c.JSON(http.StatusResetContent, gin.H{
		"message": msg,
	})
}

// SearchList 查询列表成功
func SearchList[T any](msg string, c *gin.Context, data []T) {
	if len(data) == 0 {
		data = []T{}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  utils.If(msg == "", "success", msg),
		"data": data,
	})
}

// SearchOne 查询单个
func SearchOne[T any](msg string, c *gin.Context, data T) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  utils.If(msg == "", "success", msg),
		"data": data,
	})
}

// SearchByPage 分页查询成功
func SearchByPage[T any](c *gin.Context, data []T, count int64) {
	if data == nil {
		data = []T{}
	}
	c.JSON(http.StatusOK, gin.H{
		"results": data,
		"count":   count,
	})
}

// CreateOk 创建成功 201
func CreateOk(msg string, c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": msg,
	})
}

// HandleOk 查询成功 200
func HandleOk(msg string, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": msg,
	})
}

// DeleteOk 删除成功 204
func DeleteOk(msg string, c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{
		"message": msg,
	})
}

// TokenNull 未授权 401
func TokenNull(msg string, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": msg,
	})
}
