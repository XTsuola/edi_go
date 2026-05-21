package routes

import (
	my "go_project/config"
	"go_project/models"
	"go_project/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 总类分组
var groupSlice = []models.GroupSlice{
	{ID: "TestModule", Name: "拟合模型"},
	{ID: "InversionModel", Name: "反演模型"},
	{ID: "VswrModel", Name: "天线模型"},
	{ID: "ChannelModel", Name: "链路模型"},
}

// GetGroupIDByName 根据分组ID获取分组名
func GetGroupIDByName(name string) string {
	for _, g := range groupSlice {
		if g.Name == name {
			return g.ID
		}
	}
	return ""
}

// 获取总类列表
func categoriesLevel1List(c *gin.Context) {
	db := my.DB.Table("categories_level1").Where("is_deleted = ?", false)
	name := c.Query("search")
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	groupId := c.Query("group_id")
	if groupId != "" {
		db = db.Where("group_id = ?", groupId)
	}
	var count int64
	db.Count(&count)
	var list []models.CategoriesLevel1List
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	page, _ := strconv.Atoi(c.Query("page"))
	offset := (page - 1) * pageSize
	db = db.Offset(offset).Limit(pageSize).Order("id ASC")
	err := db.Find(&list).Error
	if err != nil {
		MyErr(err.Error(), c)
		return
	}
	SearchByPage[models.CategoriesLevel1List](c, list, count)
}

// 新增总类
func categoriesLevel1Add(c *gin.Context) {
	var params models.CategoriesLevel1AddParams
	if err := c.ShouldBindJSON(&params); err != nil {
		ParamsErr(err.Error(), c)
		return
	}
	var count int64
	my.DB.Table("categories_level1").Where("name = ? AND is_deleted = ?", params.GroupName, false).Count(&count)
	if count > 0 {
		ParamsErr("名称已存在，不能重复添加", c)
		return
	}
	// 查询categories_id最大值
	var categoryId int
	my.DB.Table("categories_level1").Select("MAX(category_id)").Scan(&categoryId)
	// 查询sort_order最大值
	var count2 int64
	var sortOrder int
	my.DB.Table("categories_level1").Where("is_deleted = ?", false).Count(&count2)
	if count2 == 0 {
		sortOrder = 0
	} else {
		my.DB.Table("categories_level1").Select("MAX(sort_order)").Scan(&sortOrder)
		sortOrder += 1
	}
	nowTime := utils.NowTimestamptz()
	data := models.CategoriesLevel1AddData{
		Name:        params.Name,
		CategoryId:  categoryId + 1,
		GroupId:     GetGroupIDByName(params.GroupName),
		GroupName:   params.GroupName,
		Description: params.Description,
		SortOrder:   sortOrder,
		CreatedTime: nowTime,
		UpdatedTime: nowTime,
		IsDeleted:   false,
	}
	if err := my.DB.Table("categories_level1").Create(&data).Error; err != nil {
		MyErr(err.Error(), c)
		return
	}
	CreateOk("新增成功", c)
}

// 修改总类
func categoriesLevel1Update(c *gin.Context) {
	id := c.Param("id")
	var params models.CategoriesLevel1Update
	if err := c.ShouldBindJSON(&params); err != nil {
		ParamsErr(err.Error(), c)
		return
	}
	updateData := make(map[string]interface{})
	if params.Name != nil {
		updateData["name"] = *params.Name
	}
	if params.GroupName != nil {
		updateData["group_name"] = *params.GroupName
		updateData["group_id"] = GetGroupIDByName(*params.GroupName)
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
	if err := my.DB.Table("categories_level1").Where("id = ?", id).Updates(updateData).Error; err != nil {
		MyErr(err.Error(), c)
		return
	}
	HandleOk("修改成功", c)
}

// 删除总类
func categoriesLevel1Delete(c *gin.Context) {
	id := c.Param("id")
	if err := my.DB.Table("categories_level1").Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
	}).Error; err != nil {
		MyErr(err.Error(), c)
		return
	}
	DeleteOk("删除成功", c)
}
