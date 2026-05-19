package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CategoriesLevel2AddData struct {
	Name         string             `json:"name"`
	CategoryId   int                `json:"category_id"`
	ParentId     int                `json:"parent_id"`
	ParamGroupId uuid.UUID          `json:"param_group_id"`
	Description  string             `json:"description"`
	SortOrder    int                `json:"sort_order"`
	CreatedTime  pgtype.Timestamptz `json:"created_time"`
	UpdatedTime  pgtype.Timestamptz `json:"updated_time"`
	IsDeleted    bool               `json:"is_deleted"`
}

type CategoriesLevel2List struct {
	ID          int                `json:"id" gorm:"primaryKey"`
	Name        string             `json:"name"`
	CategoryId  int                `json:"category_id"`
	Parent      int                `json:"parent" gorm:"column:parent_id"`
	ParentName  string             `json:"parent_name"`
	ParamGroup  uuid.UUID          `json:"param_group" gorm:"column:param_group_id"`
	Description string             `json:"description"`
	SortOrder   int                `json:"sort_order"`
	CreatedTime pgtype.Timestamptz `json:"created_time"`
	UpdatedTime pgtype.Timestamptz `json:"updated_time"`
}

type CategoriesLevel2AddParams struct {
	Name        string    `json:"name"`
	Parent      int       `json:"parent"`
	ParamGroup  uuid.UUID `json:"param_group"`
	Description string    `json:"description"`
}

type CategoriesLevel2Update struct {
	Name        *string    `json:"name"`
	Parent      *int       `json:"parent"`
	ParamGroup  *uuid.UUID `json:"param_group"`
	Description *string    `json:"description"`
}
