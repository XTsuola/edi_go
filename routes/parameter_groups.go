package routes

import (
	my "go_project/config"
	"go_project/models"
	"go_project/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 获取参数分组列表
func parameterGroupList(c *gin.Context) {
	db := my.DB.Table("parameter_groups").Where("is_deleted = ?", false)
	name := c.Query("search")
	if name != "" {
		db = db.Where("group_name LIKE ?", "%"+name+"%")
	}
	var count int64
	db.Count(&count)
	var list []models.ParameterGroupsList
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	page, _ := strconv.Atoi(c.Query("page"))
	offset := (page - 1) * pageSize
	db = db.Offset(offset).Limit(pageSize).Order("id ASC")
	err := db.Find(&list).Error
	if err != nil {
		MyErr(err.Error(), c)
		return
	}
	SearchByPage[models.ParameterGroupsList](c, list, count)
}

// 新增参数分组
func parameterGroupAdd(c *gin.Context) {
	var params models.ParameterGroupsAddParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ParamsErr(err.Error(), c)
		return
	}
	var count int64
	my.DB.Table("parameter_groups").Where("name = ? AND is_deleted = ?", params.GroupName, false).Count(&count)
	if count > 0 {
		ParamsErr("名称已存在，不能重复添加", c)
		return
	}
	nowTime := utils.NowTimestamptz()
	data := models.ParameterGroupsAddData{
		GroupId:     uuid.New(),
		GroupName:   params.GroupName,
		Description: params.Description,
		CreatedTime: nowTime,
		UpdatedTime: nowTime,
		IsDeleted:   false,
	}
	if err := my.DB.Table("parameter_groups").Create(&data).Error; err != nil {
		MyErr(err.Error(), c)
		return
	}
	CreateOk("新增成功", c)
}

// 修改参数分组
func parameterGroupUpdate(c *gin.Context) {
	id := c.Param("id")
	var params models.ParameterGroupsUpdate
	if err := c.ShouldBindJSON(&params); err != nil {
		MyErr(err.Error(), c)
		return
	}
	updateData := make(map[string]interface{})
	if params.GroupName != nil {
		updateData["group_name"] = *params.GroupName
	}
	if params.Description != nil {
		updateData["description"] = *params.Description
	}
	// 无任何要更新字段，直接返回
	if len(updateData) == 0 {
		HandleOk("修改成功", c)
		return
	}
	updateData["updated_time"] = utils.NowTimestamptz()
	if err2 := my.DB.Table("parameter_groups").Where("id = ?", id).Updates(updateData).Error; err2 != nil {
		MyErr(err2.Error(), c)
		return
	}
	HandleOk("修改成功", c)
}

// 删除参数分组
func parameterGroupDelete(c *gin.Context) {
	id := c.Param("id")
	if err := my.DB.Table("parameter_groups").Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
	}).Error; err != nil {
		MyErr(err.Error(), c)
		return
	}
	DeleteOk("删除成功", c)
}
