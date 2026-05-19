package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ParameterTemplatesList struct {
	ID              int                `json:"id" gorm:"primaryKey"`
	ParamName       string             `json:"param_name"`
	ParamKey        string             `json:"param_key"`
	DataType        string             `json:"data_type"`
	ParamGroupId    uuid.UUID          `json:"param_group_id"`
	Unit            string             `json:"Unit" gorm:"column:Unit"`
	OptionalUnits   json.RawMessage    `json:"OptionalUnits" gorm:"column:OptionalUnits"`
	IsRequired      bool               `json:"is_required"`
	ValidationRules json.RawMessage    `json:"validation_rules"`
	SortOrder       int                `json:"sort_order"`
	Description     string             `json:"description"`
	CreatedTime     pgtype.Timestamptz `json:"created_time"`
	UpdatedTime     pgtype.Timestamptz `json:"updated_time"`
}

type ParameterTemplatesAddData struct {
	ParamName       string             `json:"param_name"`
	ParamKey        string             `json:"param_key"`
	DataType        string             `json:"data_type"`
	Unit            string             `json:"Unit" gorm:"column:Unit"`
	OptionalUnits   json.RawMessage    `json:"OptionalUnits" gorm:"column:OptionalUnits"`
	IsRequired      bool               `json:"is_required"`
	ValidationRules json.RawMessage    `json:"validation_rules"`
	SortOrder       int                `json:"sort_order"`
	Description     string             `json:"description"`
	CreatedTime     pgtype.Timestamptz `json:"created_time"`
	UpdatedTime     pgtype.Timestamptz `json:"updated_time"`
	IsDeleted       bool               `json:"is_deleted"`
	ParamGroupId    uuid.UUID          `json:"param_group_id"`
	ValueType       string             `json:"value_type"`
}

type ParameterTemplatesAddParams struct {
	ParamName       string          `json:"param_name"`
	ParamGroupId    uuid.UUID       `json:"param_group_id"`
	ParamKey        string          `json:"param_key"`
	DataType        string          `json:"data_type"`
	Unit            string          `json:"Unit" gorm:"column:Unit"`
	OptionalUnits   json.RawMessage `json:"OptionalUnits" gorm:"column:OptionalUnits"`
	IsRequired      bool            `json:"is_required"`
	ValidationRules json.RawMessage `json:"validation_rules"`
	SortOrder       int             `json:"sort_order"`
	Description     string          `json:"description"`
}

type ParameterTemplatesUpdate struct {
	ParamName       *string          `json:"param_name"`
	ParamGroupId    *uuid.UUID       `json:"param_group_id"`
	ParamKey        *string          `json:"param_key"`
	DataType        *string          `json:"data_type"`
	Unit            *string          `json:"Unit" gorm:"column:Unit"`
	OptionalUnits   *json.RawMessage `json:"OptionalUnits" gorm:"column:OptionalUnits"`
	IsRequired      *bool            `json:"is_required"`
	ValidationRules *json.RawMessage `json:"validation_rules"`
	SortOrder       *int             `json:"sort_order"`
	Description     *string          `json:"description"`
}
