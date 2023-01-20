package data

import (
	"github.com/tashfi04/printbin-server/models"
	"gorm.io/gorm"
)

type fileData struct{}

var file fileData

func File() *fileData {
	return &file
}

func (*fileData) ListUserFiles(db *gorm.DB, currentPage, pageLimit int, searchParam string, userID uint) ([]models.File, error) {

	var dbResponse []models.File

	fileModel := models.File{}

	tx := db.Table(fileModel.TableName()).
		Select([]string{"tracking_id", "status", "created_at"}).
		Where("user_id = ?", userID)

	if searchParam != "" {
		tx.Where("tracking_id ILIKE ?", searchParam)
	}

	if err := tx.Order("id DESC").
		Offset((currentPage - 1) * pageLimit).
		Limit(pageLimit).
		Find(&dbResponse).Error; err != nil {
		return nil, err
	}
	return dbResponse, nil
}

func (*fileData) CountUserFiles(db *gorm.DB, searchParam string, userID uint) (int64, error) {

	var dbResponse int64

	fileModel := models.File{}

	tx := db.Table(fileModel.TableName()).
		Select("COUNT(1)").
		Where("user_id = ?", userID)

	if searchParam != "" {
		tx.Where("tracking_id ILIKE ?", searchParam)
	}

	if err := tx.Find(&dbResponse).Error; err != nil {
		return 0, err
	}
	return dbResponse, nil
}

func (*fileData) Create(db *gorm.DB, file *models.File, printPageCount uint) error {

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(file).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&models.User{ID: file.UserID}).
		Update("print_page_count", printPageCount).Error; err != nil {
		//Update("is_submitted", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (*fileData) Exists(db *gorm.DB, id uint) (bool, error) {

	fileModel := models.File{}

	var dbResponse bool

	subQuery := db.Table(fileModel.TableName()).
		Select("1").
		Where("id = ? AND status = ?", id, models.FileStatusTypeToValue[models.FileStatusTypePending])

	if err := db.Table(fileModel.TableName()).
		Select("EXISTS (?)", subQuery).
		Find(&dbResponse).Error; err != nil {
		return false, err
	}
	return dbResponse, nil
}

func (*fileData) Update(db *gorm.DB, id uint) error {

	fileModel := models.File{
		ID: id,
	}

	if err := db.Model(&fileModel).
		Update("status", 2).Error; err != nil {
		return err
	}
	return nil
}
