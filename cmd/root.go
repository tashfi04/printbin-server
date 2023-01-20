package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tashfi04/printbin-server/cmd/migration"
	"github.com/tashfi04/printbin-server/config"
	"github.com/tashfi04/printbin-server/utils"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "printbin",
	Short: "printbin is a platform used for printing code during programming contests",
	Long: `A platform for printing code during programming contests. Using this platform, contestants can submit their
	code for printing without needing the volunteers to manually take the file to the printer for printing.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {

	utils.Logger().SetFormatter(&logrus.JSONFormatter{})
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmd.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(migration.RootCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if err := config.Init(); err != nil {
		utils.Logger().Fatalln(err)
	}
}
