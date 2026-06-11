package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ModelAll struct {
	ID          uuid.UUID          `json:"id" gorm:"primaryKey"`
	ModelName   string             `json:"model_name"`
	LibraryName string             `json:"library_name"`
	Version     string             `json:"version"`
	Description string             `json:"description"`
	Parameters  string             `json:"parameters"`
	ShareId     uuid.UUID          `json:"share_id"`
	Status      string             `json:"status"`
	IsPublic    bool               `json:"is_public"`
	CreatedTime pgtype.Timestamptz `json:"created_time"`
	UpdatedTime pgtype.Timestamptz `json:"updated_time"`
	IsDeleted   bool               `json:"is_deleted"`
	TotalId     int                `json:"total_id"`
	CategoryId  int                `json:"category_id"`
	UserId      uuid.UUID          `json:"user_id"`
	Extra       string             `json:"extra"`
	DocumentId  uuid.UUID          `json:"document_id"`
}
