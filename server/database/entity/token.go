package entity

import "time"

type TokenType string

var (
	TOKENTYPE_ADMIN TokenType = "ADMIN"
)

type Token struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	Value     string    `gorm:"not null"`
	Type      TokenType `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
}
