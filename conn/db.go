package conn

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tashfi04/printbin-server/config"
	"github.com/tashfi04/printbin-server/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var db *gorm.DB

// Connect sets the db client of database using configuration cfg
func Connect(cfg *config.DBCfg) error {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.Username, cfg.Password, cfg.Name, cfg.Port)

	newLogger := logger.New(
		logrus.New(),
		logger.Config{
			SlowThreshold:             time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(cfg.MaxConnLifetime)

	db = gormDB

	utils.Logger().Infoln("Database connection successful")

	return nil
}

func ConnectDB() error {
	cfg := config.DB()
	return Connect(cfg)
}

// DB returns the database instance
func DB() *gorm.DB {
	return db
}
