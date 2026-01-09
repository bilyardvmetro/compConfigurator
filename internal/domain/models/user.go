package models

import "time"

type User struct {
	UserID       int64     `json:"user_id"`
	Email        string    `json:"email"`
	Nickname     *string   `json:"nickname,omitempty"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
	RoleCode     string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	PasswordHash string    `json:"-"`
}
