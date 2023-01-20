package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID             uint           `gorm:"primaryKey" json:"-"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Username       string         `gorm:"not null" json:"username"`
	Password       string         `gorm:"not null" json:"password"`
	TeamName       string         `json:"team_name"`
	RoomNumber     string         `json:"room_number"`
	PrintPageCount uint           `gorm:"default:0" json:"print_page_count"`
	Role           uint           `gorm:"default:0" json:"role"`
}

func (*User) TableName() string {
	return "users"
}
