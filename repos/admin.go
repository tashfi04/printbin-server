package repos

import (
	"github.com/tashfi04/printbin-server/data"
	"github.com/tashfi04/printbin-server/models"
	"gorm.io/gorm"
)

type AdminRepository interface {
	CreateUser(db *gorm.DB, users []models.User) error
}

type adminRepository struct{}

var adminRepo adminRepository

func AdminRepo() AdminRepository {
	return &adminRepo
}

func (*adminRepository) CreateUser(db *gorm.DB, users []models.User) error {

	if err := data.User().BulkCreateUser(db, users); err != nil {
		return err
	}
	return nil
}
