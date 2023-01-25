package repos

import (
	"github.com/tashfi04/printbin-server/data"
	"github.com/tashfi04/printbin-server/dtos"
	"github.com/tashfi04/printbin-server/models"
	"gorm.io/gorm"
	"math"
)

type PrintFileRepository interface {
	UpdateStatus(db *gorm.DB, trackingID string) error
	ListFiles(db *gorm.DB, currentPage, pageLimit int, status, searchParam string, roomList []string) (*dtos.ListFilesResp, error)
}

type printFileRepository struct{}

var printFileRepo printFileRepository

func PrintFileRepo() PrintFileRepository {
	return &printFileRepo
}

func (*printFileRepository) UpdateStatus(db *gorm.DB, trackingID string) error {

	exists, err := data.File().Exists(db, trackingID)
	if err != nil {
		return err
	}

	if !exists {
		return gorm.ErrRecordNotFound
	}

	if err = data.File().Update(db, trackingID); err != nil {
		return err
	}

	return nil
}

func (*printFileRepository) ListFiles(db *gorm.DB, currentPage, pageLimit int, status, searchParam string, roomList []string) (*dtos.ListFilesResp, error) {

	statusValue := models.FileStatusTypeToValue[models.FileStatusType(status)]

	fileList, err := data.PrintFile().ListFiles(db, currentPage, pageLimit, statusValue, searchParam, roomList)
	if err != nil {
		return nil, err
	}

	totalFiles, err := data.PrintFile().CountFiles(db, statusValue, searchParam, roomList)
	if err != nil {
		return nil, err
	}

	respFileList := make([]dtos.FileInfo, 0)

	for _, file := range fileList {

		respFile := dtos.FileInfo{
			TrackingID: file.TrackingID,
			Status:     models.ValueToFileStatusType[file.Status],
			CreatedAt:  file.CreatedAt,
			TeamName:   file.TeamName,
			RoomNumber: file.RoomNumber,
		}

		respFileList = append(respFileList, respFile)
	}

	resp := dtos.ListFilesResp{
		Files:      respFileList,
		TotalPages: int64(math.Ceil(float64(totalFiles) / float64(pageLimit))),
	}

	return &resp, nil
}
