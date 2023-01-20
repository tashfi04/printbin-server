package repos

import (
	"github.com/tashfi04/printbin-server/data"
	"github.com/tashfi04/printbin-server/dtos"
	"github.com/tashfi04/printbin-server/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByUsername(db *gorm.DB, username string) (*models.User, error)
	GetMinimalUserInfo(db *gorm.DB, userID uint) (*dtos.MinimalUserInfo, error)
}

type userRepository struct{}

var userRepo userRepository

func UserRepo() UserRepository {
	return &userRepo
}

func (*userRepository) GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {

	user, err := data.User().GetUserByUsername(db, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (*userRepository) GetMinimalUserInfo(db *gorm.DB, userID uint) (*dtos.MinimalUserInfo, error) {

	user, err := data.User().GetMinimalUserInfo(db, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
