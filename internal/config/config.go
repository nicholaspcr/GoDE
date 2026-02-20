// Package config provides shared configuration utilities for environment variable loading.
package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// Load loads the configuration from a file and the environment.
// configPath, if non-empty, points to a specific config file; otherwise the
// function searches for .<appName>.yaml in $HOME and the current directory.
func Load(appName, configPath string, cfg any) error {
	// Create a new Viper instance to avoid global state issues
	v := viper.New()

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		v.AddConfigPath(home)
		v.AddConfigPath(".env")
		v.AddConfigPath(".")

		v.SetConfigType("yaml")
		v.SetConfigName("." + appName)
	}

	// Enable environment variable support with proper mapping
	// This allows STORE_TYPE to map to store.type, STORE_SQLITE_FILEPATH to store.sqlite.filepath, etc.
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file if it exists
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// Automatically bind all environment variables by walking the struct
	bindEnvRecursive(v, reflect.TypeOf(cfg).Elem(), "")
	return v.Unmarshal(cfg)
}

// bindEnvRecursive recursively binds environment variables for all struct fields
func bindEnvRecursive(v *viper.Viper, t reflect.Type, prefix string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Get the yaml/mapstructure tag, fallback to json tag
		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			tag = field.Tag.Get("yaml")
		}
		if tag == "" {
			tag = field.Tag.Get("json")
		}

		// Handle embedded structs with squash/inline tags
		isSquashed := strings.Contains(tag, "squash") || strings.Contains(tag, "inline")
		if isSquashed {
			// For squashed/inline embedded structs, use parent's prefix
			fieldType := field.Type
			if fieldType.Kind() == reflect.Struct {
				bindEnvRecursive(v, fieldType, prefix)
			}
			continue
		}

		if tag == "" || tag == "-" {
			continue
		}

		// Handle tags with options like ",omitempty"
		if commaIdx := strings.Index(tag, ","); commaIdx != -1 {
			tag = tag[:commaIdx]
		}

		// Build the key path
		var key string
		if prefix == "" {
			key = tag
		} else {
			key = prefix + "." + tag
		}

		// If the field is a struct, recurse into it
		fieldType := field.Type
		if fieldType.Kind() == reflect.Struct {
			bindEnvRecursive(v, fieldType, key)
		} else {
			// Bind the environment variable
			envVar := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
			_ = v.BindEnv(key, envVar)
		}
	}
}
