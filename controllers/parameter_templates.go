package controllers

import (
	my "go_project/config"
	"go_project/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取参数模板列表
func parameterTemplatesList(c *gin.Context) {
	db := my.DB.Table("parameter_templates").Where("is_deleted = ?", false)
	name := c.Query("search")
	paramGroupId := c.Query("param_group_id")
	if paramGroupId != "" {
		db = db.Where("param_group_id", paramGroupId)
	}
	if name != "" {
		db = db.Where("param_name LIKE ?", "%"+name+"%")
	}
	var count int64
	db.Count(&count)
	var list []models.ParameterTemplatesList
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	page, _ := strconv.Atoi(c.Query("page"))
	offset := (page - 1) * pageSize
	db = db.Offset(offset).Limit(pageSize).Order("id ASC")
	err := db.Find(&list).Error
	if err != nil {
		MyErr(err.Error(), c)
		return
	}
	SearchByPage[models.ParameterTemplatesList](c, list, count)
}

// 新增参数模板
func parameterTemplatesAdd(c *gin.Context) {
	var params models.ParameterTemplatesAddParams
	if err := c.ShouldBindJSON(&params); err != nil {
		MyErr(err.Error(), c)
		return
	}
	var count int64
	my.DB.Table("parameter_templates").Where("param_key = ? AND is_deleted = ?", params.ParamKey, false).Count(&count)
	if count > 0 {
		ParamsErr("参数键名称已存在，不能重复添加", c)
		return
	}
	nowTime := NowTimestamptz()
	data := models.ParameterTemplatesAddData{
		ParamName:       params.ParamName,
		ParamGroupId:    params.ParamGroupId,
		ParamKey:        params.ParamKey,
		DataType:        params.DataType,
		Unit:            params.Unit,
		OptionalUnits:   params.OptionalUnits,
		IsRequired:      params.IsRequired,
		ValidationRules: params.ValidationRules,
		SortOrder:       params.SortOrder,
		Description:     params.Description,
		CreatedTime:     nowTime,
		UpdatedTime:     nowTime,
		IsDeleted:       false,
	}
	if err2 := my.DB.Table("parameter_templates").Create(&data).Error; err2 != nil {
		MyErr(err2.Error(), c)
		return
	}
	CreateOk("新增成功", c)
}

// 修改参数模板
func parameterTemplatesUpdate(c *gin.Context) {
	id := c.Param("id")
	var params models.ParameterTemplatesUpdate
	if err := c.ShouldBindJSON(&params); err != nil {
		MyErr(err.Error(), c)
		return
	}
	updateData := make(map[string]interface{})
	if params.ParamName != nil {
		updateData["param_name"] = *params.ParamName
	}
	if params.ParamGroupId != nil {
		updateData["param_group_id"] = *params.ParamGroupId
	}
	if params.ParamKey != nil {
		updateData["param_key"] = *params.ParamKey
	}
	if params.DataType != nil {
		updateData["data_type"] = *params.DataType
	}
	if params.Unit != nil {
		updateData["Unit"] = *params.Unit
	}
	if params.OptionalUnits != nil {
		updateData["OptionalUnits"] = *params.OptionalUnits
	}
	if params.IsRequired != nil {
		updateData["is_required"] = *params.IsRequired
	}
	if params.ValidationRules != nil {
		updateData["validation_rules"] = *params.ValidationRules
	}
	if params.SortOrder != nil {
		updateData["sort_order"] = *params.SortOrder
	}
	if params.Description != nil {
		updateData["description"] = *params.Description
	}
	// 无任何要更新字段，直接返回
	if len(updateData) == 0 {
		HandleOk("修改成功", c)
		return
	}
	updateData["updated_time"] = NowTimestamptz()
	if err2 := my.DB.Table("parameter_templates").Where("id = ?", id).Updates(updateData).Error; err2 != nil {
		MyErr(err2.Error(), c)
		return
	}
	HandleOk("修改成功", c)
}

// 删除参数模板
func parameterTemplatesDelete(c *gin.Context) {
	id := c.Param("id")
	if err := my.DB.Table("parameter_templates").Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
	}).Error; err != nil {
		MyErr(err.Error(), c)
		return
	}
	DeleteOk("删除成功", c)
}
