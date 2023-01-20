package migration

import (
	"github.com/tashfi04/printbin-server/conn"
	"github.com/tashfi04/printbin-server/utils"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "migration",
	Short: "Run database migrations",
	Long:  `Migration is a tool to generate and modify database tables`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := conn.ConnectDB(); err != nil {
			utils.Logger().Errorf("Can't connect to database: %v", err)
			return err
		}
		return nil
	},
}
