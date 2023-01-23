package data

import (
	"github.com/tashfi04/printbin-server/dtos"
	"github.com/tashfi04/printbin-server/models"
	"gorm.io/gorm"
)

type userData struct{}

var user userData

func User() *userData {
	return &user
}

func (*userData) BulkCreateUser(db *gorm.DB, users []models.User) error {

	if err := db.CreateInBatches(users, 100).Error; err != nil {
		return err
	}
	return nil
}

func (*userData) CreateUser(db *gorm.DB, user *models.User) error {

	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (*userData) GetUserByID(db *gorm.DB, userID uint) (*models.User, error) {

	dbResponse := models.User{}

	if err := db.Model(&models.User{}).
		Where("id = ?", userID).
		Take(&dbResponse).Error; err != nil {
		return nil, err
	}
	return &dbResponse, nil
}

func (*userData) GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {

	dbResponse := models.User{}

	if err := db.Model(&models.User{}).
		Where("username = ?", username).
		Take(&dbResponse).Error; err != nil {
		return nil, err
	}
	return &dbResponse, nil
}

func (*userData) GetMinimalUserInfo(db *gorm.DB, userID uint) (*models.User, error) {

	dbResponse := models.User{}

	userModel := models.User{}

	if err := db.Table(userModel.TableName()).
		Select([]string{"username", "role", "team_name", "room_number", "print_page_count"}).
		Where("id = ?", userID).
		Take(&dbResponse).Error; err != nil {
		return nil, err
	}
	return &dbResponse, nil
}
