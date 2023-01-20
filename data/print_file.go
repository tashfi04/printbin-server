package data

import (
	"fmt"
	"github.com/tashfi04/printbin-server/dtos"
	"github.com/tashfi04/printbin-server/models"
	"gorm.io/gorm"
)

type printFileData struct{}

var printFile printFileData

func PrintFile() *printFileData {
	return &printFile
}

func (*printFileData) ListFiles(db *gorm.DB, currentPage, pageLimit int, status uint, searchParam string, roomList []string) ([]dtos.ListFilesInfo, error) {

	var dbResponse []dtos.ListFilesInfo

	fileModel := models.File{}
	userModel := models.User{}

	tx := db.Table(fmt.Sprintf("%s AS f", fileModel.TableName())).
		Select([]string{"f.tracking_id AS tracking_id", "f.status AS status", "f.created_at AS created_at",
			"u.team_name AS team_name", "u.room_number AS room_number"}).
		Joins(fmt.Sprintf("LEFT JOIN %s AS u ON u.id = f.user_id", userModel.TableName())).
		Where("f.status = ? AND u.room_number IN ?", status, roomList)

	if searchParam != "" {
		tx.Where("tracking_id ILIKE ? OR team_name ILIKE ?", searchParam, searchParam)
	}

	if status == models.FileStatusTypeToValue[models.FileStatusTypePending] {
		tx.Order("f.id")
	} else {
		tx.Order("f.id DESC")
	}

	if err := tx.Offset((currentPage - 1) * pageLimit).
		Limit(pageLimit).
		Find(&dbResponse).Error; err != nil {
		return nil, err
	}
	return dbResponse, nil
}

func (*printFileData) CountFiles(db *gorm.DB, status uint, searchParam string, roomList []string) (int64, error) {

	var dbResponse int64

	fileModel := models.File{}
	userModel := models.User{}

	tx := db.Table(fmt.Sprintf("%s AS f", fileModel.TableName())).
		Select("COUNT(1)").
		Joins(fmt.Sprintf("LEFT JOIN %s AS u ON u.id = f.user_id", userModel.TableName())).
		Where("f.status = ? AND u.room_number IN ?", status, roomList)

	if searchParam != "" {
		tx.Where("tracking_id ILIKE ? OR team_name ILIKE ?", searchParam, searchParam)
	}

	if err := tx.Find(&dbResponse).Error; err != nil {
		return 0, err
	}
	return dbResponse, nil
}
