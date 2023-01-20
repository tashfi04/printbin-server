package config

import "github.com/spf13/viper"

// CryptoCfg stores crypto related configs
type CryptoCfg struct {
	Key string
}

var crypto CryptoCfg

// LoadCrypto reads config from file
func LoadCrypto() {

	crypto = CryptoCfg{

		Key: viper.GetString("crypto.key"),
	}
}

// Crypto returns current app config
func Crypto() *CryptoCfg {
	return &crypto
}
