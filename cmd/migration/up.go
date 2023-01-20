package migration

import (
	"github.com/spf13/cobra"
	"github.com/tashfi04/printbin-server/config"
	"github.com/tashfi04/printbin-server/conn"
	"github.com/tashfi04/printbin-server/data"
	"github.com/tashfi04/printbin-server/migration"
	"github.com/tashfi04/printbin-server/models"
	"github.com/tashfi04/printbin-server/repos"
	"github.com/tashfi04/printbin-server/utils"
	"gorm.io/gorm"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Populate tables in database",
	Long:  `Populate tables in database`,
	Run:   upDatabase,
}

func init() {
	RootCmd.AddCommand(upCmd)
}

func upDatabase(cmd *cobra.Command, args []string) {

	utils.Logger().Infoln("Populating database...")

	db := conn.DB()

	if err := db.AutoMigrate(migration.Models...); err != nil {
		utils.Logger().Infoln(err)
	}

	utils.Logger().Infoln("Creating admin...")

	createAdmin(db)

	utils.Logger().Infoln("Database populated successfully!")
}

func createAdmin(db *gorm.DB) {

	username := config.Admin().AdminUsername
	password := config.Admin().AdminPassword

	encryptedPassword, err := utils.Encrypt(password)
	if err != nil {
		utils.Logger().Errorln(err)
	}

	_, err = repos.UserRepo().GetUserByUsername(db, username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err = data.User().CreateUser(db, &models.User{
				Username: username,
				Password: encryptedPassword,
				Role:     1,
			}); err != nil {
				utils.Logger().Errorln(err)
			}
		} else {
			utils.Logger().Infoln(err)
		}
	}
}
