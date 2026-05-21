package models

type UploadChunk struct {
	CurrentChunk  int  `json:"current_chunk"`
	IsAllUploaded bool `json:"is_all_uploaded"`
}
