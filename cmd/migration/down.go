package migration

import (
	"github.com/tashfi04/printbin-server/conn"
	"github.com/tashfi04/printbin-server/migration"
	"github.com/tashfi04/printbin-server/utils"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Drop tables in database",
	Long:  `Drop tables in database`,
	Run:   downDatabase,
}

func init() {
	RootCmd.AddCommand(downCmd)
}

func downDatabase(cmd *cobra.Command, args []string) {

	utils.Logger().Infoln("Dropping database table...")

	db := conn.DB()

	if err := db.Migrator().DropTable(migration.Models...); err != nil {
		utils.Logger().Infoln(err)
	}

	utils.Logger().Infoln("Database dropped successfully!")
}
