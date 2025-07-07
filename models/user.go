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
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Name      string    `json:"name"`
	Birthday  time.Time `json:"birthday"`
	Ideology  *Ideology `json:"ideology"` // nullable
	CreatedAt time.Time `json:"created_at"`
}
