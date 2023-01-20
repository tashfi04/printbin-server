package models

import (
	"gorm.io/gorm"
	"time"
)

type FileStatusType string

const (
	FileStatusTypePending   FileStatusType = "pending"
	FileStatusTypeCompleted FileStatusType = "completed"
)

var (
	FileStatusTypeToValue = map[FileStatusType]uint{
		FileStatusTypePending:   1,
		FileStatusTypeCompleted: 2,
	}

	ValueToFileStatusType = map[uint]FileStatusType{
		1: FileStatusTypePending,
		2: FileStatusTypeCompleted,
	}
)

type File struct {
	ID         uint           `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	TrackingID string         `gorm:"not null" json:"tracking_id,omitempty"`
	Status     uint           `gorm:"default:1" json:"status,omitempty"`
	UserID     uint           `gorm:"not null" json:"user_id,omitempty"`
	User       User           `json:"user"`
}

func (*File) TableName() string {
	return "files"
}
