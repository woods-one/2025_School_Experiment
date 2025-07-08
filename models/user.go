package models

import (
	"time"
)

type Ideology string

const (
	Right Ideology = "right"
	Left  Ideology = "left"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       string    `json:"user_id"`
	Birthday     time.Time `json:"birthday"`
	Ideology     *Ideology `json:"ideology"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
