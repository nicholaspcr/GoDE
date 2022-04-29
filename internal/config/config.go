// Package config contains all the methods related to configuration of a binary.
// Be it a CLI or a server, all configuration shall be processed through the
// methods of provided by this package.
package config

import (
	"github.com/spf13/viper"
)

// Config contains the methods used to read and process config variables for
// executables.
type Config struct {
	*viper.Viper
}

// New config is generated that is compatible with the viper interface
func New() *Config {
	return &Config{
		Viper: viper.New(),
	}
}
