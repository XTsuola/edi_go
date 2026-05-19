package controllers

import (
	"github.com/gin-gonic/gin"
	my "go_project/config"
	"go_project/models"
	"strconv"
)

func categoriesLevel2List(c *gin.Context) {
	query := my.DB.Table("categories_level2").Where("is_deleted = ?", false)
	name := c.Query("search")
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	parentId := c.Query("parent_id")
	if parentId != "" {
		query = query.Where("parent_id = ?", parentId)
	}
	var count int64
	query.Count(&count)
	var list []models.CategoriesLevel2List
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	page, _ := strconv.Atoi(c.Query("page"))
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize).Order("id ASC")
	err := query.Find(&list).Error
	if err != nil {
		MyErr(err.Error(), c)
		return
	}
	for i, obj := range list {
		var categories1Obj models.CategoriesLevel1List
		if err2 := my.DB.Table("categories_level1").Where("category_id = ?", obj.Parent).Find(&categories1Obj).Error; err2 != nil {
			MyErr(err2.Error(), c)
			return
		}
		list[i].ParentName = categories1Obj.Name
	}
	SearchByPage[models.CategoriesLevel2List](c, list, count)
}

// 新增分类
func categoriesLevel2Add(c *gin.Context) {
	var params models.CategoriesLevel2AddParams
	if err := c.ShouldBindJSON(&params); err != nil {
		MyErr(err.Error(), c)
		return
	}
	var count int64
	my.DB.Table("categories_level2").Where("name = ? AND is_deleted = ?", params.Name, false).Count(&count)
	if count > 0 {
		ParamsErr("名称已存在，不能重复添加", c)
		return
	}
	// 查询categories_id最大值
	var categoryId int
	my.DB.Table("categories_level2").Select("MAX(category_id)").Scan(&categoryId)
	//// 查询sort_order最大值
	//var count2 int64
	//var sortOrder int
	//my.DB.Table("categories_level2").Where("is_deleted = ?", false).Count(&count2)
	//if count2 == 0 {
	//	sortOrder = 0
	//} else {
	//	my.DB.Table("categories_level2").Select("MAX(sort_order)").Scan(&sortOrder)
	//	sortOrder += 1
	//}
	nowTime := NowTimestamptz()
	data := models.CategoriesLevel2AddData{
		Name:         params.Name,
		CategoryId:   categoryId + 1,
		ParentId:     params.Parent,
		ParamGroupId: params.ParamGroup,
		Description:  params.Description,
		SortOrder:    0,
		CreatedTime:  nowTime,
		UpdatedTime:  nowTime,
		IsDeleted:    false,
	}
	if err2 := my.DB.Table("categories_level2").Create(&data).Error; err2 != nil {
		MyErr(err2.Error(), c)
		return
	}
	CreateOk("新增成功", c)
}

// 修改分类
func categoriesLevel2Update(c *gin.Context) {
	id := c.Param("id")
	var params models.CategoriesLevel2Update
	if err := c.ShouldBindJSON(&params); err != nil {
		MyErr(err.Error(), c)
		return
	}
	updateData := make(map[string]interface{})
	if params.Name != nil {
		updateData["name"] = *params.Name
	}
	if params.Parent != nil {
		updateData["parent_id"] = *params.Parent
	}
	if params.ParamGroup != nil {
		updateData["param_group_id"] = *params.ParamGroup
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
	if err2 := my.DB.Table("categories_level2").Where("id = ?", id).Updates(updateData).Error; err2 != nil {
		MyErr(err2.Error(), c)
		return
	}
	HandleOk("修改成功", c)
}

// 删除子类
func categoriesLevel2Delete(c *gin.Context) {
	id := c.Param("id")
	if err := my.DB.Table("categories_level2").Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
	}).Error; err != nil {
		MyErr(err.Error(), c)
		return
	}
	DeleteOk("删除成功", c)
}
