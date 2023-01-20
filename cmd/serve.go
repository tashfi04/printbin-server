package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/tashfi04/printbin-server/session"

	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tashfi04/printbin-server/api"
	"github.com/tashfi04/printbin-server/config"
	"github.com/tashfi04/printbin-server/conn"
	"github.com/tashfi04/printbin-server/utils"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve starts the http server",
	Long:  `Serve starts the http server`,
	RunE:  serve,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := session.NewCookieStore(); err != nil {
			utils.Logger().Errorln("Failed to create new session cookie store: ", err)
			return err
		}
		if err := conn.ConnectDB(); err != nil {
			utils.Logger().Errorln("Failed to connect Database: ", err)
			return err
		}
		return nil
	},
}

func init() {

	serveCmd.PersistentFlags().IntP("p", "p", 8001, "port on which the server will listen for http")

	err := viper.BindPFlag("app.port", serveCmd.PersistentFlags().Lookup("p"))
	if err != nil {
		utils.Logger().Panicln("error binding flag", err)
	}

	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) error {

	appCfg := config.App()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT)

	r := chi.NewMux()
	r.Mount("/", api.Router())

	httpServer := &http.Server{

		Addr:         ":" + strconv.Itoa(appCfg.Port),
		Handler:      r,
		ReadTimeout:  appCfg.ReadTimeout,
		WriteTimeout: appCfg.WriteTimeout,
		IdleTimeout:  appCfg.IdleTimeout,
	}

	go func() {

		if err := httpServer.ListenAndServe(); err != nil {

			utils.Logger().Errorln(err)
			return
		}
	}()

	utils.Logger().Infof("Http server listening on port %d...", appCfg.Port)

	<-stop
	utils.Logger().Infoln("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpServer.Shutdown(ctx)

	utils.Logger().Infoln("Http server gracefully shutdown")

	return nil
}
