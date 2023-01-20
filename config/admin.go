package config

import (
	"github.com/spf13/viper"
)

// AdminCfg holds the admin panel configuration
type AdminCfg struct {
	AdminUsername string
	AdminPassword string
}

var adminPanel AdminCfg

// LoadAdmin loads admin panel configuration
func LoadAdmin() {

	adminPanel = AdminCfg{
		AdminUsername: viper.GetString("admin.username"),
		AdminPassword: viper.GetString("admin.password"),
	}
}

// Admin returns the default admin panel configuration
func Admin() *AdminCfg {
	return &adminPanel
}
