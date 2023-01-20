package config

import (
	"time"

	"github.com/spf13/viper"
)

// AppCfg stores app related configs
type AppCfg struct {
	Host             string
	Port             int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
	HttpTimeout      time.Duration
	CorsAllowedHosts []string
	StoragePath      string
	RoomList         []string
	RoomListMap      map[string]bool
	UserPrintLimit   int
}

var app AppCfg

// LoadApp reads config from file
func LoadApp() {

	app = AppCfg{

		Host:             viper.GetString("app.host"),
		Port:             viper.GetInt("app.port"),
		ReadTimeout:      viper.GetDuration("app.read_timeout") * time.Second,
		WriteTimeout:     viper.GetDuration("app.write_timeout") * time.Second,
		IdleTimeout:      viper.GetDuration("app.idle_timeout") * time.Second,
		HttpTimeout:      viper.GetDuration("app.http_timeout") * time.Second,
		CorsAllowedHosts: viper.GetStringSlice("app.cors_allowed_hosts"),
		StoragePath:      viper.GetString("app.storage_path"),
		RoomList:         viper.GetStringSlice("app.room_list"),
		UserPrintLimit:   viper.GetInt("app.user_print_limit"),
	}

	app.RoomListMap = map[string]bool{}
	for _, v := range app.RoomList {
		app.RoomListMap[v] = true
	}
}

// App returns current app config
func App() *AppCfg {
	return &app
}
