package dtos

import (
	"github.com/tashfi04/printbin-server/models"
	"time"
)

type UserFileListResp struct {
	Files      []models.File `json:"files"`
	TotalFiles int64         `json:"total_files"`
}

// ListFilesInfo db resp
type ListFilesInfo struct {
	TrackingID string    `json:"tracking_id"`
	Status     uint      `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	TeamName   string    `json:"team_name"`
	RoomNumber string    `json:"room_number"`
}

// FileInfo api resp
type FileInfo struct {
	TrackingID string                `json:"tracking_id"`
	Status     models.FileStatusType `json:"status"`
	CreatedAt  time.Time             `json:"created_at"`
	TeamName   string                `json:"team_name,omitempty"`
	RoomNumber string                `json:"room_number,omitempty"`
}

// ListFilesResp api resp
type ListFilesResp struct {
	Files      []FileInfo `json:"files"`
	TotalPages int64      `json:"total_pages"`
}
