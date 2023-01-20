package config

import (
	"github.com/spf13/viper"
)

// SessionCfg stores session related configs
type SessionCfg struct {
	SessionKey string
}

var session SessionCfg

// LoadSession reads config from file
func LoadSession() {

	session = SessionCfg{

		SessionKey: viper.GetString("session.session_key"),
	}
}

// Session returns current app config
func Session() *SessionCfg {
	return &session
}
