package models

type PageType[T any] struct {
	Rows  []T   `json:"rows"`
	Total int64 `json:"total"`
}

type LoginResult struct {
	Access  string    `json:"access"`
	Refresh string    `json:"refresh"`
	User    UserLogin `json:"user"`
}
