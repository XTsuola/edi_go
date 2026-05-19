package models

import "github.com/jackc/pgx/v5/pgtype"

type CategoriesLevel1AddData struct {
	Name        string             `json:"name"`
	CategoryId  int                `json:"category_id"`
	GroupId     string             `json:"group_id"`
	GroupName   string             `json:"group_name"`
	Description string             `json:"description"`
	SortOrder   int                `json:"sort_order"`
	CreatedTime pgtype.Timestamptz `json:"created_time"`
	UpdatedTime pgtype.Timestamptz `json:"updated_time"`
	IsDeleted   bool               `json:"is_deleted"`
}

type CategoriesLevel1List struct {
	ID          int                `json:"id" gorm:"primaryKey"`
	Name        string             `json:"name"`
	CategoryId  int                `json:"category_id"`
	GroupId     string             `json:"group_id"`
	GroupName   string             `json:"group_name"`
	Description string             `json:"description"`
	SortOrder   int                `json:"sort_order"`
	CreatedTime pgtype.Timestamptz `json:"created_time"`
	UpdatedTime pgtype.Timestamptz `json:"updated_time"`
}

type CategoriesLevel1AddParams struct {
	Name        string `json:"name"`
	GroupName   string `json:"group_name"`
	Description string `json:"description"`
}

type GroupSlice struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoriesLevel1Update struct {
	Name        *string `json:"name"`
	GroupId     string  `json:"group_id"`
	GroupName   *string `json:"group_name"`
	Description *string `json:"description"`
}
