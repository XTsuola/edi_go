package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type RegisterParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginParams struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type LogoutParams struct {
	RefreshToken string `json:"refresh"`
}

type RegisterUserData struct {
	ID          uuid.UUID          `json:"id"`
	Username    string             `json:"username"`
	Email       string             `json:"email"`
	Password    string             `json:"password"`
	IsActive    bool               `json:"si_active"`
	IsSuperuser bool               `json:"is_superuser"`
	IsStaff     bool               `json:"is_staff"`
	DataJoined  pgtype.Timestamptz `json:"data_joined"`
	Role        string             `json:"role"`
}
