package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ParameterGroupsList struct {
	ID          int                `json:"id" gorm:"primaryKey"`
	GroupId     string             `json:"group_id"`
	GroupName   string             `json:"group_name"`
	Description string             `json:"description"`
	CreatedTime pgtype.Timestamptz `json:"created_time"`
	UpdatedTime pgtype.Timestamptz `json:"updated_time"`
}

type ParameterGroupsAddData struct {
	GroupId     uuid.UUID          `json:"group_id"`
	GroupName   string             `json:"group_name"`
	Description string             `json:"description"`
	CreatedTime pgtype.Timestamptz `json:"created_time"`
	UpdatedTime pgtype.Timestamptz `json:"updated_time"`
	IsDeleted   bool               `json:"is_deleted"`
}

type ParameterGroupsAddParams struct {
	GroupName   string `json:"group_name"`
	Description string `json:"description"`
}

type ParameterGroupsUpdate struct {
	GroupName   *string `json:"group_name"`
	Description *string `json:"description"`
}
