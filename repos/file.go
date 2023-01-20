package repos

import (
	"github.com/tashfi04/printbin-server/data"
	"github.com/tashfi04/printbin-server/dtos"
	"github.com/tashfi04/printbin-server/models"
	"gorm.io/gorm"
	"math"
)

type FileRepository interface {
	ListUserFiles(db *gorm.DB, currentPage, pageLimit int, searchParam string, userID uint) (*dtos.ListFilesResp, error)
	SubmitFile(db *gorm.DB, userID, printPageCount uint, trackingID string) error
}

type fileRepository struct{}

var fileRepo fileRepository

func FileRepo() FileRepository {
	return &fileRepo
}

func (*fileRepository) ListUserFiles(db *gorm.DB, currentPage, pageLimit int, searchParam string, userID uint) (*dtos.ListFilesResp, error) {

	fileList, err := data.File().ListUserFiles(db, currentPage, pageLimit, searchParam, userID)
	if err != nil {
		return nil, err
	}

	totalFiles, err := data.File().CountUserFiles(db, searchParam, userID)
	if err != nil {
		return nil, err
	}

	respFileList := make([]dtos.FileInfo, 0)

	for _, file := range fileList {

		respFile := dtos.FileInfo{
			TrackingID: file.TrackingID,
			Status:     models.ValueToFileStatusType[file.Status],
			CreatedAt:  file.CreatedAt,
		}

		respFileList = append(respFileList, respFile)
	}

	resp := dtos.ListFilesResp{
		Files:      respFileList,
		TotalPages: int64(math.Ceil(float64(totalFiles) / float64(pageLimit))),
	}

	return &resp, nil
}

func (*fileRepository) SubmitFile(db *gorm.DB, userID, printPageCount uint, trackingID string) error {

	file := models.File{
		TrackingID: trackingID,
		UserID:     userID,
	}

	if err := data.File().Create(db, &file, printPageCount); err != nil {
		return err
	}

	return nil
}
