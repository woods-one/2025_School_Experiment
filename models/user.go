package models

import (
	"time"
)

type Ideology string

// イデオロギーの種類
const (
	Right  Ideology = "right"
	Left   Ideology = "left"
	Center Ideology = "center"
)

// ユーザーの情報
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       string    `json:"user_id"`
	Birthday     time.Time `json:"birthday"`
	Ideology     *Ideology `json:"ideology"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
