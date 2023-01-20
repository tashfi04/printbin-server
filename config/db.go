package config

import (
	"github.com/spf13/viper"
	"time"
)

// DBCfg holds the database configuration
type DBCfg struct {
	Host            string
	Username        string
	Password        string
	Name            string
	Port            int
	MaxIdleConn     int
	MaxOpenConn     int
	MaxConnLifetime time.Duration
}

var database DBCfg

// LoadDB loads DB configuration
func LoadDB() {
	database = DBCfg{
		Host:            viper.GetString("database.host"),
		Username:        viper.GetString("database.user"),
		Password:        viper.GetString("database.password"),
		Name:            viper.GetString("database.name"),
		Port:            viper.GetInt("database.port"),
		MaxIdleConn:     viper.GetInt("database.max_idle_connection"),
		MaxOpenConn:     viper.GetInt("database.max_open_connection"),
		MaxConnLifetime: viper.GetDuration("database.max_connection_lifetime") * time.Second,
	}
}

// DB returns the default DB configuration
func DB() *DBCfg {
	return &database
}
