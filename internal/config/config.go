package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Load loads the configuration from a file and the environment.
func Load(appName string, cfg interface{}) error {
	if viper.GetString("config") != "" {
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(".env")
		viper.AddConfigPath(".")

		viper.SetConfigType("yaml")
		viper.SetConfigName("." + appName)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())
	return viper.Unmarshal(cfg)
}
