package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UserBase struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type UserLogin struct {
	ID        string             `json:"id" gorm:"primaryKey"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	Role      string             `json:"role"`
	IsActive  bool               `json:"is_active"`
	LastLogin pgtype.Timestamptz `json:"last_login"`
}
type UserAll struct {
	UserLogin
	Password string `json:"password"`
}

type ChangePasswordParams struct {
	OldPassword  string `json:"old_password"`
	NewPassword  string `json:"new_password"`
	NewPassword2 string `json:"new_password2"`
}
