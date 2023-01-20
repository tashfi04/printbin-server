package migration

import (
	"github.com/spf13/cobra"
	"github.com/tashfi04/printbin-server/conn"
	"github.com/tashfi04/printbin-server/migration"
	"github.com/tashfi04/printbin-server/utils"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets the tables in database",
	Long:  `Resets the tables in database`,
	Run:   resetDatabase,
}

func init() {
	RootCmd.AddCommand(resetCmd)
}

func resetDatabase(cmd *cobra.Command, args []string) {

	utils.Logger().Infoln("Resetting database...")

	db := conn.DB()

	if err := db.Migrator().DropTable(migration.Models...); err != nil {
		utils.Logger().Infoln(err)
	}

	if err := db.AutoMigrate(migration.Models...); err != nil {
		utils.Logger().Infoln(err)
	}

	utils.Logger().Infoln("Database reset successfully!")
}
