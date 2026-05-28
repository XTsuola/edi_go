package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UploadChunk struct {
	CurrentChunk  int  `json:"current_chunk"`
	IsAllUploaded bool `json:"is_all_uploaded"`
}

type ModelPackagesAll struct {
	ID            int                `json:"id" gorm:"primaryKey"`
	FileUUID      string             `json:"file_uuid"`
	DownloadCount int                `json:"download_count"`
	FileName      string             `json:"file_name"`
	TotalChunks   int                `json:"total_chunks"`
	FileMd5       string             `json:"file_md5"`
	StorageType   string             `json:"storage_type"`
	StoragePath   string             `json:"storage_path"`
	FileSize      int8               `json:"file_size"`
	Status        string             `json:"status"`
	CreatedTime   pgtype.Timestamptz `json:"created_time"`
	UpdatedTime   pgtype.Timestamptz `json:"updated_time"`
	DeviceModelId uuid.UUID          `json:"device_model_id"`
	UploadedById  uuid.UUID          `json:"uploaded_by_id"`
}
