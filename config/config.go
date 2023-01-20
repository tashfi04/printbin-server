package config

import (
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

// Init loads all configs
func Init() error {

	viper.SetEnvPrefix("printbin")
	viper.BindEnv("env")
	viper.BindEnv("config_name")
	viper.BindEnv("config_secret_name")
	viper.BindEnv("config_path")

	configName := viper.GetString("config_name")
	configSecretName := viper.GetString("config_secret_name")
	configPath := viper.GetString("config_path")

	if configName == "" {
		return fmt.Errorf("CONFIG_NAME missing")
	}
	if configSecretName == "" {
		return fmt.Errorf("CONFIG_SECRET_NAME missing")
	}
	if configPath == "" {
		return fmt.Errorf("CONFIG_PATH missing")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("%s named \"%s\"", err.Error(), configPath)
	}

	viper.SetConfigName(configSecretName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.MergeInConfig(); err != nil {
		return fmt.Errorf("%s named \"%s\"", err.Error(), configPath)
	}

	LoadApp()
	LoadDB()
	LoadAdmin()
	LoadSession()
	LoadCrypto()

	return nil
}
